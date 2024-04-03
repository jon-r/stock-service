export const LOGS_MODEL_NAME = "Log";

export const USERS_MODEL_NAME = "User";

export const TICKERS_MODEL_NAME = "Ticker";

export const TRANSACTIONS_MODEL_NAME = "Transaction";

export type TableNames = {
  logs: string;
  users: string;

  // OLD
  tickers: string;
  transactions: string;
  stocks: string;
};

export function getDatabaseTableEnvVariables(tableNames: TableNames) {
  return {
    DB_STOCKS_TABLE_NAME: tableNames.stocks,
    DB_LOGS_TABLE_NAME: tableNames.logs,

    // OLD
    DB_TICKERS_TABLE_NAME: tableNames.tickers,
    DB_USERS_TABLE_NAME: tableNames.users,
    DB_TRANSACTIONS_TABLE_NAME: tableNames.transactions,
  };
}
