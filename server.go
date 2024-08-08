package main

import (
	"app/axel/worker/process"
	"app/axel/worker/usecase"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

var redisPools = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "redis:6379")
	},
}

// var pool *work.WorkerPool
var processManager *process.ProcessManager

var (
	ErrInvalidEvent    = errors.New("invalid event data")
	ErrEventFull       = errors.New("event is full")
	ErrTicketNotFound  = errors.New("ticket not found")
	ErrTicketNotEnough = errors.New("not enough tickets available")
	ErrEventNotFound   = errors.New("event not found")
	ErrNotArgsFound    = errors.New("arguments not found")
	NamespaceWork      = "work_namespace"
)

type Event struct {
	ID           string
	Name         string
	Location     string
	Organization string
	Date         time.Time
	ImageURL     string
	Capacity     int
	Price        float64
	PartnerID    int
}

func SomeHandler(w http.ResponseWriter, r *http.Request) {
	data := Event{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func EventHundler(w http.ResponseWriter, r *http.Request) {
	// id := r.URL.Path[len("/event/"):]
	vars := mux.Vars(r)
	id := vars["id"]
	eve := Event{ID: id}

	queueService := usecase.NewQueueService(redisPools)
	queueService.StarQueue(id)
	queueService.SendEmailIn("patton@gmail.com")

	schedulerService := usecase.NewSchedulerService(redisPools)
	schedulerService.StarScheduler(200, "AXELPATTO@GMAIL.COM")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(eve)
}

func CalcHundler(w http.ResponseWriter, r *http.Request) {
	// id := r.URL.Path[len("/event/"):]
	vars := mux.Vars(r)
	id := vars["id"]
	eve := Event{ID: id}
	num, _ := strconv.ParseInt(id, 10, 64)

	queueService := usecase.NewCalculateSalaryService(redisPools)
	queueService.Calc(num, 10)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(eve)
}

func init() {
	processManager = process.NewProcessManager(redisPools, 10, NamespaceWork)
	//processManager.StartSignal()
}

func main() {

	fmt.Println("Iniciando servidor web")

	r := mux.NewRouter()
	r.HandleFunc("/index", SomeHandler)
	r.HandleFunc("/event/{id}", EventHundler)
	r.HandleFunc("/calc/{id}", CalcHundler)

	server := &http.Server{
		Addr:    ":8001",
		Handler: r,
	}

	// Canal para escutar sinais do sistema operacional
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		// Recebido sinal de interrupção, iniciando o graceful shutdown
		log.Println("Recebido sinal de interrupção, iniciando o graceful shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Erro no graceful shutdown: %v\n", err)
		}
		//pool.Stop()
		processManager.Stop()

		close(idleConnsClosed)
	}()
	// Iniciando o servidor HTTP
	log.Println("Servidor HTTP rodando na porta 8001")

	//pool.Start()
	processManager.Start()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Erro ao iniciar o servidor HTTP: %v\n", err)
	}

	<-idleConnsClosed
	log.Println("Servidor HTTP finalizado")

}
