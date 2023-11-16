package rest

import (
	"FIO_App/pkg/adapters/producer"
	"FIO_App/pkg/logging"
	"FIO_App/pkg/repo"
	"github.com/gin-gonic/gin"
	"io"
)

func SetupRoutes(repository repo.IRepository, producer producer.IProducer, logger logging.Logger) *gin.Engine {

	router := gin.Default()

	gin.DefaultWriter = io.MultiWriter(logger.Writer())
	gin.DefaultErrorWriter = io.MultiWriter(logger.Writer())

	router.GET("/echo", func(c *gin.Context) {
		c.String(200, "Check")
	})

	handlerH := NewHandler(repository, producer)

	router.POST("/people", handlerH.CreatePerson)
	router.PATCH("/people/:id", handlerH.EditPerson)
	router.DELETE("/people/:id", handlerH.DeletePerson)
	router.GET("/people", handlerH.GetPeople)

	router.POST("/kafka/produce", handlerH.ProduceMessage)

	return router
}
