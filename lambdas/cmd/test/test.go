package main

import (
	"fmt"
	"log"
	"time"
)

// TESTING

// 1. poll every few seconds to see if theres anything in the SQS queue
//		take it all, and delete once its been passed?
// 2. if there is, take those things and put them into local buffers (one for each provider)
// 3. ticker for each provider (with different durations) to see if theres anything in the buffer, and periodically ping the worker function

type JobAction struct {
	Provider string
	Type     string
	TickerId string
	Attempts int
}

var testData1 = []JobAction{
	{Provider: "Fast", Type: "Load", TickerId: "aaa"},
	{Provider: "Fast", Type: "Load", TickerId: "bbb"},
	{Provider: "Slow", Type: "Load", TickerId: "ccc"},
	{Provider: "Slow", Type: "Load", TickerId: "ddd"},
	{Provider: "Fast", Type: "Load", TickerId: "eee"},
	{Provider: "Slow", Type: "Load", TickerId: "fff"},
	{Provider: "Fast", Type: "Load", TickerId: "ggg"},
	{Provider: "Nope", Type: "Load", TickerId: "hhh"},
	{Provider: "Fast", Type: "Load", TickerId: "iii"},
	{Provider: "Fast", Type: "Load", TickerId: "jjj"},
	{Provider: "Slow", Type: "Load", TickerId: "kkk"},
}

var testData2 = []JobAction{}

var testData3 = []JobAction{
	{Provider: "Fast", Type: "Load", TickerId: "lll"},
	{Provider: "Fast", Type: "Load", TickerId: "mmm"},
	{Provider: "Slow", Type: "Load", TickerId: "nnn"},
	{Provider: "Slow", Type: "Load", TickerId: "ooo"},
}

var testData4 = []JobAction{
	{Provider: "Fast", Type: "Load", TickerId: "ppp"},
	{Provider: "Fast", Type: "Load", TickerId: "qqq"},
	{Provider: "Fast", Type: "Load", TickerId: "rrr"},
	{Provider: "Fast", Type: "Load", TickerId: "sss"},
}

var queueFast = make(chan JobAction, 20)
var queueSlow = make(chan JobAction, 20)
var done = make(chan bool)

func sortJobs(jobs []JobAction) {
	for _, job := range jobs {
		if job.Provider == "Fast" {
			queueFast <- job
		} else if job.Provider == "Slow" {
			queueSlow <- job
		}
	}
}

func handler(queue chan JobAction, name string, speed int) {
	slowTicker := time.NewTicker(time.Duration(speed) * time.Second)

	for {
		select {
		case <-done:
			return
		case <-slowTicker.C:
			select {
			case job, ok := <-queue:
				if ok {
					log.Printf("New %s job %+v\n", name, job)
				}
			default:
				log.Printf("No %s jobs\n", name)
			}
		}

	}
}

func test() {
	sqs := make(chan []JobAction, 6)

	sqs <- testData1
	sqs <- testData2
	sqs <- testData3
	sqs <- testData4

	queueTicker := time.NewTicker(2 * time.Second)

	go func() {
		for {
			select {
			case <-done:
				log.Println("done?")
				return
			case t := <-queueTicker.C:
				select {
				case jobs, ok := <-sqs:
					log.Printf("Tick %v\n", t.Unix())
					if ok {
						sortJobs(jobs)
					}
				default:
					log.Println("Empty Queue")
				}
			}
		}
	}()

	go handler(queueFast, "Fast", 3)
	go handler(queueSlow, "Slow", 5)

	time.Sleep(60 * time.Second)
	queueTicker.Stop()
	done <- true
	fmt.Println("DONE")
}
