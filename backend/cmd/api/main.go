package main

import (
	"log"

	todoapp "todo-backend/internal/application/todo"
	tododomain "todo-backend/internal/domain/todo"
	"todo-backend/internal/infrastructure/database"
	"todo-backend/internal/infrastructure/persistence"
)

func main() {
	// Load database configuration
	dbConfig := database.NewConfigFromEnv()
	
	// Connect to database
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()
	
	// Run database migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Verify database connectivity
	if err := database.HealthCheck(db); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}
	
	// Initialize repository
	todoRepository := persistence.NewTodoRepository(db)
	
	// Initialize domain service
	domainService := tododomain.NewDomainService(todoRepository)
	
	// Initialize use cases
	createTodoUseCase := todoapp.NewCreateTodoUseCase(todoRepository, domainService)
	getTodosUseCase := todoapp.NewGetTodosUseCase(todoRepository, domainService)
	updateTodoUseCase := todoapp.NewUpdateTodoUseCase(todoRepository, domainService)
	deleteTodoUseCase := todoapp.NewDeleteTodoUseCase(todoRepository, domainService)
	
	// TODO: Initialize handlers, middleware, and router after Interface layer is implemented
	_ = createTodoUseCase
	_ = getTodosUseCase
	_ = updateTodoUseCase
	_ = deleteTodoUseCase
	
	log.Println("Infrastructure layer successfully initialized")
	log.Println("TODO: Complete Interface layer implementation")
}
