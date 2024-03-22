import {Duration, Stack, type StackProps} from "aws-cdk-lib";
import type {Construct} from "constructs";
import type {TableNames} from "./helpers/db.ts";
import * as go from "@aws-cdk/aws-lambda-go-alpha";
import * as events from "aws-cdk-lib/aws-events";
import * as targets from "aws-cdk-lib/aws-events-targets";

type DataManagerStackProps = StackProps & {
    tableNames: TableNames
}

class DataManagerStack extends Stack {
    constructor(app: Construct, id: string, props: DataManagerStackProps) {
        super(app, id, props);

        // nightly rule
        const rule = new events.Rule(this, "DataManagerPoll", {
            schedule: events.Schedule.rate(Duration.days(1)),
        })

        // manager lambda - batches tickers to fetch latest data
    }
}