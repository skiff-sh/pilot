package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/skiff-sh/pilot/pkg/behavior/behaviortype"
	"github.com/skiff-sh/pilot/server/pkg/protoenc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/skiff-sh/pilot/api/go/pilot"

	"github.com/bufbuild/protovalidate-go"
	"github.com/skiff-sh/config/contexts"
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

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			st, _ := status.FromError(err)
			if st != nil && st.Code() != codes.OK {
				httpStatus := http.StatusInternalServerError
				switch st.Code() {
				case codes.InvalidArgument:
					httpStatus = http.StatusBadRequest
				case codes.NotFound:
					httpStatus = http.StatusNotFound
				case codes.AlreadyExists:
					httpStatus = http.StatusConflict
				case codes.Unauthenticated:
					httpStatus = http.StatusUnauthorized
				}
				return ctx.Status(httpStatus).SendString(st.Message())
			}
			return fiber.DefaultErrorHandler(ctx, err)
		},
		AppName: "pilot",
		StructValidator: StructValidatorFunc(func(a any) error {
			if v, ok := a.(proto.Message); ok {
				return val.Validate(v)
			}
			return nil
		}),
		JSONEncoder: func(v any) ([]byte, error) {
			if msg, ok := v.(proto.Message); ok {
				return protoenc.ProtoMarshaller.Marshal(msg)
			} else {
				return json.Marshal(v)
			}
		},
		JSONDecoder: func(data []byte, v any) error {
			if msg, ok := v.(proto.Message); ok {
				return protoenc.ProtoUnmarshaller.Unmarshal(data, msg)
			} else {
				return json.Unmarshal(data, v)
			}
		},
	})

	store := cmap.New[behaviortype.Interface]()

	apiV1 := app.Group("/api/v1")
	if logger.Enabled(context.TODO(), slog.LevelDebug) {
		apiV1.Use(func(c fiber.Ctx) error {
			reqLogger := logger.With(
				slog.String("url", c.Request().URI().String()),
				slog.String("method", c.Method()),
			)
			reqLogger.Debug("New incoming request.")

			err := c.Next()
			if err != nil {
				reqLogger.Debug("Request failed.", slog.String("err", err.Error()))
				return err
			}

			reqLogger.Debug("Request succeeded.")
			return nil
		})
	}
	apiV1.Post("/provoke/:name", controller.NewProvokeHandler(store))

	out := &Server{
		Conf:                conf,
		Logger:              logger,
		PilotGRPCController: controller.NewPilotGRPC(val, store),
		Fiber:               app,
	}

	return out, nil
}

type Server struct {
	Conf                *config.Config
	Logger              *slog.Logger
	PilotGRPCController *controller.PilotGRPC
	Fiber               *fiber.App
}

func (s *Server) Start(ctx context.Context) error {
	logger := contexts.GetLogger(ctx)
	srv := serverapp.New(serverapp.DefaultServerInterceptors())

	grpc_health_v1.RegisterHealthServer(srv, &serverapp.HealthCheck{})
	pilot.RegisterPilotServiceServer(srv, s.PilotGRPCController)

	list, err := net.Listen("tcp", s.Conf.GRPC.Addr.String())
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error {
		err = srv.Serve(list)
		if err != nil {
			logger.Error("Failed to start gRPC handler.", slog.String("err", err.Error()))
			return err
		}

		return nil
	})

	group.Go(func() error {
		err = s.Fiber.Listen(s.Conf.HTTP.Addr.String())
		if err != nil {
			logger.Error("Failed to start HTTP handler.", slog.String("err", err.Error()))
			return err
		}
		return nil
	})

	err = group.Wait()
	if err != nil {
		logger.Error("Failed to start an application handler.", slog.String("err", err.Error()))
		return err
	}

	return nil
}

var _ fiber.StructValidator = StructValidatorFunc(nil)

type StructValidatorFunc func(a any) error

func (s StructValidatorFunc) Validate(out any) error {
	return s(out)
}
