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
      queueName: "DataEntryQueue",
      visibilityTimeout: Duration.seconds(30),
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
        SCHEDULER_FULL_ACCESS_POLICY_ARN,
      ],
    });
    const workerFunction = new go.GoFunction(this, "DataEntryWorkerFunction", {
      entry: "lambdas/cmd/data-worker",
      role: workerFunctionRole,
      environment: {
        POLYGON_API_KEY: import.meta.env.VITE_POLYGON_IO_API_KEY,

        DB_LOGS_TABLE_NAME: props.tableNames.logs,
        DB_TICKERS_TABLE_NAME: props.tableNames.tickers,

        // todo add failed items to DL queue
        //  https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/sqs-example-dead-letter-queues.html
        SQS_QUEUE_URL: queue.queueUrl,
        SQS_DL_QUEUE_URL: deadLetterQueue.queueUrl,

        EVENTBRIDGE_RULE_ARN: rule.ruleArn,
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
        EVENTBRIDGE_RULE_ARN: rule.ruleArn,
      },
    });

    rule.addTarget(new targets.LambdaFunction(tickerFunction));

    this.dataTickerProps = {
      eventRuleArn: rule.ruleArn,
      eventsQueueUrl: queue.queueUrl,
    };
  }
}
