package book

import (
	"github.com/gofiber/fiber"
	"github.com/tikam300/go-fiber/database"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

func GetBooks(c *fiber.Ctx) {
	db := database.DBConn
	var books []Book
	db.Find(&books)
	c.JSON(books)
}

func GetBook(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn
	var book Book
	db.Find(&book, id)
	c.JSON(book)
}

func NewBook(c *fiber.Ctx) {
	db := database.DBConn
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		c.Status(503).Send(err)
		return
	}

	db.Create(&book)
	c.JSON(book)
}

func DeleteBook(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn

	var book Book
	db.First(&book, id)
	if book.Title == "" {
		c.Status(500).Send("No book found with given Id")
		return
	}
	db.Delete(&book)
	c.Send("Book sucessfully deleted")
}

func UpdateBook(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn

	// Parse request body into a temporary struct
	updatedData := new(Book)
	if err := c.BodyParser(updatedData); err != nil {
		c.Status(400).Send("Invalid request body")
		return
	}

	// Find existing book
	var book Book
	result := db.First(&book, id)
	if result.RowsAffected == 0 {
		c.Status(404).Send("No book found with given ID")
		return
	}

	// Only update non-empty fields
	if updatedData.Author != "" {
		book.Author = updatedData.Author
	}
	if updatedData.Title != "" {
		book.Title = updatedData.Title
	}
	if updatedData.Rating != 0 {
		book.Rating = updatedData.Rating
	}

	// Save changes
	if err := db.Save(&book).Error; err != nil {
		c.Status(500).Send("Failed to update book")
		return
	}

	c.Send("Book updated successfully")
}
