import * as go from '@aws-cdk/aws-lambda-go-alpha';
import { Duration, Stack, type StackProps } from 'aws-cdk-lib';
import * as iam from 'aws-cdk-lib/aws-iam';
// import * as events from 'aws-cdk-lib/aws-events';
// import * as targets from 'aws-cdk-lib/aws-events-targets';
// import * as lambdaEvents from 'aws-cdk-lib/aws-lambda-event-sources';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import type { Construct } from 'constructs';

import {
  DB_FULL_ACCESS_POLICY_ARN,
  SQS_FULL_ACCESS_POLICY_ARN,
} from './helpers/iam.ts';

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

    const managerFunctionRole = this.#newIamRole(
      'DataEntryManager',
      'lambda.amazonaws.com',
      [DB_FULL_ACCESS_POLICY_ARN, SQS_FULL_ACCESS_POLICY_ARN],
    );

    // orchestrator lambda - creates the list of things to fetch (database?), sends the first queue item
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const managerFunction = new go.GoFunction(
      this,
      'DataEntryManagerFunction',
      {
        entry: 'lambdas/dataManager',
        role: managerFunctionRole,
        environment: {
          // eslint-disable-next-line @typescript-eslint/naming-convention
          SQS_QUEUE_URL: queue.queueUrl,
        },
      },
    );

    const workerFunctionRole = this.#newIamRole(
      'DataEntryWorker',
      'lambda.amazonaws.com',
      [DB_FULL_ACCESS_POLICY_ARN],
    );

    // worker lambda - reads the list, fetches the data, queues up the next fetch, then parses the fetch result
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const workerFunction = new go.GoFunction(this, 'DataEntryWorkerFunction', {
      entry: 'lambdas/dataWorker',
      role: workerFunctionRole,
    });

    /* fixme no queue trigger for now
    const invokeEventSource = new lambdaEvents.SqsEventSource(queue);
    workerFunction.addEventSource(invokeEventSource);
    */
  }

  // todo prob move this to a helper file, the api will need it also
  #newIamRole(name: string, service: string, policyArns: string[]) {
    return new iam.Role(this, `${name}Role`, {
      assumedBy: new iam.ServicePrincipal(service),
      managedPolicies: policyArns.map((policyArn, i) =>
        iam.ManagedPolicy.fromManagedPolicyArn(
          this,
          `${name}Policy_${i}`,
          policyArn,
        ),
      ),
    });
  }
}
