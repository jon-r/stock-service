import * as go from '@aws-cdk/aws-lambda-go-alpha';
import { Duration, Stack, type StackProps } from 'aws-cdk-lib';
import * as events from 'aws-cdk-lib/aws-events';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import * as lambdaEvents from 'aws-cdk-lib/aws-lambda-event-sources';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import type { Construct } from 'constructs';

export class DataEntryStack extends Stack {
  constructor(app: Construct, id: string, props?: StackProps) {
    super(app, id, props);

    // orchestrator lambda - creates the list of things to fetch (database?), sends the first queue item
    const managerFunction = new go.GoFunction(
      this,
      'DataEntryManagerFunction',
      {
        entry: 'lambdas/dataManager',
      },
    );
    // todo permissions to read/write databases

    // worker lambda - reads the list, fetches the data, queues up the next fetch, then parses the fetch result
    const workerFunction = new go.GoFunction(this, 'DataEntryWorkerFunction', {
      entry: 'lambdas/dataWorker',
    });
    // todo permissions to read/write databases

    // event trigger - starts the orchestrator to trigger at regular points
    //   (daily? hourly?) maybe do it overnight so data is ready next day
    const rule = new events.Rule(this, 'DataEntryScheduler', {
      schedule: events.Schedule.rate(Duration.hours(1)),
    });

    rule.addTarget(new targets.LambdaFunction(managerFunction));

    // SQS - delayed events throttled to match api thresholds
    const deadLetterQueue = new sqs.Queue(this, 'DataEntryDeadLetterQueue', {
      queueName: 'DataEntryDeadLetterQueue',
      retentionPeriod: Duration.days(7),
    });

    const queue = new sqs.Queue(this, 'DataEntryQueue', {
      queueName: 'DataEntryQueue',
      visibilityTimeout: Duration.seconds(30),
      deadLetterQueue: {
        maxReceiveCount: 1,
        queue: deadLetterQueue,
      },
    });

    const invokeEventSource = new lambdaEvents.SqsEventSource(queue);
    workerFunction.addEventSource(invokeEventSource);
  }
}
