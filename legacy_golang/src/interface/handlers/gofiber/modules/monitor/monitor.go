package gofibermonitorhandler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/max38/golang-clean-code-architecture/src/config"
	enities "github.com/max38/golang-clean-code-architecture/src/domain/entities/monitor"
	gofiberentities "github.com/max38/golang-clean-code-architecture/src/interface/handlers/gofiber/entities"
)

type IMontitorHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandler struct {
}

// HealthCheck godoc
// @Summary Show the Name and Version of server.
// @Description get the Name and Version of server.
// @Tags Monitor
// @Accept */*
// @Produce json
// @Success 200 {object} enities.Monitor
// @Router /api/v1/ [get]
func (h *monitorHandler) HealthCheck(c *fiber.Ctx) error {
	// Pass thorugh Use Case
	var responseMonitor = &enities.Monitor{
		Name:    config.Config.String("GOAPP_NAME"),
		Version: config.Config.String("GOAPP_VERSION"),
	}
	// Return response
	return gofiberentities.NewResponse(c).Success(fiber.StatusOK, responseMonitor).Response()
}

func MonitorHandler() IMontitorHandler {
	return &monitorHandler{}
}
