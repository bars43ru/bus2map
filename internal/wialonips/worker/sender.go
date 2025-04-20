package worker

import (
	"context"
	"fmt"
	"github.com/bars43ru/bus2map/internal/ctype"
	"github.com/bars43ru/bus2map/pkg/xslog"
	"github.com/cenkalti/backoff"
	"github.com/imkira/go-observer/v2"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"time"

	pb "github.com/bars43ru/bus2map/api/bustracking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const BackOffDelayReconnectInterval = 5 * time.Second

type Sender struct {
	tracker observer.Property[ctype.RawGPSData]
	address string
	source  string
}

func NewSender(
	tracker observer.Property[ctype.RawGPSData],
	address string,
	source string,
) *Sender {
	return &Sender{
		tracker: tracker,
		address: address,
		source:  source,
	}
}

// Run запускает клиент в цикле, отправляя данные из канала
func (s *Sender) Run(ctx context.Context) error {
	subscriber := s.tracker.Observe()
	err := backoff.RetryNotify(
		func() error {
			conn, err := grpc.NewClient(s.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return fmt.Errorf("new connection to %s: %w", s.address, err)
			}
			defer func() {
				err := conn.Close()
				if err != nil {
					slog.WarnContext(ctx, "close connection", xslog.Error(err))
				}
			}()
			if err := s.waitForConnectionReady(ctx, conn, 5*time.Second); err != nil {
				return fmt.Errorf("failure to establish connection %s: %w", s.address, err)
			}

			client := pb.NewBusTrackingServiceClient(conn)
			stream, err := client.IngestGPSData(ctx)
			if err != nil {
				return fmt.Errorf("call IngestGPSData: %w", err)
			}
			err = s.eventLoop(ctx, stream, subscriber)
			if err != nil {
				return fmt.Errorf("eventLoop: %w", err)
			}
			return nil
		},
		backoff.WithContext(backoff.NewConstantBackOff(BackOffDelayReconnectInterval), ctx),
		func(err error, duration time.Duration) {
			if err != nil {
				slog.ErrorContext(ctx, "there was a difficulty in sending the data", xslog.Error(err))
			}
		},
	)
	return err
}

func (s *Sender) eventLoop(
	ctx context.Context,
	streamingClient grpc.ClientStreamingClient[pb.GPSData, pb.IngestGPSDataResponse],
	observer observer.Stream[ctype.RawGPSData],
) error {
	var gpsData ctype.RawGPSData
	for {
		select {
		case <-ctx.Done():
			err := streamingClient.CloseSend()
			if err != nil {
				slog.ErrorContext(ctx, "close stream", xslog.Error(err))
			}
			return nil
		case <-observer.Changes():
			gpsData = observer.Next()
		}
		if gpsData.IsEmpty() {
			continue
		}
		pbGpsData := s.gpsData(gpsData)

		if err := streamingClient.Send(pbGpsData); err != nil {
			return fmt.Errorf("send GPS data: %w", err)
		}
	}
}

func (s *Sender) gpsData(gps ctype.RawGPSData) *pb.GPSData {
	return &pb.GPSData{
		Uid:       gps.UID,
		Latitude:  gps.Latitude,
		Longitude: gps.Longitude,
		Speed:     gps.Speed,
		Course:    gps.Course,
		Time:      timestamppb.New(gps.Time),
		Source:    s.source,
	}
}

func (s *Sender) waitForConnectionReady(ctx context.Context, cc *grpc.ClientConn, wait time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, wait)
	defer cancel()
	for {
		s := cc.GetState()
		if s == connectivity.Ready {
			return nil
		}

		if !cc.WaitForStateChange(ctx, s) {
			return fmt.Errorf("client connection state change timeout")
		}
	}
}
