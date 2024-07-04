package main

/*
import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

func main() {
	jobQueue = make(chan Job, 100) // Tamanho máximo da fila de jobs

	// Inicia os workers
	numWorkers := 3 // Número de workers concorrentes
	for i := 1; i <= numWorkers; i++ {
		worker := NewWorker(i, jobQueue)
		worker.Start()
		workers = append(workers, worker)
	}

	// Configuração do roteamento HTTP com Gorilla Mux
	r := mux.NewRouter()
	r.HandleFunc("/agendar", AgendarHandler).Methods("GET")

	// Inicia o servidor HTTP
	http.ListenAndServe(":8001", r)
}

// AgendarHandler manipula a requisição para agendar um job
func AgendarHandler(w http.ResponseWriter, r *http.Request) {
	valor := r.URL.Query().Get("valor")

	if valor == "" {
		http.Error(w, "Parâmetro 'valor' não fornecido", http.StatusBadRequest)
		return
	}

	job := Job{Valor: valor}
	JobScheduler(job, 30*time.Second, jobQueue)

	fmt.Fprintf(w, "Job agendado com sucesso para o valor: %s\n", valor)
}
*/
