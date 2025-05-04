package main

import (
	"context"
	"github.com/bars43ru/bus2map/internal/ctype"
	config2 "github.com/bars43ru/bus2map/internal/wialonips/config"
	service2 "github.com/bars43ru/bus2map/internal/wialonips/service"
	"github.com/bars43ru/bus2map/internal/wialonips/worker"
	"github.com/bars43ru/bus2map/pkg/xslog"
	"github.com/imkira/go-observer/v2"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bars43ru/bus2map/pkg/tcp"
)

type Worker interface {
	Run(ctx context.Context) error
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("loading .env file", xslog.Error(err))
		os.Exit(-1)
	}
	cfg, err := config2.New()
	if err != nil {
		slog.Error("new config", xslog.Error(err))
		os.Exit(-1)
	}
	setupLogger(cfg.Logger)

	ctx, cancel := context.WithCancel(context.Background())
	slog.InfoContext(ctx, "starting BusTracking server receiver")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		slog.InfoContext(ctx, "signal received and start main gracefully shutdown...")
		cancel()
	}()

	tracker := observer.NewProperty[ctype.RawGPSData](ctype.RawGPSData{})
	handlerWialonIPSStream := service2.NewReciver(tracker)

	var workers []Worker
	workers = append(workers, worker.NewSender(tracker, cfg.CoordinateAddr, cfg.Source))

	tpcServer, err := tcp.New(cfg.WialonIPS.Addr, handlerWialonIPSStream)
	if err != nil {
		slog.Error("close connection with wialon ips", xslog.Error(err))
		return
	}
	workers = append(workers, tpcServer)

	group, ctxGroup := errgroup.WithContext(ctx)
	for _, w := range workers {
		_w := w
		group.Go(func() error {
			return _w.Run(ctxGroup)
		})
	}
	err = group.Wait()
	if err != nil {
		slog.Error("end worker in WialonIPS", xslog.Error(err))
	}
	slog.InfoContext(ctx, "graceful stopped WialonIPS server receiver")
}

func setupLogger(cfg config2.Logger) {
	handlers := []slog.Handler{
		slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.Level}),
	}
	//if cfg.ToFile {
	//	logWriter := &lumberjack.Logger{
	//		Filename: "./logs/current.log",
	//		MaxSize:  10,
	//		MaxAge:   30,
	//		Compress: true,
	//	}
	//	handlers = append(handlers, slog.NewTextHandler(logWriter, &slog.HandlerOptions{Level: cfg.Level}))
	//}
	logger := slog.New(xslog.NewMultiHandler(handlers...))
	slog.SetDefault(logger)
}
