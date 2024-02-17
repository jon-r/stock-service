import * as go from '@aws-cdk/aws-lambda-go-alpha';
import { Duration, Stack, type StackProps } from 'aws-cdk-lib';
import * as iam from 'aws-cdk-lib/aws-iam';
// import * as events from 'aws-cdk-lib/aws-events';
// import * as targets from 'aws-cdk-lib/aws-events-targets';
// import * as lambdaEvents from 'aws-cdk-lib/aws-lambda-event-sources';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import type { Construct } from 'constructs';

import { DB_FULL_ACCESS_POLICY_ARN } from './helpers/db.ts';

export class DataEntryStack extends Stack {
  constructor(app: Construct, id: string, props?: StackProps) {
    super(app, id, props);

    // orchestrator lambda - creates the list of things to fetch (database?), sends the first queue item
    const managerFunctionRole = this.#newDbRole(
      'DataEntryManager',
      DB_FULL_ACCESS_POLICY_ARN,
    );

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const managerFunction = new go.GoFunction(
      this,
      'DataEntryManagerFunction',
      {
        entry: 'lambdas/dataManager',
        role: managerFunctionRole,
      },
    );

    // worker lambda - reads the list, fetches the data, queues up the next fetch, then parses the fetch result
    const workerFunctionRole = this.#newDbRole(
      'DataEntryWorker',
      DB_FULL_ACCESS_POLICY_ARN,
    );

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const workerFunction = new go.GoFunction(this, 'DataEntryWorkerFunction', {
      entry: 'lambdas/dataWorker',
      role: workerFunctionRole,
    });

    // event trigger - starts the orchestrator to trigger at regular points
    //   (daily? hourly?) maybe do it overnight so data is ready next day
    /* FIXME no scheduler for now
    const rule = new events.Rule(this, 'DataEntryScheduler', {
      schedule: events.Schedule.rate(Duration.hours(1)),
    });

    rule.addTarget(new targets.LambdaFunction(managerFunction));
    */

    // SQS - delayed events throttled to match api thresholds
    const deadLetterQueue = new sqs.Queue(this, 'DataEntryDeadLetterQueue', {
      queueName: 'DataEntryDeadLetterQueue',
      retentionPeriod: Duration.days(7),
    });

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const queue = new sqs.Queue(this, 'DataEntryQueue', {
      queueName: 'DataEntryQueue',
      visibilityTimeout: Duration.seconds(30),
      deadLetterQueue: {
        maxReceiveCount: 1,
        queue: deadLetterQueue,
      },
    });

    /* fixme no queue trigger for now
    const invokeEventSource = new lambdaEvents.SqsEventSource(queue);
    workerFunction.addEventSource(invokeEventSource);
    */
  }

  #newDbRole(name: string, policyArn: string) {
    return new iam.Role(this, `${name}Role`, {
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
      managedPolicies: [
        iam.ManagedPolicy.fromManagedPolicyArn(
          this,
          `${name}Policy`,
          policyArn,
        ),
      ],
    });
  }
}
