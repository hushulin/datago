package main

// thanks to https://github.com/Learn-by-doing/csrf-examples
import (
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"main/database"
	"main/model"
	"main/routes"
	"net/http"
)

//go:embed views/*
var views embed.FS

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("fish.db"))
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection Opened to Database")
	database.DBConn.AutoMigrate(&model.FishAddress{})
	database.DBConn.AutoMigrate(&model.MdLog{})
	database.DBConn.AutoMigrate(&model.MdClick{})
	database.DBConn.AutoMigrate(&model.Attachment{})
	fmt.Println("Database Migrated")
}

func main() {
	engine := html.NewFileSystem(http.FS(views), ".html").Layout("views/layouts/main")
	// Fiber instance
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(logger.New())
	initDatabase()
	routes.RegisterRoutes(app)
	app.Static("/js/", "./res/js")
	log.Fatal(app.Listen(":3001"))
}
