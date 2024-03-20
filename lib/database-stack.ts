import { Stack, type StackProps } from "aws-cdk-lib";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import type { TablePropsV2 } from "aws-cdk-lib/aws-dynamodb";
import type { Construct } from "constructs";

import {
  LOGS_MODEL_NAME,
  TICKERS_MODEL_NAME,
  TRANSACTIONS_MODEL_NAME,
  type TableNames,
  USERS_MODEL_NAME,
} from "./helpers/db.js";

export class DatabaseStack extends Stack {
  tableNames: TableNames;

  constructor(app: Construct, id: string, props?: StackProps) {
    super(app, id, props);

    // logs table
    const logsTable = this.#newDynamodbTable(LOGS_MODEL_NAME);

    // users table
    const usersTable = this.#newDynamodbTable(USERS_MODEL_NAME);

    const tickersTable = this.#newDynamodbTable(TICKERS_MODEL_NAME);

    const transactionsTable = this.#newDynamodbTable(TRANSACTIONS_MODEL_NAME);

    this.tableNames = {
      logs: logsTable.tableName,
      users: usersTable.tableName,
      tickers: tickersTable.tableName,
      transactions: transactionsTable.tableName,
    };
  }

  #newDynamodbTable(modelName: string, props?: Partial<TablePropsV2>) {
    return new dynamodb.TableV2(this, `${modelName}Table`, {
      partitionKey: {
        name: `${modelName}Id`,
        type: dynamodb.AttributeType.STRING,
      },
      ...props,
    });
  }
}
