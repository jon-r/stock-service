import * as go from "@aws-cdk/aws-lambda-go-alpha";
import { Duration, Stack, type StackProps } from "aws-cdk-lib";
import * as events from "aws-cdk-lib/aws-events";
import * as targets from "aws-cdk-lib/aws-events-targets";
import * as sqs from "aws-cdk-lib/aws-sqs";
import type { Construct } from "constructs";

import type { TableNames } from "./helpers/db.ts";
import {
  DB_FULL_ACCESS_POLICY_ARN,
  LAMBDA_INVOKE_POLICY_ARN,
  SQS_FULL_ACCESS_POLICY_ARN,
  newLambdaIamRole,
} from "./helpers/iam.ts";

type DataEntryStackProps = StackProps & {
  tableNames: TableNames;
};

export class DataEntryStack extends Stack {
  constructor(app: Construct, id: string, props: DataEntryStackProps) {
    super(app, id, props);

    // SQS - delayed events throttled to match remote thresholds
    const deadLetterQueue = new sqs.Queue(this, "DataEntryDeadLetterQueue", {
      queueName: "DataEntryDeadLetterQueue",
      retentionPeriod: Duration.days(7),
    });

    const queue = new sqs.Queue(this, "DataEntryQueue", {
      queueName: "DataEntryQueue",
      visibilityTimeout: Duration.seconds(30),
      deadLetterQueue: {
        maxReceiveCount: 1,
        queue: deadLetterQueue,
      },
    });

    // worker lambda - fetches and compiles third party data
    const workerFunctionRole = newLambdaIamRole(this, "DataEntryWorker", {
      serviceName: "lambda.amazonaws.com",
      policyARNs: [DB_FULL_ACCESS_POLICY_ARN],
    });
    const workerFunction = new go.GoFunction(this, "DataEntryWorkerFunction", {
      entry: "lambdas/cmd/data-worker",
      role: workerFunctionRole,
      environment: {
        DB_LOGS_TABLE_NAME: props.tableNames.logs,
        DB_TICKERS_TABLE_NAME: props.tableNames.tickers,

        SQS_QUEUE_URL: queue.queueUrl,
        // todo add failed items to DL queue
        //  https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/sqs-example-dead-letter-queues.html
        SQS_DL_QUEUE_URL: deadLetterQueue.queueUrl,
      },
    });

    // poll lambda - reads the queue in a throttled way to pass the events on to the worker function
    const rule = new events.Rule(this, "DataEntryPoll", {
      schedule: events.Schedule.rate(Duration.minutes(1)),
        // todo be disabled initially?
    });
    const tickerFunctionRole = newLambdaIamRole(this, "DataEntryTicker", {
      serviceName: "lambda.amazonaws.com",
      policyARNs: [SQS_FULL_ACCESS_POLICY_ARN, LAMBDA_INVOKE_POLICY_ARN],
    });

    const tickerFunction = new go.GoFunction(this, "DataEntryPollerFunction", {
      entry: "lambdas/cmd/data-ticker",
      role: tickerFunctionRole,
      timeout: Duration.minutes(5),
      environment: {
        SQS_QUEUE_URL: queue.queueUrl,
        LAMBDA_WORKER_NAME: workerFunction.functionName,
      },
    });

    rule.addTarget(new targets.LambdaFunction(tickerFunction));
  }
}
