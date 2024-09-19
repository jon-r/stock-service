import * as cdk from "aws-cdk-lib";
import * as codebuild from "aws-cdk-lib/aws-codebuild";
import * as s3 from "aws-cdk-lib/aws-s3";
import { Construct } from "constructs";

import type { LambdaTarget } from "./helpers/lambdas.ts";

type LambdaPipelineStackProps = cdk.StackProps & {
  lambdaTargets: LambdaTarget[];
};

// https://aws.amazon.com/blogs/devops/complete-ci-cd-with-aws-codecommit-aws-codebuild-aws-codedeploy-and-aws-codepipeline/
export class LambdaPipelineStack extends cdk.Stack {
  constructor(app: Construct, id: string, props: LambdaPipelineStackProps) {
    super(app, id, props);

    // s3 to store the artifacts
    const artifactBucket = new s3.Bucket(this, "ArtifactBucket", {
      // todo remove once this works
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
    });

    // trigger build when github main branch updates
    const githubSource = codebuild.Source.gitHub({
      owner: import.meta.env.VITE_GITHUB_NAME,
      repo: import.meta.env.VITE_GITHUB_REPO,
      webhookFilters: [
        codebuild.FilterGroup.inEventOf(codebuild.EventAction.PUSH).andBranchIs(
          "main",
        ),
      ],
    });

    // codebuild create & save each lambda artifact
    props.lambdaTargets.map((fn) => {
      new codebuild.Project(this, `Codebuild_${fn.name}`, {
        projectName: `${fn.name}_builder`,
        buildSpec: this.#goLambdaBuildSpec(fn),
        source: githubSource,
        environment: {
          computeType: codebuild.ComputeType.LAMBDA_2GB,
          buildImage: codebuild.LinuxArmLambdaBuildImage.AMAZON_LINUX_2_GO_1_21,
        },

        artifacts: codebuild.Artifacts.s3({
          bucket: artifactBucket,
          name: fn.name,
        }),
      });
    });

    // codedeploy to (manually) update the lambdas
  }

  #goLambdaBuildSpec(target: LambdaTarget) {
    const outputFile = `./build/${target.name}`;

    return codebuild.BuildSpec.fromObject({
      version: "0.2",
      phases: {
        build: {
          commands: [`go build -o ${outputFile} ./${target.path}`],
        },
      },
      artifacts: {
        files: [outputFile],
      },
    });
  }
}
