package gateway

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCServer(interceptors ...grpc.UnaryServerInterceptor) *grpc.Server {
	return grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))
}

func ServeGRPC(ctx context.Context, lis net.Listener, server *grpc.Server) error {
	errCh := make(chan error, 1)
	go func() { errCh <- server.Serve(lis) }()
	select {
	case <-ctx.Done():
		server.GracefulStop()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func NewGRPCClientConn(ctx context.Context, target string, interceptors ...grpc.UnaryClientInterceptor) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(interceptors...))
}
