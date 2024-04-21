package get_contest

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Contest struct {
	ContestID    int       `json:"contest_id"`
	ContestTitle string    `json:"contest_title"`
	ContestDesc  string    `json:"contest_description"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	IsPublic     string    `json:"is_public"`
	CreatorEmail string    `json:"creator_email"`
}

type ContestView struct {
	ContestID    int       `json:"contest_id"`
	ContestTitle string    `json:"contest_title"`
	ContestDesc  string    `json:"contest_description"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

func GetAllContests(c *fiber.Ctx) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	dbName := os.Getenv("DATABASE_NAME")
	dbToken := os.Getenv("DATABASE_TOKEN")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, dbToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		return c.Status(500).SendString(err.Error())
	}

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping db %s: %s\n", url, err)
		return c.Status(500).SendString(err.Error())
	}
	defer db.Close()

	contests, err := fetchAllContests(db)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(contests)
}

func fetchAllContests(db *sql.DB) ([]Contest, error) {
	rows, err := db.Query("SELECT contest_id, contest_title, contest_description, contest_start_time, contest_end_time, is_public, creator_email FROM contest")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var contests []Contest

	for rows.Next() {
		var contest Contest
		var creatorEmail sql.NullString

		err := rows.Scan(&contest.ContestID, &contest.ContestTitle, &contest.ContestDesc, &contest.StartTime, &contest.EndTime, &contest.IsPublic, &creatorEmail)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		if creatorEmail.Valid {
			contest.CreatorEmail = creatorEmail.String
		} else {
			contest.CreatorEmail = ""
		}

		contests = append(contests, contest)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
		return nil, err
	}

	return contests, nil
}

func ContestList(c *fiber.Ctx) error {
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
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email in token claims",
		})
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	dbName := os.Getenv("DATABASE_NAME")
	dbToken := os.Getenv("DATABASE_TOKEN")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, dbToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		return c.Status(500).SendString(err.Error())
	}

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping db %s: %s\n", url, err)
		return c.Status(500).SendString(err.Error())
	}
	defer db.Close()

	participatingContests, err := getParticipatingContests(userEmail, db)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	createdContests, err := getCreatedContests(userEmail, db)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"participating_contests": participatingContests,
		"created_contests":       createdContests,
	})
}

func getParticipatingContests(userEmail string, db *sql.DB) ([]ContestView, error) {
	rows, err := db.Query("SELECT c.contest_id, c.contest_title, c.contest_description, c.contest_start_time, c.contest_end_time FROM contest c JOIN contest_user cu ON c.contest_id = cu.contest_id WHERE cu.user_email = ?", userEmail)
	if err != nil {
		fmt.Println("Error fetching participating contests:", err)
		return nil, err
	}
	defer rows.Close()

	return parseContests(rows)
}

func getCreatedContests(userEmail string, db *sql.DB) ([]ContestView, error) {
	rows, err := db.Query("SELECT contest_id, contest_title, contest_description, contest_start_time, contest_end_time FROM contest WHERE creator_email = ?", userEmail)
	if err != nil {
		fmt.Println("Error fetching created contests:", err)
		return nil, err
	}
	defer rows.Close()

	return parseContests(rows)
}

func parseContests(rows *sql.Rows) ([]ContestView, error) {
	var contests []ContestView
	for rows.Next() {
		var contest ContestView

		if err := rows.Scan(&contest.ContestID, &contest.ContestTitle, &contest.ContestDesc, &contest.StartTime, &contest.EndTime); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		// fmt.Println(&contest.ContestID, &contest.ContestTitle, &contest.ContestDesc, &contest.StartTime, &contest.EndTime)
		contests = append(contests, contest)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
		return nil, err
	}

	return contests, nil
}
