package main

import (
	"log"
	"time"

	"github.com/Auxesia23/task_management/internal/database"
	"github.com/Auxesia23/task_management/internal/handlers"
	"github.com/Auxesia23/task_management/internal/repositories"
	"github.com/Auxesia23/task_management/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.NewPostgreSQLDB()
	if err != nil {
		log.Fatal("Error connecting to PostgreSQL database")
	}

	// Initialize Repositories
	userRepository := repositories.NewUserRepository(db)
	projectRepository := repositories.NewProjectRepository(db)
	invitationRepository := repositories.NewInvitationRepository(db)

	// Initialize Services
	userService := services.NewUserService(userRepository)
	projectService := services.NewProjectService(projectRepository)
	invitationService := services.NewInvitationService(invitationRepository, projectRepository)

	// Initialize Handlers
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	invitationHandler := handlers.NewInvitationHandler(invitationService)

	cfg := config{
		name:         "Task Management",
		port:         "8080",
		readTimeout:  5 * time.Second,
		writeTimeout: 5 * time.Second,
		idleTimeout:  10 * time.Second,
	}

	app := NewApplication(cfg, userHandler, projectHandler, invitationHandler)
	r := app.mount()
	log.Fatal(app.run(r))
}
