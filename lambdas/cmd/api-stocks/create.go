package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"

	"jon-richards.com/stock-app/db"
	"jon-richards.com/stock-app/providers"
)

//type response struct {
//	Message string                        `json:"greeting"`
//	Event   events.APIGatewayProxyRequest `json:"request"`
//}

type RequestParams struct {
	Provider providers.ProviderName `json:"provider"`
	Ticker   string                 `json:"ticker"`
}

var dbService = db.NewDatabaseService()

func createStockIndex(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var err error

	// 1. get ticket and provider from post request body
	var params RequestParams
	err = json.Unmarshal([]byte(request.Body), &params)

	if err != nil {
		return nil, err
	}

	// 2. fetch the ticket details (based on the above)
	err, details := providers.GetStockIndexDetails(params.Provider, params.Ticker)

	if err != nil {
		return nil, err
	}

	// 3. insert this ^ data into the stockIndex table

	// 4. insert a 'prepopulate history' action to the jobs table

	// 5. return 'ok'
}

func create(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// todo would be switch if multiple endpoints
	return createStockIndex(request)
}
