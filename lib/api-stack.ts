import * as go from "@aws-cdk/aws-lambda-go-alpha";
import { CfnOutput, Stack, type StackProps } from "aws-cdk-lib";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import type { Construct } from "constructs";

import { addCorsOptions } from "./helpers/api.ts";

export class ApiStack extends Stack {
  apiUrl: string;

  constructor(app: Construct, id: string, props?: StackProps) {
    super(app, id, props);

    // auth middleware
    // TODO https://github.com/aws-samples/aws-cdk-examples/blob/master/typescript/api-gateway-lambda-token-authorizer/lib/stack/gateway-lambda-auth-stack.ts#L98C1-L127C2
    // const lambdaAuthFunction

    // user controller
    const usersControllerFunction = new go.GoFunction(
      this,
      "UsersControllerFunction",
      {
        entry: "lambdas/cmd/api-users",
      },
    );
    const usersIntegration = new apigateway.LambdaIntegration(
      usersControllerFunction,
    );
    // TODO roles

    // observability controller (checking logs)
    const logsControllerFunction = new go.GoFunction(
      this,
      "LogsControllerFunction",
      {
        entry: "lambdas/cmd/api-logs",
      },
    );
    const logsIntegration = new apigateway.LambdaIntegration(
      logsControllerFunction,
    );
    // TODO roles

    // stock indexes controller
    const stocksControllerFunction = new go.GoFunction(
      this,
      "StocksControllerFunction",
      {
        entry: "lambdas/cmd/api-stocks",
      },
    );
    const stocksIntegration = new apigateway.LambdaIntegration(
      stocksControllerFunction,
    );
    // TODO roles

    const api = new apigateway.RestApi(this, "stockAppApi", {
      restApiName: "Stock App API",
      // defaultMethodOptions: {
      //     authorizer: lambdaAuthFunction
      // }
    });

    const usersApi = api.root.addResource("users").addResource("{path+}");
    usersApi.addMethod("GET", usersIntegration);
    addCorsOptions(usersApi);

    const logsApi = api.root.addResource("logs").addResource("{path+}");
    logsApi.addMethod("GET", logsIntegration);
    addCorsOptions(logsApi);

    const stocksApi = api.root.addResource("stocks").addResource("{path+}");
    stocksApi.addMethod("GET", stocksIntegration);
    addCorsOptions(stocksApi);

    this.apiUrl = api.url;
    new CfnOutput(this, "API Url", { value: this.apiUrl });
  }
}
