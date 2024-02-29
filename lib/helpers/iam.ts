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

type NewLambdaIamRoleProps = {
  serviceName: string;
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

  return new iam.Role(scope, `${name}Role`, {
    assumedBy: new iam.ServicePrincipal(props.serviceName),
    managedPolicies: lambdaPolicies,
  });
}
