package main

import (
	"github.com/gofiber/fiber/v2"
)

type Book struct {
	ID     int    `json:"id"`
	Tittle string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func main() {
	app := fiber.New()

	books = append(books, Book{ID: 1, Tittle: "PEACHCER", Author: "PPEACH21"})
	books = append(books, Book{ID: 2, Tittle: "PEERAPAT", Author: "SAENGPHOEM"})

	app.Get("/books",getBooks)
	app.Get("/books/:id",getBook)
	app.Post("/books",createBook)
	app.Put("/editbook/:id",editBook)
	app.Delete("/deletebook/:id",deleteBook)
	app.Post("/upload",uploadFile)

	app.Listen(":8080")
}