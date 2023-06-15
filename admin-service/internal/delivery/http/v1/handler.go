package v1

import (
	"admin/internal/service"
	"admin/pkg/auth"
	"log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Service
	tokenManager auth.TokenManager
	responseCh   chan []byte
}

func NewHandler(service *service.Service, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     service,
		tokenManager: tokenManager,
		responseCh:   make(chan []byte),
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		go h.consumeResponseMessagesCourses()
		go h.consumeResponseMessagesStudents()
		h.initAdminRoutes(v1)
	}
}

func (h *Handler) consumeResponseMessagesCourses() {
	err := h.services.Kafka.ConsumeMessages("courses-response", h.handleResponseMessage)
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) consumeResponseMessagesStudents() {
	err := h.services.Kafka.ConsumeMessages("students-response", h.handleResponseMessage)
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) handleResponseMessage(message string) {
	h.responseCh <- []byte(message)
}
