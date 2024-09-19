import * as cdk from "aws-cdk-lib";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import type { Construct } from "constructs";

import type { TableNames } from "./helpers/db.js";

export class DatabaseStack extends cdk.Stack {
  tableNames: TableNames;

  constructor(app: Construct, id: string, props?: cdk.StackProps) {
    super(app, id, props);

    const stockTable = this.#newDynamoDBTable("Stock");

    const logsTable = this.#newDynamoDBTable("Log");

    this.tableNames = {
      stocks: stockTable.tableName,
      logs: logsTable.tableName,
    };
  }

  #newDynamoDBTable(modelName: string) {
    return new dynamodb.TableV2(this, `${modelName}Table`, {
      partitionKey: this.#dynamoDBAttribute("PK"),
      sortKey: this.#dynamoDBAttribute("SK"),
      globalSecondaryIndexes: [
        {
          indexName: "GSI",
          partitionKey: this.#dynamoDBAttribute("PK"),
          sortKey: this.#dynamoDBAttribute("DT"),
        },
      ],

      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });
  }

  #dynamoDBAttribute(name: string | number): dynamodb.Attribute {
    if (typeof name === "string") {
      return { type: dynamodb.AttributeType.STRING, name };
    }
    return { type: dynamodb.AttributeType.NUMBER, name: String(name) };
  }
}
