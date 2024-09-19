import * as cdk from "aws-cdk-lib";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import type { Construct } from "constructs";

import {
  TICKERS_MODEL_NAME,
  TRANSACTIONS_MODEL_NAME,
  type TableNames,
  USERS_MODEL_NAME,
} from "./helpers/db.js";

export class DatabaseStack extends cdk.Stack {
  tableNames: TableNames;

  constructor(app: Construct, id: string, props?: cdk.StackProps) {
    super(app, id, props);

    const stockTable = this.#newDynamoDBTable("Stock");

    const logsTable = this.#newDynamoDBTable("Log");

    // OLD
    const usersTable = this.#legacyDynamodbTable(USERS_MODEL_NAME);
    const tickersTable = this.#legacyDynamodbTable(TICKERS_MODEL_NAME);
    const transactionsTable = this.#legacyDynamodbTable(
      TRANSACTIONS_MODEL_NAME,
    );

    this.tableNames = {
      stocks: stockTable.tableName,
      logs: logsTable.tableName,

      // OLD
      users: usersTable.tableName,
      tickers: tickersTable.tableName,
      transactions: transactionsTable.tableName,
    };
  }

  #legacyDynamodbTable(
    modelName: string,
    props?: Partial<dynamodb.TablePropsV2>,
  ) {
    return new dynamodb.TableV2(this, `${modelName}Table`, {
      // todo STK-114 make a dynamoDB attribute method
      partitionKey: {
        name: `${modelName}Id`,
        type: dynamodb.AttributeType.STRING,
      },

      removalPolicy: cdk.RemovalPolicy.DESTROY,

      ...props,
    });
  }

  #newDynamoDBTable(modelName: string) {
    return new dynamodb.TableV2(this, `${modelName}Table`, {
      partitionKey: { type: dynamodb.AttributeType.STRING, name: "PK" },
      sortKey: { type: dynamodb.AttributeType.STRING, name: "SK" },
      globalSecondaryIndexes: [
        {
          indexName: "GSI",
          partitionKey: { type: dynamodb.AttributeType.STRING, name: "PK" },
          sortKey: { type: dynamodb.AttributeType.STRING, name: "DT" },
        },
      ],

      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });
  }
}
