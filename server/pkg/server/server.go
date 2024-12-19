package server

import (
	"context"
	"log/slog"
	"net"

	"github.com/bufbuild/protovalidate-go"
	"github.com/skiff-sh/config/contexts"
	pilot "github.com/skiff-sh/pilot/api/go"
	"github.com/skiff-sh/pilot/server/pkg/config"
	"github.com/skiff-sh/pilot/server/pkg/controller"
	"github.com/skiff-sh/serverapp"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func New(conf *config.Config, logger *slog.Logger) (*Server, error) {
	val, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	out := &Server{
		Conf:            conf,
		Logger:          logger,
		PilotController: controller.NewController(val),
	}

	return out, nil
}

type Server struct {
	Conf            *config.Config
	Logger          *slog.Logger
	PilotController *controller.Controller
}

func (s *Server) Start(ctx context.Context) error {
	logger := contexts.GetLogger(ctx)
	srv := serverapp.New(serverapp.DefaultServerInterceptors())

	grpc_health_v1.RegisterHealthServer(srv, &serverapp.HealthCheck{})
	pilot.RegisterPilotServiceServer(srv, s.PilotController)

	list, err := net.Listen("tcp", s.Conf.Server.Addr.String())
	if err != nil {
		return err
	}

	err = srv.Serve(list)
	if err != nil {
		logger.Error("Failed to serve.", slog.String("err", err.Error()))
		return err
	}

	return nil
}
