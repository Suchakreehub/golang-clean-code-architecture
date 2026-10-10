package gofiberserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	gofibercrudhandler "github.com/max38/golang-clean-code-architecture/src/interface/handlers/gofiber/modules/crud"
	gofibermonitorhandler "github.com/max38/golang-clean-code-architecture/src/interface/handlers/gofiber/modules/monitor"
	gofiberuserhandler "github.com/max38/golang-clean-code-architecture/src/interface/handlers/gofiber/modules/user"

	// Replace this with your own docs package
	_ "github.com/max38/golang-clean-code-architecture/src/infrastructure/gofiber/docs"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	SwaggerModule()
	CRUDModule()
}

type moduleFactory struct {
	router      fiber.Router
	server      *gofiberServer
	middlewares IMiddlewares
}

func InitModule(router fiber.Router, server *gofiberServer, middlewares IMiddlewares) IModuleFactory {
	return &moduleFactory{
		router:      router,
		server:      server,
		middlewares: middlewares,
	}
}

func (m *moduleFactory) MonitorModule() {
	var monitorHandler = gofibermonitorhandler.MonitorHandler()

	m.router.Get("/", monitorHandler.HealthCheck)
}

func (m *moduleFactory) SwaggerModule() {
	m.server.app.Get("/swagger/*", swagger.HandlerDefault) // default

	m.server.app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))
}

func (m *moduleFactory) UsersModule() {
	// var userRepository = postgresuserrepository.UserRepository()
	// var userUsecase = userusecase.UserUsecase(m.server.userRepository)
	var userHandler = gofiberuserhandler.UserHandler(m.server.application.UserUsecase)

	var router = m.router.Group("/users")

	router.Post("/register", userHandler.RegisterUser)
	router.Post("/login", userHandler.UserLogin)
	router.Post("/logout", userHandler.UserLogout)
	router.Post("/refresh-token", userHandler.RefreshToken)

	router.Post("/profile", m.middlewares.JwtAuth(m.server.application.UserUsecase), userHandler.GetUserProfile)

}

func (m *moduleFactory) CRUDModule() {
	var router = m.router.Group("/data")
	var crudHandler = gofibercrudhandler.CRUDHandler()
	router.Get("/:entity_model_slug/describe", m.middlewares.JwtAuth(m.server.application.UserUsecase), crudHandler.Describe)
	router.Get("/:entity_model_slug/", m.middlewares.JwtAuth(m.server.application.UserUsecase), crudHandler.RetriveList)
	router.Get("/:entity_model_slug/:id", m.middlewares.JwtAuth(m.server.application.UserUsecase), crudHandler.Retrive)
}
