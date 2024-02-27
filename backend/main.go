package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

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

	app.Get("/google_login", func(c *fiber.Ctx) error {
		url := googleOauthConfig.AuthCodeURL(oauthStateString)
		return c.Redirect(url)
	})

	app.Get("/google_callback", func(c *fiber.Ctx) error {
		code := c.Query("code")
		token, err := googleOauthConfig.Exchange(context.Background(), code)
		if err != nil {
			return err
		}

		userInfo, err := getUserInfo(token.AccessToken)
		if err != nil {
			return err
		}

		// Now you have userInfo, you can use it as needed
		// For example, store the email in the session or database
		sessionEmail := userInfo.Email

		fmt.Println("session email --> ", sessionEmail)

		return c.JSON(userInfo)
	})

	// Your existing routes...

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

func getUserInfo(accessToken string) (*UserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo UserInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

type UserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Verified bool   `json:"verified_email"`
}
