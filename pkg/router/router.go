package router

import (
	"FIO_App/pkg/adapters/producer"
	"FIO_App/pkg/handlers"
	"FIO_App/pkg/repo"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(storage repo.IRepository, producer producer.IProducer) *gin.Engine {
	router := gin.Default()
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

	handlerH := handlers.NewHandler(storage, producer)

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
