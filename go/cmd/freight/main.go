package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jhon-Henkel/vehicle-tracker-full-cycle-immersion/go/internal/freight/entity"
	"github.com/Jhon-Henkel/vehicle-tracker-full-cycle-immersion/go/internal/freight/infra/repository"
	"github.com/Jhon-Henkel/vehicle-tracker-full-cycle-immersion/go/internal/freight/usecase"
	"github.com/Jhon-Henkel/vehicle-tracker-full-cycle-immersion/go/pkg/kafka"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

)

var (
	routesCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "routes_created_total",
			Help: "Total number of routes created",
		},
	)
	routesStarted = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "routes_started_total",
			Help: "Total number of routes started",
		},
		[]string{"status"},
	)
	errorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
	)
)

func init() {
	prometheus.MustRegister(routesCreated)
	prometheus.MustRegister(routesStarted)
	prometheus.MustRegister(errorsTotal)
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/routes?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":8080", nil)
	} ()

	msgChan := make(chan *ckafka.Message)
	topics := []string{"route"}
	servers := "host.docker.internal:9094"
	go kafka.Consumer(topics, servers, msgChan)

	repository := repository.NewRouteRepositoryMysql(db)
	freight := entity.NewFreight(10)
	createRouteUseCase := usecase.NewCreateRouteUseCase(repository, freight)
	changeRouteStatusUseCase := usecase.NewChangeRouteStatusUseCase(repository)

	for msg := range msgChan {
		input := usecase.CreateRouteInput{}
		json.Unmarshal(msg.Value, &input)

		switch input.Event {
		case "RouteCreated":
			_, err := createRouteUseCase.Execute(input)
			if err != nil {
				fmt.Println(err)
				errorsTotal.Inc()
			} else {
				routesCreated.Inc()
			}
		case "RouteStarted":
			input := usecase.ChangeRouteStatusInput{}
			json.Unmarshal(msg.Value, &input)
			_, err := changeRouteStatusUseCase.Execute(input)
			if err != nil {
				fmt.Println(err)
				errorsTotal.Inc()
			} else {
				routesStarted.WithLabelValues("started").Inc()
			}
		case "RouteFinished":
			input := usecase.ChangeRouteStatusInput{}
			json.Unmarshal(msg.Value, &input)
			_, err := changeRouteStatusUseCase.Execute(input)
			if err != nil {
				fmt.Println(err)
				errorsTotal.Inc()
			} else {
				routesStarted.WithLabelValues("finished").Inc()
			}
		}
	}
}
