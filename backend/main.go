package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
	"github.com/MumukshTayal/online-judge/get_leaderboard"
	"github.com/MumukshTayal/online-judge/get_problem"
	"github.com/MumukshTayal/online-judge/get_submissions"

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

type Runny struct {
	ProblemID      string    `json:"problem_id"`
	ContestID      string    `json:"contest_id"`
	UserEmail      string    `json:"user_email"`
	Code           string    `json:"code"`
	Language       string    `json:"language"`
	SubmissionTime time.Time `json:"submission_date_time"`
}

type Testcase struct {
	TestInput  []byte
	TestOutput []byte
}

type Constraint struct {
	TimeLimit int
	MemLimit  int
}

type PrepareForJuding struct {
	TestInpt   string `json:"test_input"`
	TestOutput string `json:"test_output"`
	TestCode   string `json:"test_code"`
	TimeLimit  int    `json:"time_limit"`
	MemLimit   int    `json:"memory_limit"`
	Language   string `json:"language"`
}

func RetrieveTextFromBlob(inputBlob []byte, outputBlob []byte) (string, string, error) {
	inputText := string(inputBlob)
	outputText := string(outputBlob)
	return inputText, outputText, nil
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

	app.Post("/api/run_code", func(c *fiber.Ctx) error {
		fmt.Println("INSIDE ADD TO QUEUE!!!!")
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
			fmt.Println("Error parsing token:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		userEmail, ok := claims["email"].(string)
		fmt.Println("USER_EMAIL:", userEmail)
		if !ok {
			fmt.Println("Invalid token claims")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid token claims",
			})
		}

		var runny Runny
		if err := c.BodyParser(&runny); err != nil {
			return err
		}

		problemID := runny.ProblemID
		contestId := runny.ContestID
		lang := runny.Language
		fmt.Println("Contest:", contestId, "Problem:", problemID)
		err2 := godotenv.Load(".env")
		if err2 != nil {
			log.Fatalf("Some error occurred. Err: %s", err2)
		}

		dbName := os.Getenv("DATABASE_NAME")
		dbToken := os.Getenv("DATABASE_TOKEN")

		url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, dbToken)

		db, err := sql.Open("libsql", url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
			return err
		}
		defer db.Close()
		input, output, err := fetchTestcases(db, problemID)
		if err != nil {
			return err
		}

		time_limit, mem_limit, err := fetchConstraints(db, contestId)
		if err != nil {
			return err
		}

		sendToJudge := PrepareForJuding{
			TestInpt:   input,
			TestOutput: output,
			TestCode:   runny.Code,
			TimeLimit:  time_limit,
			MemLimit:   mem_limit,
			Language:   lang,
		}
		data, err := json.Marshal(sendToJudge)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		judgeURL := "http://127.0.0.1:3001/judge/add_to_queue"
		resp, err := http.Post(judgeURL, "application/json", bytes.NewReader(data))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Received non-OK response status code: %d", resp.StatusCode)
			return c.SendStatus(http.StatusInternalServerError)
		}
		respBody, err := io.ReadAll(resp.Body)
		// fmt.Println(string(respBody))
		return c.SendString(string(respBody))
	})

	app.Post("/api/submit", func(c *fiber.Ctx) error {
		fmt.Println("INSIDE ADD TO QUEUE!!!!")
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
			fmt.Println("Error parsing token:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		userEmail, ok := claims["email"].(string)
		fmt.Println("USER_EMAIL:", userEmail)
		if !ok {
			fmt.Println("Invalid token claims")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid token claims",
			})
		}

		var runny Runny
		if err := c.BodyParser(&runny); err != nil {
			return err
		}

		problemID := runny.ProblemID
		contestId := runny.ContestID
		lang := runny.Language
		fmt.Println("Contest:", contestId, "Problem:", problemID)
		err2 := godotenv.Load(".env")
		if err2 != nil {
			log.Fatalf("Some error occurred. Err: %s", err2)
		}

		dbName := os.Getenv("DATABASE_NAME")
		dbToken := os.Getenv("DATABASE_TOKEN")

		url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, dbToken)

		db, err := sql.Open("libsql", url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
			return err
		}
		defer db.Close()
		input, output, err := fetchTestcases(db, problemID)
		if err != nil {
			return err
		}

		time_limit, mem_limit, err := fetchConstraints(db, contestId)
		if err != nil {
			return err
		}
		sendToJudge := PrepareForJuding{
			TestInpt:   input,
			TestOutput: output,
			TestCode:   runny.Code,
			TimeLimit:  time_limit,
			MemLimit:   mem_limit,
			Language:   lang,
		}
		data, err := json.Marshal(sendToJudge)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		judgeURL := "http://127.0.0.1:3001/judge/add_to_queue"
		resp, err := http.Post(judgeURL, "application/json", bytes.NewReader(data))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Received non-OK response status code: %d", resp.StatusCode)
			return c.SendStatus(http.StatusInternalServerError)
		}
		respBody, err := io.ReadAll(resp.Body)
		respString := string(respBody)
		lines := strings.Split(respString, "\n")
		var totalTestcases, testCasesPassed int
		var timeElapsed float64
		for _, line := range lines {
			fields := strings.Split(line, ":")
			if len(fields) != 2 {
				continue
			}
			fieldName := strings.TrimSpace(fields[0])
			fieldValue := strings.TrimSpace(fields[1])

			switch fieldName {
			case "Total Testcases":
				totalTestcases, _ = strconv.Atoi(fieldValue)
			case "Test Cases Passed":
				testCasesPassed, _ = strconv.Atoi(fieldValue)
			case "Time Elapsed (msec)":
				timeElapsed, _ = strconv.ParseFloat(fieldValue, 64)
			}
		}

		var result string
		if testCasesPassed == totalTestcases {
			result = "AC"
		} else if testCasesPassed < totalTestcases {
			result = "WA"
		}

		if timeElapsed > float64(time_limit) {
			result = "TLE"
		}

		fmt.Println("Result:", result)

		stmt, err := db.Prepare("INSERT INTO submission (problem_id, user_email, code, language, result, execution_time, memory_used, submission_date_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(runny.ProblemID, userEmail, runny.Code, runny.Language, result, timeElapsed, mem_limit, runny.SubmissionTime)

		// fmt.Println("YOOOO", runny.ProblemID, userEmail, runny.Code, runny.Language, result, timeElapsed, mem_limit, runny.SubmissionTime)
		if err != nil {
			return err
		}

		return c.SendString("Submission Added Successfully")
	})

	//app.Post("/api/run_code", run_code.RunCode)
	// app.Post("/api/submit", add_submission.AddSubmission)
	app.Post("/api/create_contest", add_contest.AddContest)
	app.Post("/api/create_problem", add_problem.AddProblem)
	app.Post("/edit_profile", edit_userProfile.EditUserProfile)
	app.Post("/api/add_testcase", add_testcase.AddTestCase)

	app.Get("/api/get_all_contests", get_contest.GetAllContests)
	app.Get("/api/get_all_problems", get_problem.GetAllProblems)
	app.Get("/api/get_all_submissions", get_submissions.GetAllSubmissions)
	app.Get("/api/get_contest_list", get_contest.ContestList)
	app.Get("/api/get_contest_details/:contestId", get_contest_details.GetContestDetails)
	app.Get("/api/get_problem/:problemId", get_problem.GetProblemByProblemId)
	app.Get("/api/get_leaderboard", get_leaderboard.GetLeaderboard)
	app.Get("/api/get_submissions", get_submissions.GetSubmissions)

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

func fetchTestcases(db *sql.DB, problemID string) (string, string, error) {
	var testcases Testcase

	row := db.QueryRow("SELECT testcase_input, testcase_output FROM testcase WHERE problem_id = ?", problemID)
	err := row.Scan(&testcases.TestInput, &testcases.TestOutput)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No testcases found with ID:", problemID)
			return "", "", err
		}

		fmt.Fprintf(os.Stderr, "failed to scan problem: %v\n", err)
		return "", "", err
	}

	inputText, outputText, err := RetrieveTextFromBlob(testcases.TestInput, testcases.TestOutput)
	if err != nil {
		return "", "", err
	}

	return inputText, outputText, nil
}
func fetchConstraints(db *sql.DB, contestID string) (int, int, error) {
	var constraint Constraint

	row := db.QueryRow("SELECT time_limit, memory_limit FROM allowed_list WHERE contest_id = ?", contestID)
	err := row.Scan(&constraint.TimeLimit, &constraint.MemLimit)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No constraints for:", contestID)
			return 0, 0, err
		}

		fmt.Fprintf(os.Stderr, "failed to scan constraints: %v\n", err)
		return 0, 0, err
	}

	time_limit := &constraint.TimeLimit
	memory_limit := &constraint.MemLimit

	return *time_limit, *memory_limit, nil
}
