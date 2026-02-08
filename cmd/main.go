package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ClickHouse/clickhouse-go/v2"

	"analytics-service/internal/adapter/clickhouse"
	"analytics-service/internal/handler/kafka"
	"analytics-service/internal/service"
)

const (
	clickhouseDSN = "clickhouse://localhost:9000/default"
	kafkaBroker   = "localhost:9092"
	kafkaTopic    = "transactions.created"
	kafkaGroupID  = "analytics-service"
)

func main() {
	db, err := sql.Open("clickhouse", clickhouseDSN)
	if err != nil {
		log.Fatal("clickhouse connection error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("clickhouse ping error:", err)
	}

	log.Println("ClickHouse connected")

	repo := clickhouse.NewUserStatRepository(db)
	analyticsService := service.NewAnalyticsService(repo)
	consumer := kafka.NewConsumer(analyticsService)

	reader := kafka.NewKafkaReader(
		kafkaBroker,
		kafkaTopic,
		kafkaGroupID,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := reader.Start(ctx, consumer); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("shutting down...")
}
