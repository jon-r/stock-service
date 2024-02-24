import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";

import { DataEntryStack } from "./data-entry-stack.ts";
import { DatabaseStack } from "./database-stack.js";

export class StockAppStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    new DatabaseStack(scope, "DatabaseStack", props);

    new DataEntryStack(scope, "DataEntryStack", props);
  }
}
