import * as iam from "aws-cdk-lib/aws-iam";
import type { Construct } from "constructs";

export const SQS_FULL_ACCESS_POLICY_ARN =
  "arn:aws:iam::aws:policy/AmazonSQSFullAccess";
export const SQS_READ_ONLY_POLICY_ARN =
  "arn:aws:iam::aws:policy/AmazonSQSReadOnlyAccess";

export const DB_FULL_ACCESS_POLICY_ARN =
  "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess";
export const DB_READ_ONLY_POLICY_ARN =
  "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess";

export const LAMBDA_BASIC_POLICY_ARN =
  "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole";
export const LAMBDA_INVOKE_POLICY_ARN =
  "arn:aws:iam::aws:policy/service-role/AWSLambdaRole";

export const EVENTS_FULL_ACCESS_POLICY_ARN =
  "arn:aws:iam::aws:policy/CloudWatchEventsFullAccess";

export const enum LambdaEnvVariables {
  EventBridgeRuleName = "EVENTBRIDGE_RULE_NAME",
  DbLogsTableName = "DB_LOGS_TABLE_NAME",
  DbStocksTableName = "DB_STOCKS_TABLE_NAME",
  LambdaTickerName = "LAMBDA_TICKER_NAME",
  LambdaTickerTimeout = "TICKER_TIMEOUT",
  LambdaWorkerName = "LAMBDA_WORKER_NAME",
  PolygonApiKey = "POLYGON_API_KEY",
  SqsDlqUrl = "SQS_DLQ_URL",
  SqsQueueUrl = "SQS_QUEUE_URL",
}

// this means we dont mistype env variables
export type KnownEnvVariables = Partial<{
  [key in LambdaEnvVariables]: string;
}>;

type NewLambdaIamRoleProps = {
  policyARNs: string[];
};

export function newLambdaIamRole(
  scope: Construct,
  name: string,
  props: NewLambdaIamRoleProps,
) {
  const lambdaPolicies = [...props.policyARNs, LAMBDA_BASIC_POLICY_ARN].map(
    (policyArn, i) =>
      iam.ManagedPolicy.fromManagedPolicyArn(
        scope,
        `${name}Policy_${i}`,
        policyArn,
      ),
  );

  const principle = new iam.ServicePrincipal(
    "lambda.amazonaws.com",
  ) as iam.IPrincipal;

  return new iam.Role(scope, `${name}Role`, {
    assumedBy: principle,
    managedPolicies: lambdaPolicies,
  }) as iam.IRole;
}
