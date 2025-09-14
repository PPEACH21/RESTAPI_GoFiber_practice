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
		fmt.Println(user)
	}

	if user == nil{
		return c.Status(fiber.StatusNotFound).SendString("User is Empty")
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

func createUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and Password are required",
		})
	}

	email := strings.TrimSpace(user.Email)
	Password := strings.TrimSpace(user.Password)

	_, _,err:= Client.Collection("User").Add(ctx,map[string] interface{}{
		"email" : email,
		"password" :Password,

	})
	if err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"email":     email,
		"password":    Password,
	})
}

func editUser(c *fiber.Ctx) error{
	userID := c.Params("id")
	
	data := Client.Collection("User").Doc(userID)

	_,err := data.Get(ctx)
	if err !=nil{
		return  c.Status(fiber.StatusNotFound).SendString("User not Found")
	}

	datanew := new(User)
	if err := c.BodyParser(datanew); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	email := strings.TrimSpace(datanew.Email)
	password := strings.TrimSpace(datanew.Password)

	_,err = data.Update(ctx,[]firestore.Update{
		{
			Path: "email",
			Value: email,
		},
		{
			Path: "password",
			Value: password,
		},
	})
	if err != nil {
    	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"email" : email,
		"password" : password,
	})
}

func deleteUser(c *fiber.Ctx)error{
	UserID := c.Params("id")
	
	data:= Client.Collection("User").Doc(UserID);

	_,err:=data.Get(ctx);
	if err != nil{
		return  c.Status(fiber.StatusNotFound).SendString("User not Found")
	}

	_,err = data.Delete(ctx)
	if err!=nil{
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	
	return  c.Status(fiber.StatusOK).SendString("Delete Success")
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
			Title:u.Title,
			Author: u.Author,
		})
	}
	if book==nil{
		return c.Status(fiber.StatusNotFound).SendString("Book is Empty")
	}
	return c.JSON(book)
}

func getBookID(c *fiber.Ctx) error {
	BookID := c.Params("id")

	data,err:= Client.Collection("Books").Doc(BookID).Get(ctx)
	if err != nil{
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	m := data.Data()
	fmt.Println(m)
	return c.JSON(m)
}

func createBook(c *fiber.Ctx) error {
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if book.Title == "" || book.Author == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title and author are required",
		})
	}

	title := strings.TrimSpace(book.Title)
	author := strings.TrimSpace(book.Author)

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

func editBook(c *fiber.Ctx) error {
	bookId := c.Params("id")

	bookUpdate := new(Book)
	
	if err := c.BodyParser(bookUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	title := strings.TrimSpace(bookUpdate.Title)
	author := strings.TrimSpace(bookUpdate.Author)
	
	_,err := Client.Collection("Books").Doc(bookId).Update(ctx,[]firestore.Update{
		{
			Path: "title",
			Value: title,
	
		},	
		{
			Path: "author",
			Value: author,	
		},	

	})
	if err != nil {
    	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"title" : title,
		"author" : author,
	})
}

func deleteBook(c *fiber.Ctx) error {
	bookId := c.Params("id")
	if bookId == "" {
        return c.Status(fiber.StatusBadRequest).SendString("Missing book ID")
    }

	docRef := Client.Collection("Books").Doc(bookId)

	_,err := docRef.Get(ctx)
	if err != nil{
		 return c.Status(fiber.StatusNotFound).SendString("Book not found")
	}

	_,err = docRef.Delete(ctx)
	if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
    }

	return c.Status(fiber.StatusOK).SendString("Delete Success")
}

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
