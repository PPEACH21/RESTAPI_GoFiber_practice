package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Book struct {
	Title string `json:"title"`
	Author string `json:"author"`
}

type User = struct{
	Email 	string 	`json:"email"`
	Password string	`json:"password"`
}

var Client *firestore.Client
var ctx = context.Background() 

func main() {
	app := fiber.New()
	
	opt := option.WithCredentialsFile("FirebaseKey.json")
	Fbapp, err := firebase.NewApp(ctx, nil,opt)
	if err != nil {
	log.Fatalln(err)
	}

	Client, err = Fbapp.Firestore(ctx)
	if err != nil {
	log.Fatalln(err)
	}
	defer Client.Close()


	app.Post("/login",login)

	app.Get("/users",getUsers)
	// JWT_SECRET
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	app.Use(checkMiddleware)

		app.Get("/users/:id",getUser)
		app.Post("/users",createUser)
		app.Put("/edituser/:id",editUser)
		app.Delete("/deleteusers/:id",deleteUser)


		app.Get("/books",getBooks)
		app.Get("/books/:id",getBookID)
		app.Post("/books",createBook)
		app.Put("/editbook/:id",editBook)
		app.Delete("/deletebook/:id",deleteBook)
		// app.Post("/upload",uploadFile)
		app.Get("/api/config", getEnv)
	app.Listen(":8080")
}


func login(c *fiber.Ctx) error{
	user := new(User)
	
	if err := c.BodyParser(user) ; err!=nil{
		return  c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	docs,err := Client.Collection("User").
		Where("email", "==", user.Email).
        Limit(1).
        Documents(ctx).
        Next()
	if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

	var member User
    if err := docs.DataTo(&member); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Error parsing user data")
    }
	// ตรวจสอบ password (plain text)
    if user.Password != member.Password {
        return fiber.ErrUnauthorized
    }


	// Create token
    token := jwt.New(jwt.SigningMethodHS256)

    // Set claims
    claims := token.Claims.(jwt.MapClaims)
    claims["email"] = user.Email
    claims["role"] = "admin"
    claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

    // Generate encoded token
    t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
      return c.SendStatus(fiber.StatusInternalServerError)
    }

	return c.JSON(fiber.Map{
		"message" : "Login success",
		"token": t,
	})
}

func checkMiddleware(c *fiber.Ctx)error{
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	
	if claims ["role"] != "admin"{
		return  fiber.ErrUnauthorized
	}

	start:=time.Now().In(time.FixedZone("UTC+7", 7*60*60))

	fmt.Printf(
		"URL = %s, Method = %s, Time = %s\n",
		c.OriginalURL(),c.Method(),start,
	)
	return c.Next()
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