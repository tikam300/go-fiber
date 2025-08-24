// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/gofiber/fiber"
// 	"github.com/tikam300/go-fiber/book"
// 	"github.com/tikam300/go-fiber/database"

// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"

// 	_ "modernc.org/sqlite"
// )

// func helloWorld(c *fiber.Ctx) {
// 	c.Send("Hello World")
// }

// func setupRoutes(app *fiber.App) {
// 	app.Get("/api/v1/book", book.GetBooks)
// 	app.Get("/api/v1/book/:id", book.GetBook)
// 	app.Post("/api/v1/book", book.NewBook)
// 	app.Delete("/api/v1/book/:id", book.DeleteBook)
// 	app.Put("/api/v1/book/:id", book.UpdateBook)
// }

// func initDatabase() (*gorm.DB, error) {
// 	var err error
// 	database.DBConn, err = gorm.Open(sqlite.Open("books.db"), &gorm.Config{})

// 	if err != nil {
// 		// panic("Failed to connect to database: " + err.Error())
// 		return nil, err
// 	}

// 	fmt.Println("Database connection successfully opened")
// 	return database.DBConn, nil

// }

// func main() {
// 	app := fiber.New()

// 	_, err := initDatabase()
// 	if err != nil {
// 		// panic("failed to get sql.DB from gorm")
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}

// 	setupRoutes(app)

// 	app.Listen(3000)
// }

package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber"
	"github.com/tikam300/go-fiber/book"
	"github.com/tikam300/go-fiber/database"

	"github.com/glebarez/sqlite" // pure Go sqlite driver
	"gorm.io/gorm"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/book", book.GetBooks)
	app.Get("/api/v1/book/:id", book.GetBook)
	app.Post("/api/v1/book", book.NewBook)
	app.Delete("/api/v1/book/:id", book.DeleteBook)
	app.Put("/api/v1/book/:id", book.UpdateBook)
}

func initDatabase() (*sql.DB, error) {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := database.DBConn.DB()
	if err != nil {
		return nil, err
	}
	fmt.Println("Database connection successfully opened")
	database.DBConn.AutoMigrate(&book.Book{})
	fmt.Println("Database Migrated")
	return sqlDB, nil
}

func main() {
	app := fiber.New()

	sqlDB, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlDB.Close()

	setupRoutes(app)
	app.Listen(3000)
}
