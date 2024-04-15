package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/MumukshTayal/online-judge/docker_containerize"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type TestJob struct {
	TestId int
}

var (
	jobQueue   = make(chan TestJob, 10)
	queueMutex sync.Mutex
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,HEAD",
	}))

	app.Post("/judge/check_sol", docker_containerize.Containerize)

	app.Post("/judge/add_to_queue", func(c *fiber.Ctx) error {
		var job TestJob
		err := c.BodyParser(&job)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		queueMutex.Lock()
		jobQueue <- job
		defer queueMutex.Unlock()

		return c.SendString("Test job added to the queue!")
	})

	go func() {
		for {
			// Attempt to receive a job from the queue (non-blocking)
			select {
			case job := <-jobQueue:
				check_solution()
				fmt.Println("Test Job ID:", job.TestId)
			default:
				// No job available, do nothing (avoid busy waiting)
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()

	app.Listen(":3001")
}

func check_solution() {
	url := "http://localhost:3001/judge/check_sol"

	submittedCode := docker_containerize.SubmittedCode{
		Code: "def add(numbers):\n    return 6",
	}
	body, err := json.Marshal(submittedCode)

	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("KI KI KI RRRR:", resp)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response status code:", resp.StatusCode)
		return
	}

	var responseBody []byte
	_, err = resp.Body.Read(responseBody)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response:", string(responseBody))
}
