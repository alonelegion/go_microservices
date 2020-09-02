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

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	gs.Serve(l)
}
