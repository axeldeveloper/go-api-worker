package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gocraft/web"
)

// Job representa um trabalho a ser executado
type Job struct {
	Valor string
}

// Worker representa um worker que processa jobs
type Worker struct {
	ID       int
	JobQueue chan Job
	QuitChan chan bool
}

// NewWorker cria um novo worker com um ID e uma fila de jobs
func NewWorker(id int, jobQueue chan Job) *Worker {
	return &Worker{
		ID:       id,
		JobQueue: jobQueue,
		QuitChan: make(chan bool),
	}
}

// Start inicia o worker, ele escuta por jobs no JobQueue
func (w *Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.JobQueue:
				fmt.Printf("Worker %d: processando job para o valor %s\n", w.ID, job.Valor)
				// Aqui você pode realizar o processamento do job
			case <-w.QuitChan:
				fmt.Printf("Worker %d: parando\n", w.ID)
				return
			}
		}
	}()
}

// Stop para o worker
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// JobScheduler agenda um job para ser executado após um certo tempo
func JobScheduler(job Job, delay time.Duration, jobQueue chan<- Job) {
	timer := time.NewTimer(delay)
	go func() {
		<-timer.C
		jobQueue <- job
	}()
}

var jobQueue chan Job
var workers []*Worker

func TMain() {
	jobQueue = make(chan Job, 100) // Tamanho máximo da fila de jobs

	// Inicia os workers
	numWorkers := 3 // Número de workers concorrentes
	for i := 1; i <= numWorkers; i++ {
		worker := NewWorker(i, jobQueue)
		worker.Start()
		workers = append(workers, worker)
	}

	// Configuração do roteamento HTTP com gocraft/web
	router := web.New(Worker{})
	router.Get("/agendar", AgendarHandler)

	// Inicia o servidor HTTP
	http.ListenAndServe(":8001", router)
}

// AgendarHandler manipula a requisição para agendar um job
func AgendarHandler(rw web.ResponseWriter, req *web.Request) {
	valor := req.URL.Query().Get("valor")

	if valor == "" {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Parâmetro 'valor' não fornecido\n")
		return
	}

	job := Job{Valor: valor}
	JobScheduler(job, 30*time.Second, jobQueue)

	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Job agendado com sucesso para o valor: %s\n", valor)
}
