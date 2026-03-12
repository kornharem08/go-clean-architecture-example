package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/microsoft/go-mssqldb"

	addressRepo "github.com/example/clean-architecture/infra/mssql/address"
	userRepo "github.com/example/clean-architecture/infra/mssql/user"
	addressSvc "github.com/example/clean-architecture/internal/address"
	userSvc "github.com/example/clean-architecture/internal/user"
	"github.com/example/clean-architecture/routers"
)

func main() {
	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize User repository and usecase
	userRepo := userRepo.NewRepository(db)
	userUsecase := userSvc.NewService(userRepo)

	// Initialize Address repository and usecase
	addressRepo := addressRepo.NewAddressRepository(db)
	addressUsecase := addressSvc.NewService(addressRepo)

	// Create Gin router
	router := gin.Default()

	// Setup routes: inject all usecases (addressRepo no longer passed directly)
	routers.SetupRoutes(router, userUsecase, addressUsecase)

	// Start server
	port := getPort()
	log.Printf("Server running on port %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDB initializes the database connection
func initDB() (*sql.DB, error) {
	// Get connection string from environment or use default
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Default connection string - update with your MSSQL credentials
		connStr = "server=localhost;user id=sa;password=YourPassword123;database=clean_db;port=1433"
	}

	db, err := sql.Open("mssql", connStr)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}

// getPort returns the port from environment or uses default
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
