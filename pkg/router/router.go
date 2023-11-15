package router

import (
	"FIO_App/pkg/adapters/producer"
	"FIO_App/pkg/logging"
	"FIO_App/pkg/ports/rest"
	"FIO_App/pkg/repo"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"io"
)

func SetupRoutes(storage repo.IRepository, producer producer.IProducer, logger logging.Logger) *gin.Engine {

	router := gin.Default()

	gin.DefaultWriter = io.MultiWriter(logger.Writer())
	gin.DefaultErrorWriter = io.MultiWriter(logger.Writer())

	router.GET("/echo", func(c *gin.Context) {
		c.String(200, "Check")
	})

	//graphqlHandler := handler.NewDefaultServer(
	//	graph.NewExecutableSchema(
	//		graph.Config{
	//			Resolvers: &graph.Resolver{},
	//		},
	//	),
	//)

	handlerH := rest.NewHandler(storage, producer)

	router.GET("/", playgroundHandler())
	//router.POST("/query", gin.WrapH(graphqlHandler))

	router.POST("/people", handlerH.CreatePerson)
	router.PATCH("/people/:id", handlerH.EditPerson)
	router.DELETE("/people/:id", handlerH.DeletePerson)
	router.GET("/people", handlerH.GetPeople)

	router.POST("/kafka/produce", handlerH.ProduceMessage)

	return router
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
