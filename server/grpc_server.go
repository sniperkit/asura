package server

import (
	"net"

	"google.golang.org/grpc"

	"github.com/teragrid/asura/types"
	cmn "github.com/teragrid/teralibs/common"
)

type GRPCServer struct {
	cmn.BaseService

	proto    string
	addr     string
	listener net.Listener
	server   *grpc.Server

	app types.asuraApplicationServer
}

// NewGRPCServer returns a new gRPC asura server
func NewGRPCServer(protoAddr string, app types.asuraApplicationServer) cmn.Service {
	proto, addr := cmn.ProtocolAndAddress(protoAddr)
	s := &GRPCServer{
		proto:    proto,
		addr:     addr,
		listener: nil,
		app:      app,
	}
	s.BaseService = *cmn.NewBaseService(nil, "asuraServer", s)
	return s
}

// OnStart starts the gRPC service
func (s *GRPCServer) OnStart() error {
	if err := s.BaseService.OnStart(); err != nil {
		return err
	}
	ln, err := net.Listen(s.proto, s.addr)
	if err != nil {
		return err
	}
	s.Logger.Info("Listening", "proto", s.proto, "addr", s.addr)
	s.listener = ln
	s.server = grpc.NewServer()
	types.RegisterasuraApplicationServer(s.server, s.app)
	go s.server.Serve(s.listener)
	return nil
}

// OnStop stops the gRPC server
func (s *GRPCServer) OnStop() {
	s.BaseService.OnStop()
	s.server.Stop()
}
