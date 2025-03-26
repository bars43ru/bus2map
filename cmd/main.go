package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/bars43ru/bus2map/api/bustracking"
	"github.com/bars43ru/bus2map/cmd/config"
	"github.com/bars43ru/bus2map/internal/controller"
	"github.com/bars43ru/bus2map/internal/receiver"
	"github.com/bars43ru/bus2map/internal/repository"
	"github.com/bars43ru/bus2map/internal/sender"
	"github.com/bars43ru/bus2map/internal/service"
	"github.com/bars43ru/bus2map/protocols/tcp"
	"github.com/bars43ru/bus2map/protocols/yandex"
)

type Workers interface {
	Run(ctx context.Context) error
}

type WorkerFn func(ctx context.Context) error

func (fn WorkerFn) Run(ctx context.Context) error {
	return fn(ctx)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("loading .env file", slog.Any("error", err))
		os.Exit(-1)
	}
	cfg, err := config.New()
	if err != nil {
		slog.Error("new config", slog.Any("error", err))
		os.Exit(-1)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Logger,
	}))
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	slog.InfoContext(ctx, "starting BusTracking server receiver")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		slog.InfoContext(ctx, "signal received and start main gracefully shutdown...")
		cancel()
	}()

	routeRepository := repository.NewRoute(repository.FileDatasourceRoute)
	scheduleRepository := repository.NewSchedule(repository.FileDatasourceSchedule)
	transportRepository := repository.NewTransport(repository.FileDatasourceTransport)
	busTracking := service.New(routeRepository, transportRepository, scheduleRepository)

	var workers []Workers
	workers = append(workers, routeRepository, scheduleRepository, transportRepository)

	if cfg.WialonIPS.Enabled {
		bridgeWialonIPS := receiver.BridgeWialonIPS(busTracking)

		tpcServer, err := tcp.New(cfg.WialonIPS.Addr, bridgeWialonIPS)
		if err != nil {
			slog.Error("close connection with wialon ips", slog.Any("error", err))
			return
		}
		workers = append(workers, tpcServer)
	}

	if cfg.EGTS.Enabled {
		bridgeEGTSIPS := receiver.BridgeWialonIPS(busTracking)

		tpcServer, err := tcp.New(cfg.EGTS.Addr, bridgeEGTSIPS)
		if err != nil {
			slog.Error("close connection with egts", slog.Any("error", err))
			return
		}
		workers = append(workers, tpcServer)
	}

	if cfg.Yandex.Enabled {
		cli := yandex.New(cfg.Yandex.Clid, cfg.Yandex.Url)
		worker := sender.BridgeYandex(cli, busTracking.SubscribeLocation())
		workers = append(workers, WorkerFn(worker))
	}

	if cfg.TwoGIS.Enabled {
		cli := yandex.New(cfg.TwoGIS.Clid, cfg.TwoGIS.Url)
		worker := sender.BridgeYandex(cli, busTracking.SubscribeLocation())
		workers = append(workers, WorkerFn(worker))
	}

	facade := service.New(
		routeRepository,
		transportRepository,
		scheduleRepository,
	)
	grpcSrv := grpc.NewServer()
	grpcCtrl := controller.NewBusTrackingService(facade)
	pb.RegisterBusTrackingServiceServer(grpcSrv, grpcCtrl)
	if cfg.GRPC.UseReflection {
		reflection.Register(grpcSrv)
	}
	workers = append(workers, NewGRPCSrv(grpcSrv, cfg.GRPC.ListenAddr))

	group, ctxGroup := errgroup.WithContext(ctx)
	for _, w := range workers {
		_w := w
		group.Go(func() error {
			return _w.Run(ctxGroup)
		})
	}

	err = group.Wait()
	if err != nil {
		slog.Error("end worker in BusTracking", slog.Any("error", err))
	}
	slog.InfoContext(ctx, "graceful stopped BusTracking server receiver")
}

func NewGRPCSrv(grpcSrv *grpc.Server, address string) WorkerFn {
	return func(ctx context.Context) error {
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return fmt.Errorf("open listen %s: %w", address, err)
		}
		go func() {
			<-ctx.Done()
			grpcSrv.GracefulStop()
		}()
		err = grpcSrv.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				slog.InfoContext(ctx, "grpc server has gracefully shutdown (return: %s).", slog.Any("error", err))
				return nil
			}
			slog.ErrorContext(ctx, "shutdown grpc server", slog.Any("error", err))
		}
		slog.InfoContext(ctx, "grpc server has gracefully shutdown.")
		return nil
	}
}
