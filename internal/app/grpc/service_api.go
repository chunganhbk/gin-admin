package grpc

// ServiceAPI allows individual grpc services to register the grpc server
type ServiceAPI interface {
	RegisterService(*RpcServer)
}
