package main

import (
	"encoding/json"
	"log"

	"github.com/Kalmera74/Shorty/internal/db"
	"github.com/Kalmera74/Shorty/internal/features/analytics"
	"github.com/Kalmera74/Shorty/pkg/messaging"
)

func main() {
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	analyticsStore := analytics.NewAnalyticsRepository(dbConn)
	analyticsService := analytics.NewAnalyticService(analyticsStore)

	rabbit, err := messaging.NewRabbitMQConnection()
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer rabbit.Close()

	err = rabbit.DeclareQueue("clicks")
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	log.Printf("Worker started, waiting for messages on queue: clicks")

	msgs, err := rabbit.Consume("clicks", "analytics-worker", true)
	if err != nil {
		log.Fatalf("failed to start consumer: %v", err)
	}

	for msg := range msgs {
		var click analytics.ClickEvent
		if err := json.Unmarshal(msg.Body, &click); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		log.Printf("Processing click for short_id=%s", click.ShortID)

		record := analytics.ClickModel{
			ShortID:   click.ShortID,
			IpAddress: click.Ip,
			UserAgent: click.UserAgent,
			CreatedAt: click.TimeStamp,
		}

		if _, err := analyticsService.Create(record); err != nil {
			log.Printf("failed to save click: %v", err)
			continue
		}

	}
}
