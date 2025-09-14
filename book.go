package main

import (
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
)

func getUsers(c *fiber.Ctx) error {
	var user []User
	data := Client.Collection("User").Documents(ctx)
	for{
		doc,err:= data.Next()
		if err == iterator.Done{
			break;
		}
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		
		var u User
		if err := doc.DataTo(&u); err != nil {
			log.Println("error convert:", err)
			continue
		}

		user = append(user, User{
			Email: u.Email,
			Password: u.Password,
		})
	}
	return c.JSON(user)
}

func getUser(c *fiber.Ctx) error {
	UserID := c.Params("id")

	data,err:= Client.Collection("User").Doc(UserID).Get(ctx)
	if err != nil{
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	m := data.Data()
	fmt.Println(m)
	return c.JSON(m)
}

func getBooks(c *fiber.Ctx) error {
	var book []Book
	data := Client.Collection("Books").OrderBy("createdAt", firestore.Asc).Documents(ctx)
	for{
		doc,err:= data.Next()
		if err == iterator.Done{
			break;
		}
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		
		var u Book
		if err := doc.DataTo(&u); err != nil {
			log.Println("error convert:", err)
			continue
		}

		book = append(book, Book{
			Tittle: u.Tittle,
			Author: u.Author,
		})
	}
	return c.JSON(book)
}

func createBook(c *fiber.Ctx) error {
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if book.Tittle == "" || book.Author == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title and author are required",
		})
	}

	title := strings.TrimSpace(book.Tittle)
	author := strings.TrimSpace(book.Tittle)

	_, _, err := Client.Collection("Books").Add(ctx, map[string]interface{}{
		"title":     title,
		"author":    author,
		"createdAt": firestore.ServerTimestamp,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(fiber.Map{
		"title":     title,
		"author":    author,
	})
}

// func editBook(c *fiber.Ctx) error {
// 	bookId, err := strconv.Atoi(c.Params("id"))

// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}

// 	bookUpdate := new(Book)

// 	if err := c.BodyParser(bookUpdate); err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}

// 	for i, book := range books {
// 		if book.ID == bookId {
// 			books[i].Tittle = bookUpdate.Tittle
// 			books[i].Author = bookUpdate.Author
// 			return c.JSON(books[i])
// 		}
// 	}
// 	return c.Status(fiber.StatusNotFound).SendString("Book not Found")
// }

// func deleteBook(c *fiber.Ctx) error {
// 	bookId, err := strconv.Atoi(c.Params("id"))

// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}

// 	for i, book := range books {
// 		if book.ID == bookId {
// 			books = append(books[:i], books[i+1:]...)
// 			return c.Status(fiber.StatusAccepted).SendString("Delete complete")
// 		}
// 	}
// 	return c.Status(fiber.StatusNotFound).SendString("Book not Found")
// }

// func uploadFile(c *fiber.Ctx) error{
// 	file,err := c.FormFile("image")

// 	if err != nil{
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}

// 	err = c.SaveFile(file,"./uploads/"+file.Filename)

// 	if err != nil{
// 		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
// 	}
// 	return c.SendString("File Upload Complete!")
// }

// //render HTML render Engine