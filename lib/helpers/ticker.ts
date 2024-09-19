import { LambdaEnvVariables } from "./lambdas.ts";

export const TICKER_RULE_NAME = "DataEntryTickerPoll";

export type DataTickerProps = {
  eventsQueueUrl: string;
  eventPollerFunctionName: string;
};

export function getTickerEnvVariables(ticker: DataTickerProps) {
  return {
    [LambdaEnvVariables.LambdaTickerName]: ticker.eventPollerFunctionName,
    [LambdaEnvVariables.EventBridgeRuleName]: TICKER_RULE_NAME,
    [LambdaEnvVariables.SqsQueueUrl]: ticker.eventsQueueUrl,
  };
}
