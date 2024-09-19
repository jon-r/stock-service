import { LambdaEnvVariables } from "./lambdas.ts";

export type TableNames = {
  logs: string;
  stocks: string;
};

export function getDatabaseTableEnvVariables(tableNames: TableNames) {
  return {
    [LambdaEnvVariables.DbStocksTableName]: tableNames.stocks,
    [LambdaEnvVariables.DbLogsTableName]: tableNames.logs,
  };
}
