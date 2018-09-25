package main

import (
	"fmt"
	"hlcup/services"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
	"github.com/buaazp/fasthttprouter"
	"hlcup/handlers"
	"github.com/valyala/fasthttp"
	"github.com/paulbellamy/ratecounter"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	counter := ratecounter.NewRateCounter(5 * time.Second)

	errorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
	infoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	storage := services.NewStorage(errorLogger, infoLogger)

	waitGroup := new(sync.WaitGroup)

	startTime := time.Now()

	storage.Init("/home/artyomnorin/projects/go/src/hlcup/data/full/data.zip", 4, waitGroup)

	waitGroup.Wait()

	infoLogger.Println(fmt.Sprintf("index has been built. Duration: %f", time.Now().Sub(startTime).Seconds()))

	PrintMemUsage()
	runtime.GC()
	PrintMemUsage()

	userApiHandler := handlers.NewUserApiHandler(storage, errorLogger, infoLogger, counter, "/home/artyomnorin/projects/go/src/hlcup/data/full/options.txt")
	locationApiHandler := handlers.NewLocationApiHandler(storage, errorLogger, infoLogger, counter,  "/home/artyomnorin/projects/go/src/hlcup/data/full/options.txt")
	visitApiHandler := handlers.NewVisitApiHandler(storage, errorLogger, infoLogger, counter)

	router := fasthttprouter.New()

	router.GET("/users/:user_id", userApiHandler.GetById)
	router.GET("/locations/:location_id", locationApiHandler.GetById)
	router.GET("/visits/:visit_id", visitApiHandler.GetById)

	router.GET("/users/:user_id/visits", userApiHandler.GetVisitedPlaces)
	router.GET("/locations/:location_id/avg", locationApiHandler.GetAverageMark)

	router.POST("/users/:user_id", userApiHandler.CreateOrUpdate)
	router.POST("/locations/:location_id", locationApiHandler.CreateOrUpdate)
	router.POST("/visits/:visit_id", visitApiHandler.CreateOrUpdate)

	go func() {
		for {
			infoLogger.Println(fmt.Sprintf("RPS: %d", counter.Rate()))
			time.Sleep(5 * time.Second)
		}
	}()

	infoLogger.Println("Server is listening...")
	err := fasthttp.ListenAndServe(":8080", router.Handler)

	if err != nil {
		errorLogger.Fatalln(err)
	}
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
