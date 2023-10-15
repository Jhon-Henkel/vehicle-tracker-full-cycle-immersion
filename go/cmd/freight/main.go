package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Jhon-Henkel/vehicle-tracker-full-cycle-immersion/go/internal/freight/entity"
	"github.com/Jhon-Henkel/vehicle-tracker-full-cycle-immersion/go/internal/freight/infra/repository"
	"github.com/Jhon-Henkel/vehicle-tracker-full-cycle-immersion/go/internal/freight/usecase"
	"github.com/Jhon-Henkel/vehicle-tracker-full-cycle-immersion/go/pkg/kafka"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/routes?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

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
			output, err := createRouteUseCase.Execute(input)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(output)
		case "RouteStarted", "RouteFinished":
			input := usecase.ChangeRouteStatusInput{}
			json.Unmarshal(msg.Value, &input)
			output, err := changeRouteStatusUseCase.Execute(input)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(output)
		}
	}
}
