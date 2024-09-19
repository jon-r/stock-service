import * as lambdaGo from "@aws-cdk/aws-lambda-go-alpha";
import * as cdk from "aws-cdk-lib";
import * as events from "aws-cdk-lib/aws-events";
import * as targets from "aws-cdk-lib/aws-events-targets";
import * as logs from "aws-cdk-lib/aws-logs";
import type { Construct } from "constructs";

import { type TableNames, getDatabaseTableEnvVariables } from "./helpers/db.ts";
import {
  DB_READ_ONLY_POLICY_ARN,
  EVENTS_FULL_ACCESS_POLICY_ARN,
  LAMBDA_INVOKE_POLICY_ARN,
  SQS_FULL_ACCESS_POLICY_ARN,
  newLambdaIamRole,
} from "./helpers/iam.ts";
import {
  type DataTickerProps,
  getTickerEnvVariables,
} from "./helpers/ticker.ts";

type DataManagerStackProps = cdk.StackProps & {
  dataTickerProps: DataTickerProps;
  tableNames: TableNames;
};

export class DataManagerStack extends cdk.Stack {
  constructor(app: Construct, id: string, props: DataManagerStackProps) {
    super(app, id, props);

    // Nightly rule
    const daily1AM: events.CronOptions = { hour: "1", minute: "0" };
    const rule = new events.Rule(this, "DataManagerPoll", {
      schedule: events.Schedule.cron(daily1AM),
      enabled: true,
    });

    // Manager lambda - batches tickers to fetch latest data
    const managerFunctionRole = newLambdaIamRole(this, "DataManager", {
      policyARNs: [
        SQS_FULL_ACCESS_POLICY_ARN,
        LAMBDA_INVOKE_POLICY_ARN,
        EVENTS_FULL_ACCESS_POLICY_ARN,
        DB_READ_ONLY_POLICY_ARN,
      ],
    });

    const managerFunction = new lambdaGo.GoFunction(
      this,
      "DataManagerFunction",
      {
        entry: "lambdas/cmd/data-manager",
        role: managerFunctionRole,
        environment: {
          ...getDatabaseTableEnvVariables(props.tableNames),
          ...getTickerEnvVariables(props.dataTickerProps),
        },
        logRetention: logs.RetentionDays.THREE_MONTHS,
      },
    );

    rule.addTarget(new targets.LambdaFunction(managerFunction));
  }
}
