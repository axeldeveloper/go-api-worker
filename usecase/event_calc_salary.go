package usecase

import (
	"log"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type CalculateSalaryService struct {
	enqueuer *work.Enqueuer
}

func NewCalculateSalaryService(redisPool *redis.Pool) *CalculateSalaryService {
	enqueuer := work.NewEnqueuer(namespace, redisPool)
	return &CalculateSalaryService{enqueuer: enqueuer}
}

func (qs *CalculateSalaryService) Calc(a int64, b int64) {
	payload := work.Q{"val_1": a, "val_2": b}
	_, err := qs.enqueuer.Enqueue("calculate_salary", payload)
	if err != nil {
		log.Fatal(err)
	}
}
