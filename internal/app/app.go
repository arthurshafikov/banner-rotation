package app

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arthurshafikov/banner-rotation/internal/config"
	"github.com/arthurshafikov/banner-rotation/internal/repository"
	"github.com/arthurshafikov/banner-rotation/internal/services"
	"github.com/arthurshafikov/banner-rotation/internal/transport/http"
	"github.com/arthurshafikov/banner-rotation/internal/transport/http/handler"
	"github.com/arthurshafikov/banner-rotation/pkg/postgres"
	"github.com/arthurshafikov/banner-rotation/pkg/queue"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

var envFileLocation string

func init() {
	flag.StringVar(&envFileLocation, "env", "./.env", "Path to .env file")
}

func Run() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	group, ctx := errgroup.WithContext(ctx)
	defer cancel()

	config := config.NewConfig(envFileLocation)

	db := postgres.NewSqlxDB(ctx, group, config.DSN)

	repos := repository.NewRepository(db)

	loggerKafka := log.New(os.Stdout, "kafka writer: ", 0)
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{config.QueueConfig.BrokerAddress},
		Logger:  loggerKafka,
	})
	queue := queue.NewQueue(ctx, kafkaWriter)
	group.Go(func() error {
		return queue.Dispatch()
	})

	services := services.NewServices(services.Dependencies{
		Repository:  repos,
		EGreedValue: config.MultihandedBanditConfig.EGreedValue,
		Queue:       queue,
	})

	handler := handler.NewHandler(ctx, services, http.NewRequestParser())
	s := http.NewServer(handler)
	s.Serve(ctx, group, config.ServerConfig.Port)

	if err := group.Wait(); err != nil {
		log.Println(err)
	}
}
