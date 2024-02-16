package main

import (
	"github.com/MumukshTayal/online-judge/config"
	"github.com/MumukshTayal/online-judge/controllers"
	"github.com/MumukshTayal/online-judge/get_contestUsers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.GoogleConfig()

	app.Get("/google_login", controllers.GoogleLogin)
	app.Get("/google_callback", controllers.GoogleCallback)
	app.Get("/get_contest_users", get_contestUsers.GetContestUsers)
	app.Listen(":8080")

}
