package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Book struct {
	ID     int    `json:"id"`
	Tittle string `json:"title"`
	Author string `json:"author"`
}

type User = struct{
	Email 	string 	`json:"email"`
	Password string	`json:"password"`
}

var member = User{
	Email:"User@gmail.com",
	Password:"Password123",
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
	

	app.Post("/login",login)

	//JWT_SECRET
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	app.Use(checkMiddleware)

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


func login(c *fiber.Ctx) error{
	user := new(User)
	
	if err := c.BodyParser(user) ; err!=nil{
		return  c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if user.Email!= member.Email || user.Password != member.Password{
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