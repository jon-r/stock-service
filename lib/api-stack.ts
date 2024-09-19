import * as lambdaGo from "@aws-cdk/aws-lambda-go-alpha";
import * as cdk from "aws-cdk-lib";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as logs from "aws-cdk-lib/aws-logs";
import type { Construct } from "constructs";

import { type TableNames, getDatabaseTableEnvVariables } from "./helpers/db.ts";
import {
  DB_FULL_ACCESS_POLICY_ARN,
  EVENTS_FULL_ACCESS_POLICY_ARN,
  type KnownEnvVariables,
  LAMBDA_INVOKE_POLICY_ARN,
  type LambdaTarget,
  SQS_FULL_ACCESS_POLICY_ARN,
  newLambdaIamRole,
} from "./helpers/lambdas.ts";
import {
  type DataTickerProps,
  getTickerEnvVariables,
} from "./helpers/ticker.ts";

type ApiStackProps = cdk.StackProps & {
  dataTickerProps: DataTickerProps;
  tableNames: TableNames;
};

export class ApiStack extends cdk.Stack {
  apiUrl: string;
  apiLambdas: LambdaTarget[];

  constructor(app: Construct, id: string, props: ApiStackProps) {
    super(app, id, props);

    // Auth middleware
    // TODO https://github.com/aws-samples/aws-cdk-examples/blob/master/typescript/api-gateway-lambda-token-authorizer/lib/stack/gateway-lambda-auth-stack.ts#L98C1-L127C2
    // const lambdaAuthFunction

    // User controller
    const usersControllerFunction = new lambdaGo.GoFunction(
      this,
      "UsersControllerFunction",
      {
        entry: "lambdas/cmd/api-users",
        environment: {
          ...getDatabaseTableEnvVariables(props.tableNames),
        } satisfies KnownEnvVariables,
        logRetention: logs.RetentionDays.THREE_MONTHS,
      },
    );
    const usersIntegration = new apigateway.LambdaIntegration(
      usersControllerFunction,
    );
    // TODO roles

    // Observability controller (checking logs)
    const logsControllerFunction = new lambdaGo.GoFunction(
      this,
      "LogsControllerFunction",
      {
        entry: "lambdas/cmd/api-logs",
        environment: {
          ...getDatabaseTableEnvVariables(props.tableNames),
        } satisfies KnownEnvVariables,
        logRetention: logs.RetentionDays.THREE_MONTHS,
      },
    );
    const logsIntegration = new apigateway.LambdaIntegration(
      logsControllerFunction,
    );
    // TODO roles

    // stock indexes controller
    const stocksControllerFunctionRole = newLambdaIamRole(
      this,
      "DataEntryWorker",
      {
        policyARNs: [
          SQS_FULL_ACCESS_POLICY_ARN,
          DB_FULL_ACCESS_POLICY_ARN,
          EVENTS_FULL_ACCESS_POLICY_ARN,
          LAMBDA_INVOKE_POLICY_ARN,
        ],
      },
    );
    const stocksControllerFunction = new lambdaGo.GoFunction(
      this,
      "StocksControllerFunction",
      {
        entry: "lambdas/cmd/api-stocks",
        role: stocksControllerFunctionRole,
        environment: {
          ...getTickerEnvVariables(props.dataTickerProps),
          ...getDatabaseTableEnvVariables(props.tableNames),
        } satisfies KnownEnvVariables,
        logRetention: logs.RetentionDays.THREE_MONTHS,
      },
    );
    const stocksIntegration = new apigateway.LambdaIntegration(
      stocksControllerFunction,
    );

    const api = new apigateway.RestApi(this, "stockAppApi", {
      restApiName: "Stock App API",
      // defaultMethodOptions: {
      //     authorizer: lambdaAuthFunction
      // }
    });

    const usersApi = api.root.addResource("users").addResource("{path+}");
    usersApi.addMethod("GET", usersIntegration);
    this.#addCorsOptions(usersApi);

    const logsApi = api.root.addResource("logs").addResource("{path+}");
    logsApi.addMethod("GET", logsIntegration);
    this.#addCorsOptions(logsApi);

    const stocksApi = api.root.addResource("stocks").addResource("{path+}");
    stocksApi.addMethod("GET", stocksIntegration);
    stocksApi.addMethod("POST", stocksIntegration);
    this.#addCorsOptions(stocksApi);

    this.apiUrl = api.url;
    new cdk.CfnOutput(this, "API Url", { value: this.apiUrl });

    this.apiLambdas = [
      {
        name: "api-users",
        path: "lambdas/cmd/api-users",
        arn: usersControllerFunction.functionArn,
      },
      {
        name: "api-logs",
        path: "lambdas/cmd/api-logs",
        arn: logsControllerFunction.functionArn,
      },
      {
        name: "api-stocks",
        path: "lambdas/cmd/api-stocks",
        arn: stocksControllerFunction.functionArn,
      },
    ];
  }

  #addCorsOptions(apiResource: apigateway.IResource) {
    apiResource.addMethod(
      "OPTIONS",
      new apigateway.MockIntegration({
        // In case you want to use binary media types, uncomment the following line
        // contentHandling: ContentHandling.CONVERT_TO_TEXT,
        integrationResponses: [
          {
            statusCode: "200",
            responseParameters: {
              "method.response.header.Access-Control-Allow-Headers":
                "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent'",
              "method.response.header.Access-Control-Allow-Origin": "'*'",
              "method.response.header.Access-Control-Allow-Credentials":
                "'false'",
              "method.response.header.Access-Control-Allow-Methods":
                "'OPTIONS,GET,PUT,POST,DELETE'",
            },
          },
        ],
        // In case you want to use binary media types, comment out the following line
        passthroughBehavior: apigateway.PassthroughBehavior.NEVER,
        requestTemplates: {
          "application/json": '{"statusCode": 200}',
        },
      }),
      {
        methodResponses: [
          {
            statusCode: "200",
            responseParameters: {
              "method.response.header.Access-Control-Allow-Headers": true,
              "method.response.header.Access-Control-Allow-Methods": true,
              "method.response.header.Access-Control-Allow-Credentials": true,
              "method.response.header.Access-Control-Allow-Origin": true,
            },
          },
        ],
      },
    );
  }
}
