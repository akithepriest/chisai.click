package handlers

import (
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
)

func BindHandlers(e *echo.Echo, db *mongo.Database) {
	indexGroup := e.Group("")
	indexHandler := NewIndexHandler(db)
	indexHandler.DefineRoutes(indexGroup)
}