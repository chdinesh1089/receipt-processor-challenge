package main

import (
	server "github.com/chdinesh1089/receipt-processor/server"
)

func main() {
	s := server.NewServer()
	s.Serve()
}
