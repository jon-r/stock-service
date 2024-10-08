import * as lambdaGo from "@aws-cdk/aws-lambda-go-alpha";
import * as cdk from "aws-cdk-lib";
import * as events from "aws-cdk-lib/aws-events";
import * as targets from "aws-cdk-lib/aws-events-targets";
import * as logs from "aws-cdk-lib/aws-logs";
import * as sqs from "aws-cdk-lib/aws-sqs";
import type { Construct } from "constructs";

import { type TableNames, getDatabaseTableEnvVariables } from "./helpers/db.ts";
import {
  DB_FULL_ACCESS_POLICY_ARN,
  EVENTS_FULL_ACCESS_POLICY_ARN,
  LAMBDA_INVOKE_POLICY_ARN,
  SQS_FULL_ACCESS_POLICY_ARN,
  newLambdaIamRole,
} from "./helpers/iam.ts";
import {
  type DataTickerProps,
  TICKER_RULE_NAME,
  getTickerEnvVariables,
} from "./helpers/ticker.ts";

type DataEntryStackProps = cdk.StackProps & {
  tableNames: TableNames;
};

export class DataEntryStack extends cdk.Stack {
  dataTickerProps: DataTickerProps;

  constructor(app: Construct, id: string, props: DataEntryStackProps) {
    super(app, id, props);

    const deadLetterQueue = new sqs.Queue(this, "DataEntryDeadLetterQueue", {
      queueName: "DataEntryDeadLetterQueue",
      retentionPeriod: cdk.Duration.days(7),
    });

    const queue = new sqs.Queue(this, "DataEntryQueue", {
      queueName: "DataEntryQueue",
      visibilityTimeout: cdk.Duration.minutes(4),
      deadLetterQueue: {
        maxReceiveCount: 1,
        queue: deadLetterQueue,
      },
    });

    const rule = new events.Rule(this, "DataEntryPoll", {
      schedule: events.Schedule.rate(cdk.Duration.minutes(2)),
      ruleName: TICKER_RULE_NAME,
      enabled: false,
    });

    // worker lambda - fetches and compiles third party data
    const workerFunctionRole = newLambdaIamRole(this, "DataEntryWorker", {
      policyARNs: [
        SQS_FULL_ACCESS_POLICY_ARN,
        DB_FULL_ACCESS_POLICY_ARN,
        EVENTS_FULL_ACCESS_POLICY_ARN,
      ],
    });
    const workerFunction = new lambdaGo.GoFunction(
      this,
      "DataEntryWorkerFunction",
      {
        entry: "lambdas/cmd/data-worker",
        role: workerFunctionRole,
        environment: {
          ...getDatabaseTableEnvVariables(props.tableNames),
          ...getTickerEnvVariables({
            eventRuleName: TICKER_RULE_NAME,
            eventsQueueUrl: queue.queueUrl,
            eventPollerFunctionName: "",
          }),

          POLYGON_API_KEY: import.meta.env.VITE_POLYGON_IO_API_KEY,

          SQS_DLQ_URL: deadLetterQueue.queueUrl,
        },
        logRetention: logs.RetentionDays.THREE_MONTHS,
      },
    );

    // poll lambda - reads the queue in a throttled way to pass the events on to the worker function
    const tickerFunctionRole = newLambdaIamRole(this, "DataEntryTicker", {
      policyARNs: [
        SQS_FULL_ACCESS_POLICY_ARN,
        LAMBDA_INVOKE_POLICY_ARN,
        EVENTS_FULL_ACCESS_POLICY_ARN,
      ],
    });
    const tickerTimeout = 5;
    const tickerFunction = new lambdaGo.GoFunction(
      this,
      "DataEntryPollerFunction",
      {
        entry: "lambdas/cmd/data-ticker",
        role: tickerFunctionRole,
        // Long timeout, single concurrent function only
        timeout: cdk.Duration.minutes(tickerTimeout + 0.1),
        reservedConcurrentExecutions: 1,
        // Dont reattempt
        retryAttempts: 0,
        environment: {
          ...getTickerEnvVariables({
            eventRuleName: TICKER_RULE_NAME,
            eventsQueueUrl: queue.queueUrl,
            eventPollerFunctionName: "", // Wont self invoke
          }),

          TICKER_TIMEOUT: String(tickerTimeout),
          LAMBDA_WORKER_NAME: workerFunction.functionName,
        },
        logRetention: logs.RetentionDays.THREE_MONTHS,
      },
    );

    rule.addTarget(new targets.LambdaFunction(tickerFunction));

    this.dataTickerProps = {
      eventRuleName: TICKER_RULE_NAME,
      eventsQueueUrl: queue.queueUrl,
      eventPollerFunctionName: tickerFunction.functionName,
    };
  }
}
