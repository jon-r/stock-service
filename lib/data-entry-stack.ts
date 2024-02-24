import * as go from "@aws-cdk/aws-lambda-go-alpha";
import { Duration, Stack, type StackProps } from "aws-cdk-lib";
// import * as events from 'aws-cdk-lib/aws-events';
// import * as targets from 'aws-cdk-lib/aws-events-targets';
import * as lambdaEvents from "aws-cdk-lib/aws-lambda-event-sources";
import * as sqs from "aws-cdk-lib/aws-sqs";
import type { Construct } from "constructs";

import {
  DB_FULL_ACCESS_POLICY_ARN,
  SQS_FULL_ACCESS_POLICY_ARN,
  newLambdaIamRole,
} from "./helpers/iam.ts";

export class DataEntryStack extends Stack {
  constructor(app: Construct, id: string, props?: StackProps) {
    super(app, id, props);

    // event trigger - starts the orchestrator to trigger at regular points
    //   (daily? hourly?) maybe do it overnight so data is ready next day
    /* FIXME no scheduler for now
    const rule = new events.Rule(this, 'DataEntryScheduler', {
      schedule: events.Schedule.rate(Duration.hours(1)),
    });

    rule.addTarget(new targets.LambdaFunction(managerFunction));
    */

    // SQS - delayed events throttled to match api thresholds
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

    // orchestrator lambda - creates the list of things to fetch, sends the first queue item
    const managerFunctionRole = newLambdaIamRole(this, "DataEntryManager", {
      serviceName: "lambda.amazonaws.com",
      policyARNs: [DB_FULL_ACCESS_POLICY_ARN, SQS_FULL_ACCESS_POLICY_ARN],
    });
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const managerFunction = new go.GoFunction(
      this,
      "DataEntryManagerFunction",
      {
        entry: "lambdas/cmd/data-manager",
        role: managerFunctionRole,
        environment: {
          SQS_QUEUE_URL: queue.queueUrl,
        },
      },
    );

    // worker lambda - reads the list, fetches the data, queues up the next fetch, then parses the fetch result
    const workerFunctionRole = newLambdaIamRole(this, "DataEntryWorker", {
      serviceName: "lambda.amazonaws.com",
      policyARNs: [DB_FULL_ACCESS_POLICY_ARN, SQS_FULL_ACCESS_POLICY_ARN],
    });
    const workerFunction = new go.GoFunction(this, "DataEntryWorkerFunction", {
      entry: "lambdas/cmd/data-worker",
      role: workerFunctionRole,
      environment: {
        SQS_QUEUE_URL: queue.queueUrl,
        // todo add failed items to DL queue
        //  https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/sqs-example-dead-letter-queues.html
        SQS_DL_QUEUE_URL: deadLetterQueue.queueUrl,
      },
    });
    const invokeEventSource = new lambdaEvents.SqsEventSource(queue);
    workerFunction.addEventSource(invokeEventSource);
  }
}
