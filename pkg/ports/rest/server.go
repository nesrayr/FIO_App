package rest

import (
	"FIO_App/pkg/adapters/producer"
	"FIO_App/pkg/logging"
	"FIO_App/pkg/repo"
	"net/http"
)

func NewRESTServer(addr string, repo repo.IRepository, producer producer.IProducer, logger logging.Logger) *http.Server {
	r := SetupRoutes(repo, producer, logger)
	return &http.Server{Addr: addr, Handler: r}
}
