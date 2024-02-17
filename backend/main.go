package main

import (
	"strconv"

	"github.com/MumukshTayal/online-judge/add_contest"
	"github.com/MumukshTayal/online-judge/add_problem"
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
	app.Get("/get_all_contests_users_pairs", get_contestUsers.GetContestUsers)
	app.Get("/get_problem_by_problemId", func(c *fiber.Ctx) error {
		problemIDStr := c.Query("problem_id")
		problemID, err := strconv.Atoi(problemIDStr)
		if err != nil {
			return err
		}

		return get_problem.GetProblemByProblemId(c, problemID)
	})
	app.Get("/get_problem_by_userId", func(c *fiber.Ctx) error {
		userIDStr := c.Query("user_id")
		userId, err := strconv.Atoi(userIDStr)
		if err != nil {
			return err
		}

		problems, err := get_problem.GetProblemByUserId(c, userId)
		if err != nil {
			return err
		}
		return c.JSON(problems)
	})
	app.Post("/create_contest", add_contest.AddContest)
	app.Post("/create_problem", add_problem.AddProblem)
	app.Post("/add_problems_to_contest", add_contest.AddProblemIDandContestIDtoTable)
	app.Listen(":8080")

}
