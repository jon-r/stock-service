import {Stack, StackProps} from "aws-cdk-lib";
import {Construct} from "constructs";
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import {AttributeType} from 'aws-cdk-lib/aws-dynamodb';
import {LOGS_MODEL_NAME} from "./helpers/db.js";

export class DatabaseStack extends Stack {
  constructor(app: Construct, id: string, props?: StackProps) {
    super(app, id, props);

    // logs table
    new dynamodb.Table(this, `${LOGS_MODEL_NAME}Table`, {
      partitionKey: {
        name: `${LOGS_MODEL_NAME}Id`,
        type: AttributeType.STRING
      },
      tableName: LOGS_MODEL_NAME,
    })

    // users table

    // stock indexes table

    // transactions table
  }

}
