package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/Auxesia23/task_management/internal/dto"
	"github.com/Auxesia23/task_management/internal/handlers"
	"github.com/Auxesia23/task_management/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/valyala/fasthttp"
)

type application struct {
	cfg            config
	userHandler    handlers.UserHandler
	projectHandler handlers.ProjectHandler
}

type config struct {
	name         string
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

func NewApplication(cfg config, userHandler handlers.UserHandler, projectHandler handlers.ProjectHandler) *application {
	return &application{
		cfg,
		userHandler,
		projectHandler,
	}
}

func (app *application) mount() fasthttp.RequestHandler {
	r := fiber.New()

	logFile, err := os.OpenFile("fiber.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	output := io.MultiWriter(os.Stdout, logFile)

	r.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${ip}) - ${latency} - ${ua}\n",
		Output: output,
	}))
	r.Use(healthcheck.New())

	r.Get("/metrics", monitor.New())
	r.Get("/protected", middlewares.JWTAuthMiddleware, func(c *fiber.Ctx) error {
		user := c.Locals("user").(*dto.AccessTokenClaims)
		return c.SendString("Hello, " + user.FullName)
	})

	auth := r.Group("/auth")
	{
		auth.Post("/register", app.userHandler.RegisterHandler)
		auth.Post("/login", app.userHandler.LoginHandler)
		auth.Post("/refresh", app.userHandler.RefreshHandler)

	}

	users := r.Group("/users", middlewares.JWTAuthMiddleware)
	{
		users.Get("/", app.userHandler.SearchUserhandler)
	}

	projects := r.Group("/projects")
	{
		projects.Post("/", middlewares.JWTAuthMiddleware, app.projectHandler.CreateProjectHanlder)
		projects.Get("/", app.projectHandler.GetProjectsHanlder)
		projects.Get("/:id", app.projectHandler.ReadProjectByIdHanlder)
		projects.Put("/:id", middlewares.JWTAuthMiddleware, app.projectHandler.UpdateProjectHanlder)
		projects.Delete("/:id", middlewares.JWTAuthMiddleware, app.projectHandler.DeleteProjectHanlder)
	}

	return r.Handler()
}

func (app *application) run(r fasthttp.RequestHandler) error {
	srv := &fasthttp.Server{
		Handler:      r,
		Name:         app.cfg.name,
		ReadTimeout:  app.cfg.readTimeout,
		WriteTimeout: app.cfg.writeTimeout,
		IdleTimeout:  app.cfg.idleTimeout,
	}
	log.Printf("Starting %s on port %s", app.cfg.name, app.cfg.port)
	return srv.ListenAndServe(":" + app.cfg.port)
}
