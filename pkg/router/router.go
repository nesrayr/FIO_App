package router

import (
	"FIO_App/pkg/handlers"
	"FIO_App/pkg/storage/person"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(storage person.IStorage) *gin.Engine {
	router := gin.Default()
	router.GET("/echo", func(c *gin.Context) {
		c.String(200, "Check")
	})

	handler := handlers.NewHandler(storage)

	router.POST("/people", handler.CreatePerson)
	router.PATCH("/people/:id", handler.EditPerson)
	router.DELETE("/people/:id", handler.DeletePerson)
	router.GET("/people", handler.GetPeople)

	router.POST("/kafka/produce", handler.ProduceMessage)

	return router
}
