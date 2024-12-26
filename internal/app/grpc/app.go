package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	authgrpc "github.com/paveldroo/sso-service/internal/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService authgrpc.Auth, port int) App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
	}
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))
			return status.Error(codes.Internal, "internal error")
		}),
	}
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor(recoveryOpts...), logging.UnaryServerInterceptor(InterceptorLogger(slog.Default()), loggingOpts...)))

	authgrpc.Register(srv, authService)

	return App{
		log:        log,
		grpcServer: srv,
		port:       port,
	}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(
		func(ctx context.Context, level logging.Level, msg string, fields ...any) {
			l.Log(ctx, slog.Level(level), msg, fields...)
		})
}

func (a App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("start listen tcp: %w", err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("serving grpc server: %w", err)
	}

	return nil
}

func (a App) Stop() {
	a.log.Info("stopping grpc server")
	a.grpcServer.GracefulStop()
}
