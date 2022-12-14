package micro

import "google.golang.org/grpc"

type Service interface {
	// The service name
	Name() string
	// Options returns the current options
	Options() Options
	// Run the service
	Run() error
}

type Application interface {
	Regist(grpcServer *grpc.Server)
}

func NewService(app Application, opts ...Option) Service {
	return newService(app, opts...)
}
