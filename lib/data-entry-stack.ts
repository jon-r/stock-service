import * as go from "@aws-cdk/aws-lambda-go-alpha";
import { Duration, Stack, type StackProps } from "aws-cdk-lib";
import * as events from "aws-cdk-lib/aws-events";
import * as targets from "aws-cdk-lib/aws-events-targets";
import * as sqs from "aws-cdk-lib/aws-sqs";
import type { Construct } from "constructs";

import type { TableNames } from "./helpers/db.ts";
import type { DataTickerProps } from "./helpers/events.ts";
import {
  DB_FULL_ACCESS_POLICY_ARN,
  LAMBDA_INVOKE_POLICY_ARN,
  SCHEDULER_FULL_ACCESS_POLICY_ARN,
  SQS_FULL_ACCESS_POLICY_ARN,
  newLambdaIamRole,
} from "./helpers/iam.ts";

type DataEntryStackProps = StackProps & {
  tableNames: TableNames;
};

export class DataEntryStack extends Stack {
  dataTickerProps: DataTickerProps;

  constructor(app: Construct, id: string, props: DataEntryStackProps) {
    super(app, id, props);

    const deadLetterQueue = new sqs.Queue(this, "DataEntryDeadLetterQueue", {
      queueName: "DataEntryDeadLetterQueue",
      retentionPeriod: Duration.days(7),
    });

    const queue = new sqs.Queue(this, "DataEntryQueue", {
      queueName: "DataEntryQueue.fifo",
      visibilityTimeout: Duration.seconds(30),
      fifo: true,
      deadLetterQueue: {
        maxReceiveCount: 1,
        queue: deadLetterQueue,
      },
    });

    const rule = new events.Rule(this, "DataEntryPoll", {
      schedule: events.Schedule.rate(Duration.minutes(1)),
      enabled: false,
    });

    // worker lambda - fetches and compiles third party data
    const workerFunctionRole = newLambdaIamRole(this, "DataEntryWorker", {
      policyARNs: [
        SQS_FULL_ACCESS_POLICY_ARN,
        DB_FULL_ACCESS_POLICY_ARN,
        // SCHEDULER_FULL_ACCESS_POLICY_ARN,
      ],
    });
    const workerFunction = new go.GoFunction(this, "DataEntryWorkerFunction", {
      entry: "lambdas/cmd/data-worker",
      role: workerFunctionRole,
      environment: {
        POLYGON_API_KEY: import.meta.env.VITE_POLYGON_IO_API_KEY,

        DB_LOGS_TABLE_NAME: props.tableNames.logs,
        DB_TICKERS_TABLE_NAME: props.tableNames.tickers,

        SQS_QUEUE_URL: queue.queueUrl,
        SQS_DL_QUEUE_URL: deadLetterQueue.queueUrl,

        // todo worker maybe cant invoke the ticker? how to do recursion?
        // EVENTBRIDGE_RULE_NAME: rule.ruleName,
      },
    });

    // poll lambda - reads the queue in a throttled way to pass the events on to the worker function
    const tickerFunctionRole = newLambdaIamRole(this, "DataEntryTicker", {
      policyARNs: [
        SQS_FULL_ACCESS_POLICY_ARN,
        LAMBDA_INVOKE_POLICY_ARN,
        SCHEDULER_FULL_ACCESS_POLICY_ARN,
      ],
    });
    const tickerFunction = new go.GoFunction(this, "DataEntryPollerFunction", {
      entry: "lambdas/cmd/data-ticker",
      role: tickerFunctionRole,
      // long timeout, single concurrent function only
      timeout: Duration.minutes(5),
      reservedConcurrentExecutions: 1,
      environment: {
        SQS_QUEUE_URL: queue.queueUrl,
        LAMBDA_WORKER_NAME: workerFunction.functionName,
        EVENTBRIDGE_RULE_NAME: rule.ruleName,
      },
    });

    rule.addTarget(new targets.LambdaFunction(tickerFunction));

    this.dataTickerProps = {
      eventRuleName: rule.ruleName,
      eventsQueueUrl: queue.queueUrl,
      eventPollerFunctionName: tickerFunction.functionName,
    };
  }
}
