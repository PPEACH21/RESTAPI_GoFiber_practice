package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

type Book struct {
	ID     int    `json:"id"`
	Tittle string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func main() {
	if err := godotenv.Load(); err!=nil{
		log.Fatal("load .env error")
	}
	engine := html.New("./views",".html")
	app := fiber.New( 
		fiber.Config{
			Views: engine,
	})

	books = append(books, Book{ID: 1, Tittle: "PEACHCER", Author: "PPEACH21"})
	books = append(books, Book{ID: 2, Tittle: "PEERAPAT", Author: "SAENGPHOEM"})

	app.Get("/books",getBooks)
	app.Get("/books/:id",getBook)
	app.Post("/books",createBook)
	app.Put("/editbook/:id",editBook)
	app.Delete("/deletebook/:id",deleteBook)
	app.Post("/upload",uploadFile)
	app.Get("/html",testHTML)
	app.Get("/api/config", getEnv)
	app.Listen(":8080")
}


func getEnv(c *fiber.Ctx)error{
	secret := os.Getenv("SECRET")
	if secret == ""{
		secret = "defaultsecret"
	}
	return c.JSON(fiber.Map{
			"SECRET": secret,
	})
}