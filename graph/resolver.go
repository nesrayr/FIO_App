package graph

import (
	"FIO_App/pkg/repo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	storage repo.IRepository
}
