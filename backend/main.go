package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"

	"github.com/MumukshTayal/online-judge/add_contest"
	"github.com/MumukshTayal/online-judge/add_problem"
	"github.com/MumukshTayal/online-judge/add_testcase"
	"github.com/MumukshTayal/online-judge/edit_userProfile"
	"github.com/MumukshTayal/online-judge/fetch_userProfile"
	"github.com/MumukshTayal/online-judge/get_contest"
	"github.com/MumukshTayal/online-judge/get_contest_details"
	"github.com/MumukshTayal/online-judge/get_problem"

	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random" // Change this to something more secure in production
	judgeURL          = "http://judge-server:3001"
)

type TestJob struct {
	TestId int
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

	app.Post("/api/login", func(c *fiber.Ctx) error {
		fmt.Println("we are in add user !!")
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

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fmt.Errorf("Failed to parse JWT claims")
		}

		userEmail := claims["email"].(string)
		userFirstName := claims["given_name"].(string)
		userLastName := claims["family_name"].(string)

		fmt.Println("email --> ", userEmail)

		err = godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err)
		}

		dbName := os.Getenv("DATABASE_NAME")
		dbToken := os.Getenv("DATABASE_TOKEN")

		url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, dbToken)

		db, err := sql.Open("libsql", url)
		if err != nil {
			log.Fatalf("Error connecting to database: %s", err)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Fatalf("Error pinging database: %s", err)
		}

		if err := AddUserToUserProfileTable(db, userEmail, userFirstName, userLastName); err != nil {
			return err
		}

		return c.JSON("User added successfully!")
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

	app.Post("/judge/add_to_queue", func(c *fiber.Ctx) error {
		msg := TestJob{TestId: 1}

		data, err := json.Marshal(msg)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		resp, err := http.Post(judgeURL, "application/json", bytes.NewReader(data))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Received non-OK response status code: %d", resp.StatusCode)
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.SendString("POST request sent successfully")
	})

	app.Post("/api/create_contest", add_contest.AddContest)
	app.Post("/api/create_problem", add_problem.AddProblem)
	app.Post("/edit_profile", edit_userProfile.EditUserProfile)
	app.Post("/api/add_testcase", add_testcase.AddTestCase)

	app.Get("/api/get_all_contests", get_contest.GetAllContests)
	app.Get("/api/get_contest_details/:contestId", get_contest_details.GetContestDetails)
	app.Get("/api/get_problem/:problemId", get_problem.GetProblemByProblemId)

	app.Listen(":8080")
}

func AddUserToUserProfileTable(db *sql.DB, email, firstName, lastName string) error {
	userName := fmt.Sprintf("%s %s", firstName, lastName)
	stmt, err := db.Prepare("INSERT INTO user_profile (user_email, user_name) SELECT ?, ? WHERE NOT EXISTS (SELECT 1 FROM user_profile WHERE user_email = ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(email, userName, email)
	if err != nil {
		return err
	}
	return nil
}
