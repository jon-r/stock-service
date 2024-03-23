export const TICKER_RULE_NAME = "DataEntryTickerPoll";

export type DataTickerProps = {
  eventRuleName: string;
  eventsQueueUrl: string;
  eventPollerFunctionName: string;
};

export function getTickerEnvVariables(ticker: DataTickerProps) {
  return {
    LAMBDA_TICKER_NAME: ticker.eventPollerFunctionName,
    EVENTBRIDGE_RULE_NAME: ticker.eventRuleName,
    SQS_QUEUE_URL: ticker.eventsQueueUrl,
  };
}