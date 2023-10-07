package handlers

import (
	"FIO_App/pkg/dtos"
	"FIO_App/pkg/kafka"
	"FIO_App/pkg/storage/person"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type IHandler interface {
	CreatePerson(ctx *gin.Context)
	DeletePerson(ctx *gin.Context)
	EditPerson(ctx *gin.Context)
	GetPeople(ctx *gin.Context)
	ProduceMessage(ctx *gin.Context)
}

type Handler struct {
	storage person.IStorage
}

func NewHandler(storage person.IStorage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) CreatePerson(ctx *gin.Context) {
	var payload dtos.PersonDTO
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.storage.CreatePerson(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) DeletePerson(ctx *gin.Context) {
	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = h.storage.DeletePerson(ID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) EditPerson(ctx *gin.Context) {
	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var payload dtos.PersonDTO
	if err = ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = h.storage.EditPerson(ID, payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) GetPeople(ctx *gin.Context) {
	var limit, offset int
	limitStr, ok := ctx.GetQuery("limit")
	if ok {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid data format for the limit"})
			return
		}
		if limit < 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "limit should be greater than 0"})
			return
		}
	}
	offsetStr, ok := ctx.GetQuery("offset")
	if ok {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid data format for the offset"})
			return
		}
		if offset < 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "offset should be greater than 0"})
			return
		}
	}
	var gender, nationality string
	gender, _ = ctx.GetQuery("gender")
	nationality, _ = ctx.GetQuery("nationality")

	people, err := h.storage.GetPeople(limit, offset, nationality, gender)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"people": people})
}

func (h *Handler) ProduceMessage(ctx *gin.Context) {
	var payload kafka.FIO

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := kafka.SendMessageToQueue(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
