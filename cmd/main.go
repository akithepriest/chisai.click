package main

import (
	"github.com/akithepriest/chisai.click/internal"
	"github.com/akithepriest/chisai.click/internal/server"
)

func main() {
	internal.NewLogger()
	srv := server.NewEchoServer()
	srv.Start()
}