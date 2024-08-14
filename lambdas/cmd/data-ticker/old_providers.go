package main

//
//var providerQueues = map[providers_old.ProviderName]chan jobs_old.JobQueueItem{
//	providers_old.PolygonIo: make(chan jobs_old.JobQueueItem, 20),
//}
//
//func allocateJobs(jobList *[]jobs_old.JobQueueItem) {
//	for _, job := range *jobList {
//		providerQueues[job.Action.Provider] <- job
//	}
//}
//
//func (handler DataTickerHandler) invokeWorkerTicker(provider providers_old.ProviderName, delay providers_old.SettingsDelay) {
//	var err error
//
//	duration := time.Duration(delay) * time.Second
//	ticker := handler.Clock.Ticker(duration)
//
//	for {
//		select {
//		case <-ticker.C:
//			select {
//			case job, ok := <-providerQueues[provider]:
//				if ok {
//					handler.LogService.Infow("Invoking Job",
//						"job", job,
//					)
//					err = handler.EventsService.InvokeWorker(job.Action)
//					if err != nil {
//						handler.LogService.Warnw("Failed to Invoke Worker",
//							"error", err,
//						)
//
//						err = handler.QueueService.RetryJob(job.Action, err.Error(), handler.NewUuid)
//					}
//
//					err = handler.QueueService.DeleteJob(job.RecieptHandle)
//					if err != nil {
//						handler.LogService.Warnw("Failed to delete Job from queue",
//							"error", err,
//						)
//					}
//				}
//			default:
//				// no jobs_old for this provider
//			}
//		}
//	}
//}
