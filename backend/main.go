package main

import (
	"strconv"

	"github.com/MumukshTayal/online-judge/config"
	"github.com/MumukshTayal/online-judge/controllers"
	"github.com/MumukshTayal/online-judge/get_contestUsers"
	"github.com/MumukshTayal/online-judge/get_problem"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.GoogleConfig()

	app.Get("/google_login", controllers.GoogleLogin)
	app.Get("/google_callback", controllers.GoogleCallback)
	app.Get("/get_contest_users", get_contestUsers.GetContestUsers)
	app.Get("/get_problem", func(c *fiber.Ctx) error {
		problemIDStr := c.Query("problem_id")
		problemID, err := strconv.Atoi(problemIDStr)
		if err != nil {
			return err
		}

		return get_problem.GetProblem(c, problemID)
	})
	app.Listen(":8080")

}
