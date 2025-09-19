package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Config from env
	kafkaBrokers := []string{getenv("KAFKA_BROKER", "kafka:9092")}
	kafkaTopic := getenv("KAFKA_TOPIC", "orders")
	port := getenv("PORT", "8080")

	// init store + producer
	store := NewStore()
	producer := NewKafkaProducer(kafkaBrokers, kafkaTopic)
	defer producer.Close()

	app := NewApp(store, producer)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: app.Routes(),
	}

	go func() {
		log.Printf("order-service starting on :%s (Kafka: %v, topic: %s)", port, kafkaBrokers, kafkaTopic)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exiting")
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
