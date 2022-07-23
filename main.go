package main

import (
	"log"
	"os"

	"github.com/romandkv/andrey/pkg/server"
)

const (
	port = "port"
)

func main() {
	svr := server.NewServer(
		os.Getenv(port),
		&server.JsonValueMutatorMiddleware{},
		&server.JsonKeyMutatorMiddleware{},
		&server.Marshaler{},
		server.FinalHandler{},
	)
	log.Fatalln(svr.Run())
}
