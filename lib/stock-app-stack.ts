import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";

import { ApiStack } from "./api-stack.ts";
import { DataEntryStack } from "./data-entry-stack.ts";
import { DatabaseStack } from "./database-stack.js";

export class StockAppStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const { tableNames } = new DatabaseStack(scope, "DatabaseStack", props);

    new DataEntryStack(scope, "DataEntryStack", { ...props, tableNames });

    new ApiStack(scope, "ApiStack", props);
  }
}
