package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kalmera74/Shorty/internal/db"
	"github.com/Kalmera74/Shorty/internal/features/analytics"
	"github.com/Kalmera74/Shorty/pkg/messaging"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(os.Stderr).With().Timestamp().Str("worker", "analytics").Logger()
	log.Logger = logger

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupSignalHandler(cancel)

	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	analyticsStore := analytics.NewAnalyticsRepository(dbConn)
	analyticsService := analytics.NewAnalyticService(analyticsStore)

	var mq messaging.IMessaging
	mq, err = messaging.NewRabbitMQConnection()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to RabbitMQ")
	}
	defer mq.Close()

	clickQueue := os.Getenv("CLICK_QUEUE")
	if clickQueue == "" {
		clickQueue = "clicks_queue"
	}

	if err := mq.DeclareQueue(clickQueue); err != nil {
		log.Fatal().Err(err).Msg("Failed to declare RabbitMQ queue")
	}

	log.Info().Str("queue", clickQueue).Msg("Worker started, waiting for messages")

	consumeLoop(ctx, mq, clickQueue, analyticsService)
}

func setupSignalHandler(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		log.Info().Msg("Shutting down analytics worker...")
		cancel()
	}()
}

func consumeLoop(ctx context.Context, mq messaging.IMessaging, queue string, service analytics.IAnalyticsService) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Context cancelled, exiting worker loop")
			return
		default:
		}

		msgs, err := mq.Consume(queue, "analytics-worker", true)
		if err != nil {
			log.Error().Err(err).Msg("Failed to start consumer, retrying in 5s")
			time.Sleep(5 * time.Second)
			continue
		}

		processMessages(ctx, msgs, service)

		log.Warn().Msg("RabbitMQ consumer channel closed, reconnecting...")
		time.Sleep(1 * time.Second)
	}
}

func processMessages(ctx context.Context, msgs <-chan messaging.IMessage, service analytics.IAnalyticsService) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Context cancelled, stopping message processing")
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Warn().Msg("Message channel closed")
				return
			}

			var click analytics.ClickEvent
			if err := json.Unmarshal(msg.Body(), &click); err != nil {
				log.Warn().Err(err).Msg("Failed to unmarshal click message")
				continue
			}

			log.Info().
				Uint("click_id", uint(click.ShortID)).
				Str("ip", click.Ip).
				Str("user_agent", click.UserAgent).
				Msg("Processing click")

			record := analytics.ClickModel{
				ShortID:   click.ShortID,
				IpAddress: click.Ip,
				UserAgent: click.UserAgent,
				CreatedAt: click.TimeStamp,
			}

			if _, err := service.Create(ctx, record); err != nil {
				log.Error().
					Err(err).
					Uint("click_id", uint(click.ShortID)).
					Msg("Failed to save click to DB")
			} else {
				_ = msg.Ack()
			}
		}
	}
}
