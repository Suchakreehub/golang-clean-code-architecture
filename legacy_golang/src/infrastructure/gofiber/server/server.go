package gofiberserver

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/max38/golang-clean-code-architecture/src/config"
	"github.com/max38/golang-clean-code-architecture/src/domain/entities"
)

type IServer interface {
	Start()
}

type gofiberServer struct {
	app         *fiber.App
	application entities.ApplicationEntity
}

func NewServer(application entities.ApplicationEntity) IServer {
	return &gofiberServer{
		app: fiber.New(
			fiber.Config{
				AppName:     config.Config.String("GOAPP_NAME"),
				JSONEncoder: json.Marshal,
				JSONDecoder: json.Unmarshal,
			},
		),
		application: application,
	}
}

func (s *gofiberServer) Start() {
	// Init Middleware
	var middlewares IMiddlewares = InitMiddlewares()
	s.app.Use(middlewares.Cors())
	s.app.Use(middlewares.Logger())

	// Init Module
	var routerGroupApiV1 = s.app.Group("/api/v1")
	var modules = InitModule(routerGroupApiV1, s, middlewares)
	modules.MonitorModule()
	modules.UsersModule()
	modules.SwaggerModule()
	modules.CRUDModule()

	s.app.Use(middlewares.RouterCheck())

	// Graceful Shutdown
	var channelSignalInterrupt = make(chan os.Signal, 1)
	signal.Notify(channelSignalInterrupt, os.Interrupt)
	go func() {
		_ = <-channelSignalInterrupt
		log.Println("server is shutting down...")
		_ = s.app.Shutdown()
	}()

	// Listen to host:port
	var serverListen string = fmt.Sprintf("%s:%d", config.Config.String("GOAPP_HOST"), config.Config.Int("GOAPP_PORT"))
	log.Printf("server is starting on %v", serverListen)
	s.app.Listen(serverListen)
}
