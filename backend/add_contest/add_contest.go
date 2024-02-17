package add_contest

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func AddContest(c *fiber.Ctx) error {
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

	if err := AddContestsToContestTable(c, db); err != nil {
		return err
	}

	return c.SendString("Contest added successfully!")

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

type Distribute struct {
	ContestID  int   `json:"contest_id"`
	ProblemIDs []int `json:"problems_in_contest"`
}

func AddContestsToContestTable(c *fiber.Ctx, db *sql.DB) error {
	var contest Contest
	if err := c.BodyParser(&contest); err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO contest (contest_title, contest_description, contest_start_time, contest_end_time, is_public, creator_id) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(contest.ContestTitle, contest.ContestDescription, contest.StartTime, contest.EndTime, contest.IsPublic, contest.CreatorID)
	if err != nil {
		return err
	}

	return nil
}

func fetAllContests(db *sql.DB) ([]Contest, error) {
	rows, err := db.Query("SELECT * FROM contest")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var contests []Contest

	for rows.Next() {
		var contest Contest

		if err := rows.Scan(&contest.ContestID, &contest.ContestTitle, &contest.ContestDescription, &contest.StartTime, &contest.EndTime, &contest.IsPublic, &contest.CreatorID); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		contests = append(contests, contest)
		fmt.Println(contest.ContestID, contest.ContestTitle, contest.ContestDescription, contest.StartTime, contest.EndTime, contest.IsPublic, contest.CreatorID)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}
	return contests, nil
}

func AddProblemIDandContestIDtoTable(c *fiber.Ctx) error {
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

	var distribute Distribute
	if err := c.BodyParser(&distribute); err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO distribute_problems_to_contest (contest_id, problem_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	addedProblems := make([]string, 0)
	for _, problemID := range distribute.ProblemIDs {
		_, err := stmt.Exec(distribute.ContestID, problemID)
		if err != nil {
			return err
		}
		addedProblems = append(addedProblems, strconv.Itoa(problemID))
	}

	responseMessage := fmt.Sprintf("Problems %s added to contest %d successfully", strings.Join(addedProblems, ", "), distribute.ContestID)
	return c.SendString(responseMessage)
}
