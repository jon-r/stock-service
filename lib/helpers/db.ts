export const LOGS_MODEL_NAME = "Log";

export const USERS_MODEL_NAME = "User";

export const TICKERS_MODEL_NAME = "Ticker";

export const TRANSACTIONS_MODEL_NAME = "Transaction";

export type TableNames = {
  logs: string;
  stocks: string;
};

export function getDatabaseTableEnvVariables(tableNames: TableNames) {
  return {
    DB_STOCKS_TABLE_NAME: tableNames.stocks,
    DB_LOGS_TABLE_NAME: tableNames.logs,
  };
}
