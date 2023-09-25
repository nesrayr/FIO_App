package router

import (
	"FIO_App/graph"
	"FIO_App/pkg/handlers"
	"FIO_App/pkg/storage/person"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(storage person.IStorage) *gin.Engine {
	router := gin.Default()
	router.GET("/echo", func(c *gin.Context) {
		c.String(200, "Check")
	})

	graphqlHandler := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{},
			},
		),
	)

	handlerH := handlers.NewHandler(storage)

	router.GET("/", playgroundHandler())
	router.POST("/query", gin.WrapH(graphqlHandler))

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
