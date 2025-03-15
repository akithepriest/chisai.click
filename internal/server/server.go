package server

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/akithepriest/chisai.click/internal"
	"github.com/akithepriest/chisai.click/internal/handlers"
	"github.com/akithepriest/chisai.click/internal/middlewares"
	"github.com/labstack/echo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EchoServer struct {
	e *echo.Echo
}

// Makes the server go BROOM
func NewEchoServer() *EchoServer {
	e := echo.New()
	e.Static("/static", "static")
	e.Use(middlewares.LoggingMiddleware)

	db, err := InitDB()
	if err != nil {
		internal.Logger.LogError().Msg(err.Error())
	}

	handlers.BindHandlers(e, db)
	return &EchoServer{
		e: e,
	}
}

func (s *EchoServer) Start() {
	internal.Logger.LogError().Msg(s.e.Start(":8080").Error())
}

func InitDB() (*mongo.Database, error) {
	connString := os.Getenv("MONGODB_CONN_STRING")
	if connString == "" {
		return nil, errors.New("MONGODB_CONN_STRING is not defined")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connString))
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection with mongo server: %v", err)
	}

	go func() {
		<-context.Background().Done()
		client.Disconnect(context.Background())
	}()

	db := client.Database("master")

	return db, nil
}