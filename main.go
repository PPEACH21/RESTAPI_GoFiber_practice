package main

import (
	"strconv"

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
	app.Listen(":8080")
}

func getBooks(c *fiber.Ctx)error{
	return c.JSON(books)
}

func getBook(c *fiber.Ctx)error{
	bookID,err := strconv.Atoi(c.Params("id"))

	
	if err!= nil{
		return  c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	
	for _, book := range books{
		if book.ID == bookID{
			return c.JSON(book)
		}
	} 

	return  c.Status(fiber.StatusNotFound).SendString("Book not Found")
}

func createBook(c *fiber.Ctx)error{
	book := new(Book)
	if err := c.BodyParser(book) ; err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	books = append(books,*book)
	return  c.JSON(book)
}

func editBook(c *fiber.Ctx)error{
	bookId,err := strconv.Atoi(c.Params("id"))

	if err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	bookUpdate := new(Book)

	if err := c.BodyParser(bookUpdate) ; err!= nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, book := range books{
		if book.ID == bookId{
			books[i].Tittle = bookUpdate.Tittle
			books[i].Author = bookUpdate.Author
			return c.JSON(books[i])
		}
	} 
	return  c.Status(fiber.StatusNotFound).SendString("Book not Found")
}

func deleteBook(c *fiber.Ctx)error{
	bookId,err := strconv.Atoi(c.Params("id"))
	
	if err!= nil{
		return  c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i,book := range books{
		if book.ID == bookId{
			books=append(books[:i],books[i+1:]...)
			return c.Status(fiber.StatusAccepted).SendString("Delete complete")
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("Book not Found")
}