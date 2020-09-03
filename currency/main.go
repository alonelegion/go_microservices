package main

import (
	protos "github.com/alonelegion/go_microservices/currency/protos/currency"
	"github.com/alonelegion/go_microservices/currency/server"
	"google.golang.org/grpc/reflection"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	log := hclog.Default()

	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()
	cs := server.NewCurrency(log)

	// register the currency server
	protos.RegisterCurrencyServer(gs, cs)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(gs)

	// create a TCP socket for inbound server connections
	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	// listen for requests
	_ = gs.Serve(l)
}
