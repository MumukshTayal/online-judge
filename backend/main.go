package main

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
	"strconv"

	"github.com/MumukshTayal/online-judge/add_contest"
	"github.com/MumukshTayal/online-judge/add_problem"
	"github.com/MumukshTayal/online-judge/edit_userProfile"
	"github.com/MumukshTayal/online-judge/fetch_userProfile"
	"github.com/MumukshTayal/online-judge/get_contest"
	"github.com/MumukshTayal/online-judge/get_problem"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random" // Change this to something more secure in production
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,HEAD",
	}))

	app.Post("/verify_creds", func(c *fiber.Ctx) error {
		fmt.Println("creds verified")
		return c.JSON("cool!")
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		fmt.Println("login post request received")

		authHeader := c.Get("Authorization")
		tokenStr := ""
		if authHeader != "" {
			authValue := strings.Split(authHeader, " ")
			if len(authValue) == 2 && authValue[0] == "Bearer" {
				tokenStr = authValue[1]
			}
		}

		token, err := jwt.Parse(tokenStr, nil)
		if token == nil {
			fmt.Println("err --> ", err)
			return nil
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		fmt.Println("email --> ", claims["email"])

		return c.JSON("cool!")
	})

	app.Get("/get_all_contests_users_pairs", get_contest.GetContestUsers)
	app.Get("/get_contests_by_userId", func(c *fiber.Ctx) error {
		userIDStr := c.Query("user_id")
		userId, err := strconv.Atoi(userIDStr)
		if err != nil {
			return err
		}

		contests, err := get_contest.GetContestByUserId(c, userId)
		if err != nil {
			return err
		}
		return c.JSON(contests)
	})

	app.Get("/get_contest_detail_by_contestId", func(c *fiber.Ctx) error {
		contestIDStr := c.Query("contest_id")
		contestId, err := strconv.Atoi(contestIDStr)
		if err != nil {
			return err
		}

		contests, err := get_contest.GetContestDetails(c, contestId)
		if err != nil {
			return err
		}
		return c.JSON(contests)
	})

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

	app.Get("/fetch_profile", func(c *fiber.Ctx) error {
		userIDStr := c.Query("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return err
		}

		userProfile, err := fetch_userProfile.GetUserProfileByUserId(c, userID)
		if err != nil {
			return err
		}
		return c.JSON(userProfile)
	})

	app.Post("/create_contest", add_contest.AddContest)
	app.Post("/create_problem", add_problem.AddProblem)
	app.Post("/add_problems_to_contest", add_contest.AddProblemIDandContestIDtoTable)
	app.Post("/edit_profile", edit_userProfile.EditUserProfile)

	app.Listen(":8080")
}
