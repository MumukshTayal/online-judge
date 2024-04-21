package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/MumukshTayal/online-judge/docker_containerize"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type TestJob struct {
	// TestId int
	Ctx *fiber.Ctx
}

type JobResult struct {
	Output string
	Err    error
}

type PrepareForJuding struct {
	TestInpt   string `json:"test_input"`
	TestOutput string `json:"test_output"`
	TestCode   string `json:"test_code"`
}

var (
	jobQueue   = make(chan TestJob, 10)
	resultChan = make(chan JobResult, 10)
	queueMutex sync.Mutex
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,HEAD",
	}))

	app.Post("/judge/add_to_queue", func(c *fiber.Ctx) error {
		// fmt.Println("YO YO YO YO YO YO YO!!!")

		var job TestJob
		job.Ctx = c
		// err := c.BodyParser(&job)
		// if err != nil {
		// 	return c.SendStatus(fiber.StatusBadRequest)
		// }
		queueMutex.Lock()
		jobQueue <- job
		defer queueMutex.Unlock()

		// Wait for the result
		result := <-resultChan
		if result.Err != nil {
			fmt.Println("Error processing job:", result.Err)
			return result.Err
		}
		fmt.Println(result.Output)
		// fmt.Println(result.Output)
		c.SendString(result.Output)
		return nil
	})

	go func() {
		for {
			// Attempt to receive a job from the queue (non-blocking)
			select {
			case job := <-jobQueue:
				output, err := docker_containerize.Containerize(job.Ctx)
				resultChan <- JobResult{Output: output, Err: err}
			default:
				// No job available, do nothing (avoid busy waiting)
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()

	app.Listen(":3001")
}
