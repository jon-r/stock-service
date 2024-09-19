import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";

import { ApiStack } from "./api-stack.ts";
import { DataEntryStack } from "./data-entry-stack.ts";
import { DataManagerStack } from "./data-manager-stack.ts";
import { DatabaseStack } from "./database-stack.js";
import { LambdaPipelineStack } from "./lambda-pipeline-stack.ts";

export class StockAppStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const { tableNames } = new DatabaseStack(scope, "DatabaseStack", props);

    const { dataTickerProps, dataEntryLambdas } = new DataEntryStack(
      scope,
      "DataEntryStack",
      {
        ...props,
        tableNames,
      },
    );

    const { dataManagerLambdas } = new DataManagerStack(
      scope,
      "DataManagerStack",
      {
        ...props,
        dataTickerProps,
        tableNames,
      },
    );

    const { apiLambdas } = new ApiStack(scope, "ApiStack", {
      ...props,
      dataTickerProps,
      tableNames,
    });

    new LambdaPipelineStack(scope, "LambdaPipelineStack", {
      ...props,
      lambdaTargets: [
        ...dataEntryLambdas,
        ...dataManagerLambdas,
        ...apiLambdas,
      ],
    });
  }
}
