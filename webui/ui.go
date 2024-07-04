package webui

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gocraft/work/webui"
	"github.com/gomodule/redigo/redis"
)

var (
	//redisHostPort  = flag.String("redis", ":6379", "tcp")
	//redisDatabase  = flag.String("database", "0", "my_app_namespace")
	//redisNamespace = flag.String("ns", "work", "my_app_namespace")
	//webHostPort    = flag.String("listen", ":5040", "localhost")

	redisHostPort  = flag.String("redis", ":6379", "tcp")
	redisDatabase  = flag.String("database", "0", "db0")
	redisNamespace = flag.String("ns", "my_app_namespace", "my_app_namespace")
	webHostPort    = flag.String("listen", ":5040", "localhost")
)

func RunWeb() {
	flag.Parse()
	// https://github.com/stanislas-m/amqp-work-adapter
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("Starting workwebui:")
	fmt.Println("redis = ", *redisHostPort)
	fmt.Println("database = ", *redisDatabase)
	fmt.Println("namespace = ", *redisNamespace)
	fmt.Println("listen = ", *webHostPort)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	database, err := strconv.Atoi(*redisDatabase)
	if err != nil {
		fmt.Printf("Error: %v is not a valid database value", *redisDatabase)
		return
	}

	pool := newPool(*redisHostPort, database)

	server := webui.NewServer(*redisNamespace, pool, *webHostPort)
	server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	server.Stop()

	fmt.Println("\n Quitting...")
}

func newPool(addr string, database int) *redis.Pool {

	fmt.Println("\npooll..")

	fmt.Println(addr)
	fmt.Println(database)

	return &redis.Pool{
		MaxActive:   3,
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			//return redis.DialURL(addr, redis.DialDatabase(database))
			return redis.Dial("tcp", ":6379")
		},
		Wait: true,
	}
}
