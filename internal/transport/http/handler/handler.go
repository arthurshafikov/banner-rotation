package handler

import "github.com/thewolf27/banner-rotation/internal/services"

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}
