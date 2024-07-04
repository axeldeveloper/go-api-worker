package process

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var (
	namespace = "work_namespace"
)

type Context struct {
	customerID int64
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job:", job.Name)
	return next()
}

func (c *Context) FindCustomer(job *work.Job, next work.NextMiddlewareFunc) error {
	// If there's a customer_id param, set it in the context for future middleware and handlers to use.
	if _, ok := job.Args["customer_id"]; ok {
		c.customerID = job.ArgInt64("customer_id")
		if err := job.ArgError(); err != nil {
			return err
		}
	}

	return next()
}

func (c *Context) SendEmail(job *work.Job) error {
	// Extract arguments:
	addr := job.ArgString("address")
	subject := job.ArgString("subject")
	if err := job.ArgError(); err != nil {
		return err
	}

	// Go ahead and send the email...
	fmt.Println("Send email to:", addr, "with subject:", subject, "and customer id:", c.customerID)

	return nil
}

func (c *Context) Scheduler(job *work.Job) error {
	contract := job.ArgString("contract")
	if err := job.ArgError(); err != nil {
		return err
	}

	// Go ahead and send the email...
	fmt.Println("Scheduler Send email to:", contract)

	return nil
}

func NewProcessManager(redisPool *redis.Pool, numWorkers uint, namespace string) *ProcessManager {
	ctx := Context{} // Substitua com o contexto que você está usando
	pool := work.NewWorkerPool(ctx, numWorkers, namespace, redisPool)
	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).FindCustomer)
	pool.Job("send_email", (*Context).SendEmail)
	pool.Job("scheduler_job", (*Context).Scheduler)
	return &ProcessManager{pool: pool}
}

type ProcessManager struct {
	pool *work.WorkerPool
}

func (pm *ProcessManager) Start() {
	pm.pool.Start()
}

func (pm *ProcessManager) StartSignal() {
	pm.pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	pm.Stop()
}

func (pm *ProcessManager) Stop() {
	pm.pool.Stop()
}

/*
func StartProcess(redisPool *redis.Pool) {
	pool := work.NewWorkerPool(Context{}, 10, namespace, redisPool)
	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).FindCustomer)
	// Map the name of jobs to handler functions
	pool.Job("send_email", (*Context).SendEmail)

	// Customize options:
	// pool.JobWithOptions("upload_s3", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).Export)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func StarCron(redisPool *redis.Pool) {
	fmt.Println("RUN JOB CROM")

	// ok
	pool := work.NewWorkerPool(Context{}, 10, namespace, redisPool)
	pool.PeriodicallyEnqueue("1 * * * * *", "get_repository")

	fmt.Println("CROM DE 1 MINUTO")
	pool.Job("get_repository", (*Context).GetRepository)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
	fmt.Println("\nQUITTING RUN CRON ...")
}

*/
