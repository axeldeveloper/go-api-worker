package main

/*
import (
	"app/axel/worker/common"
	"app/axel/worker/process"
	"app/axel/worker/webui"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	},
}

var secondsInTheFuture int64 = 200

var (
	enqueuer  = work.NewEnqueuer("my_app_namespace", redisPool)
	namespace = "my_app_namespace"
)

type Tjob struct {
	customerID int64
}

func (c *Tjob) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("Starting Log job: ", job.Name)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	return next()
}

func (c *Tjob) CalculateSalary(job *work.Job) error {
	fmt.Println("Starting CalculateSalary Job: ", job.Name)
	fmt.Println("Calcalar a cada 30s. ")
	fmt.Println("CalculateSalary")
	fmt.Print((rand.Float64()*5)+5, ",")
	fmt.Println(job.Name)
	return nil
}

func (c *Tjob) SendEmail(job *work.Job) error {
	fmt.Println("Starting SendEmail Job: ", job.Name)
	addr := job.ArgString("address")
	subject := job.ArgString("subject")
	customer_id := job.ArgFloat64("customer_id")

	if err := job.ArgError(); err != nil {
		return err
	}

	fmt.Println("Address: ", addr)
	fmt.Println("Subject: ", subject)
	fmt.Println("CustomerId: ", customer_id)
	// sendEmailTo(addr, subject)
	return nil
}

func (c *Tjob) Export(job *work.Job) error {
	fmt.Println("Starting Export Job: ", job.Name)
	bucket := job.ArgString("bucket")
	fmt.Println("Upload:", bucket)
	return nil
}

func (c *Tjob) Priorities(job *work.Job) error {
	fmt.Println("Starting Priorities Job: ", job.Name)
	company := job.ArgString("company_id")
	fmt.Println("company_id:", company)
	return nil
}

func (c *Tjob) Scheduler(job *work.Job) error {
	fmt.Println("Starting Scheduler Job: ", job.Name)
	email := job.ArgString("contract")
	fmt.Println("contract:", email)
	return nil
}

func (c *Tjob) ClearCacheJob(job *work.Job) error {
	fmt.Println("Starting ClearCacheJob Job: ", job.Name)
	mesages := job.ArgString("mesages")
	fmt.Println("mesages:", mesages)
	return nil
}

func enqueueJob(job string, payload work.Q) {
	fmt.Println("Enqueued:", job, "with Paylod:", payload)
	_, err := enqueuer.Enqueue(job, payload)
	if err != nil {
		log.Fatal(err)
	}
}

func enqueueEmail() {
	var q = work.Q{"address": "axel@example.com", "subject": "hello world", "customer_id": 4}
	enqueueJob("send_email", q)
}

func enqueueS3() {
	enqueueJob("export_s3", work.Q{"bucket": "my-s3-bucket"})
}

func priorityJob() {
	enqueueJob("priority_s3", work.Q{"company_id": "portal-servidor"})
}

func schedulherJob() {
	fmt.Println("Starting secondsInTheFuture : ", secondsInTheFuture)
	_, err := enqueuer.EnqueueIn("scheduler_job", secondsInTheFuture, work.Q{"contract": "axel@patton.com"})
	if err != nil {
		log.Fatal(err)
	}
}

func uniqueJob() {
	fmt.Println("Starting uniqueJob : ")
	_, err := enqueuer.EnqueueUnique("clear_cache", work.Q{"mesages": "limpar tudo"}) // job returned
	if err != nil {
		log.Fatal(err)
	}
}

func uniqueJob2() {
	// Trabalhos exclusivos
	//secondsInTheFuture := 300
	_, err := enqueuer.EnqueueUnique("clear_cache", work.Q{"object_id_": "123"}) // job returned
	//job, err = enqueuer.EnqueueUnique("clear_cache", work.Q{"object_id_": "123"}) // job == nil -- this duplicate job isn't enqueued.
	//job, err = enqueuer.EnqueueUniqueIn("clear_cache", 300, work.Q{"object_id_": "789"}) // job != nil (diff id)
	if err != nil {
		log.Fatal(err)
	}
}

func Queueing() {
	fmt.Println("Faz um enfileirador com um namespace específico")
	pool := work.NewWorkerPool(Tjob{}, 10, namespace, redisPool)
	pool.Middleware((*Tjob).Log)

	// pool.PeriodicallyEnqueue("asterisco/30 * * * *", "calculate_salary")
	// pool.Job("calculate_salary", (*Tjob).CalculateSalary)

	// Map the name of jobs to handler functions
	// pool.Job("send_email", (*Tjob).SendEmail)
	// pool.Job("export_s3",  (*Tjob).Export)
	pool.Job("scheduler_job", (*Tjob).Scheduler)
	pool.Job("clear_cache", (*Tjob).ClearCacheJob)

	// Customize options:
	// pool.JobWithOptions("priority_s3", work.JobOptions{Priority: 10, MaxFails: 1}, (*Tjob).Priorities)

	// Start processing jobs
	pool.Start()

	webui.RunWeb()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()

}

func call4() {
	fmt.Println(common.ReverseRunes("!oG ,olleH"))
	// process.StartProcess(redisPool)
	// process.StarCron(redisPool)
	process.StarQueue(redisPool, "")
	fmt.Println("\nQuitting...")
}

func init() {
	fmt.Println("\n Timezoneeee...")
	os.Setenv("TZ", "America/Sao_Paulo")

}

func main2() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatal(err)
	}
	// handle err
	time.Local = loc

	// common.Init()
	fmt.Println("Tempo padrão:", time.Now())
	time.Local = time.UTC
	fmt.Println("Tempo UTC:", time.Now())
	fmt.Println("Tempo UTC:", time.Now())

	time.Local, _ = time.LoadLocation("America/Sao_Paulo")
	fmt.Println("Tempo SP:", time.Now())
	fmt.Println("Tempo SP:", time.Now())

	Queueing()
	//enqueueEmail()
	//enqueueS3()
	//priorityJob()
	uniqueJob()
	schedulherJob()
}

*/
