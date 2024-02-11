import { Stack, StackProps } from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import { Construct } from 'constructs';

import {
  LOGS_MODEL_NAME,
  STOCK_INDEXES_MODEL_NAME,
  TRANSACTIONS_MODEL_NAME,
  USERS_MODEL_NAME,
} from './helpers/db.js';

export class DatabaseStack extends Stack {
  constructor(app: Construct, id: string, props?: StackProps) {
    super(app, id, props);

    // logs table
    new dynamodb.Table(this, `${LOGS_MODEL_NAME}Table`, {
      partitionKey: {
        name: `${LOGS_MODEL_NAME}Id`,
        type: dynamodb.AttributeType.STRING,
      },
      tableName: LOGS_MODEL_NAME,
    });

    // users table
    new dynamodb.Table(this, `${USERS_MODEL_NAME}Table`, {
      partitionKey: {
        name: `${USERS_MODEL_NAME}Id`,
        type: dynamodb.AttributeType.STRING,
      },
      tableName: USERS_MODEL_NAME,
    });

    // stock indexes table
    new dynamodb.Table(this, `${STOCK_INDEXES_MODEL_NAME}Table`, {
      partitionKey: {
        name: `${STOCK_INDEXES_MODEL_NAME}Id`,
        type: dynamodb.AttributeType.STRING,
      },
      tableName: STOCK_INDEXES_MODEL_NAME,
    });

    // transactions table
    new dynamodb.Table(this, `${TRANSACTIONS_MODEL_NAME}Table`, {
      partitionKey: {
        name: `${TRANSACTIONS_MODEL_NAME}Id`,
        type: dynamodb.AttributeType.STRING,
      },
      tableName: TRANSACTIONS_MODEL_NAME,
    });
  }
}
