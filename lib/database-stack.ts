import { Stack, type StackProps } from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import type { Construct } from 'constructs';

import {
  JOBS_MODEL_NAME,
  LOGS_MODEL_NAME,
  STOCK_INDEXES_MODEL_NAME,
  TRANSACTIONS_MODEL_NAME,
  USERS_MODEL_NAME,
} from './helpers/db.js';

export class DatabaseStack extends Stack {
  constructor(app: Construct, id: string, props?: StackProps) {
    super(app, id, props);

    // logs table
    this.#newDynamodbTable(LOGS_MODEL_NAME);

    // users table
    this.#newDynamodbTable(USERS_MODEL_NAME);

    // stock indexes table
    this.#newDynamodbTable(STOCK_INDEXES_MODEL_NAME);

    // transactions table
    this.#newDynamodbTable(TRANSACTIONS_MODEL_NAME);

    // queued api jobs
    this.#newDynamodbTable(JOBS_MODEL_NAME);
  }

  #newDynamodbTable(modelName: string) {
    return new dynamodb.Table(this, `${modelName}Table`, {
      partitionKey: {
        name: `${modelName}Id`,
        type: dynamodb.AttributeType.STRING,
      },
      tableName: `stock-app_${modelName}`,
    });
  }
}
