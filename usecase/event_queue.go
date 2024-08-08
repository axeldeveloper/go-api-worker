package usecase

import (
	"log"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var (
	namespace = "work_namespace"
)

type QueueService struct {
	enqueuer *work.Enqueuer
}

func NewQueueService(redisPool *redis.Pool) *QueueService {
	enqueuer := work.NewEnqueuer(namespace, redisPool)
	return &QueueService{enqueuer: enqueuer}
}

func (qs *QueueService) StarQueue(email string) {
	payload := work.Q{"address": email, "subject": "hello world", "customer_id": 4}
	_, err := qs.enqueuer.Enqueue("send_email", payload)
	if err != nil {
		log.Fatal(err)
	}
}

func (qs *QueueService) SendEmailIn(email string) {
	payload := work.Q{"email": email, "subject": "hello world", "customer_id": 4}
	_, err := qs.enqueuer.EnqueueIn("send_welcome_email", 300, payload)

	if err != nil {
		log.Fatal(err)
	}
}
