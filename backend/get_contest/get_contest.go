package get_contest

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func GetContestUsers(c *fiber.Ctx) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	db_name := os.Getenv("DATABASE_NAME")
	db_token := os.Getenv("DATABASE_TOKEN")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", db_name, db_token)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping db %s: %s\n", url, err)
		os.Exit(1)
	}
	defer db.Close()

	user_data, err := fetchUsers(db)
	if err != nil {
		return err
	}

	return c.JSON(user_data)

}

type User struct {
	ContestId int
	UserId    int
}

type Contest struct {
	ContestID          int       `json:"contest_id"`
	ContestTitle       string    `json:"contest_title"`
	ContestDescription string    `json:"contest_description"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	IsPublic           bool      `json:"is_public"`
	CreatorID          int       `json:"creator_id"`
}

func fetchUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM contest_user")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ContestId, &user.UserId); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		users = append(users, user)
		// fmt.Println(user.ContestId, user.UserId)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}
	return users, nil
}

func GetContestByUserId(c *fiber.Ctx, userId int) ([]int, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	db_name := os.Getenv("DATABASE_NAME")
	db_token := os.Getenv("DATABASE_TOKEN")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", db_name, db_token)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping db %s: %s\n", url, err)
		os.Exit(1)
	}
	defer db.Close()

	rows, err := db.Query("SELECT contest_id FROM contest_user WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contests []int

	for rows.Next() {
		var contestID int
		if err := rows.Scan(&contestID); err != nil {
			return nil, err
		}
		contests = append(contests, contestID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contests, nil

}

func GetContestDetails(c *fiber.Ctx, contestId int) (Contest, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	db_name := os.Getenv("DATABASE_NAME")
	db_token := os.Getenv("DATABASE_TOKEN")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", db_name, db_token)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping db %s: %s\n", url, err)
		os.Exit(1)
	}
	defer db.Close()

	var contest Contest

	row := db.QueryRow("SELECT * FROM contest WHERE contest_id = ?", contestId)

	err = row.Scan(&contest.ContestID, &contest.ContestTitle, &contest.ContestDescription, &contest.StartTime, &contest.EndTime, &contest.IsPublic, &contest.CreatorID)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No contest found with ID:", contestId)
			return contest, nil
		}

		fmt.Fprintf(os.Stderr, "failed to scan problem: %v\n", err)
		return contest, err
	}

	fmt.Println(contest.ContestTitle, contest.ContestDescription, contest.StartTime, contest.EndTime, contest.IsPublic, contest.CreatorID)
	return contest, nil
}
