package main

import (
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/providers"
)

// TESTING

// 1. poll every few seconds to see if theres anything in the SQS queue
//		take it all, and delete once its been passed?
// 2. if there is, take those things and put them into local buffers (one for each provider)
// 3. ticker for each provider (with different durationss) to see if theres anything in the buffer, and periodically ping the worker function

var testData = []jobs.JobAction{
	{Provider: providers.Fast, Type: jobs.LoadTickerHistory, TickerId: "aaa"},
	{Provider: providers.Fast, Type: jobs.LoadTickerHistory, TickerId: "bbb"},
	{Provider: providers.Slow, Type: jobs.LoadTickerHistory, TickerId: "ccc"},
	{Provider: providers.Slow, Type: jobs.LoadTickerHistory, TickerId: "ddd"},
	{Provider: providers.Fast, Type: jobs.LoadTickerHistory, TickerId: "eee"},
	{Provider: providers.Slow, Type: jobs.LoadTickerHistory, TickerId: "fff"},
	{Provider: providers.Fast, Type: jobs.LoadTickerHistory, TickerId: "ggg"},
	{Provider: providers.Fast, Type: jobs.LoadTickerHistory, TickerId: "hhh"},
	{Provider: providers.Fast, Type: jobs.LoadTickerHistory, TickerId: "iii"},
	{Provider: providers.Fast, Type: jobs.LoadTickerHistory, TickerId: "jjj"},
	{Provider: providers.Slow, Type: jobs.LoadTickerHistory, TickerId: "kkk"},
}

func test() {
	queueShort := make(chan jobs.JobAction)
	queueLong := make(chan jobs.JobAction)
}
