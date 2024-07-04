package usecase

import (
	"fmt"
	"log"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type SchedulerService struct {
	enqueuer *work.Enqueuer
}

func NewSchedulerService(redisPool *redis.Pool) *SchedulerService {
	enqueuer := work.NewEnqueuer(namespace, redisPool)
	return &SchedulerService{enqueuer: enqueuer}
}

func (qs *SchedulerService) StarScheduler(secondsInTheFuture int64, contract string) {
	fmt.Println("Starting secondsInTheFuture : ", secondsInTheFuture)
	payload := work.Q{"contract": contract}
	_, err := qs.enqueuer.EnqueueIn("scheduler_job", secondsInTheFuture, payload)
	if err != nil {
		log.Fatal(err)
	}
}
