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
		fmt.Println("YO YO YO YO YO YO YO!!!")
		// fmt.Println("BODY:", c.Body())
		// authHeader := c.Get("Authorization")
		// tokenStr := ""
		// if authHeader != "" {
		// 	authValue := strings.Split(authHeader, " ")
		// 	if len(authValue) == 2 && authValue[0] == "Bearer" {
		// 		tokenStr = authValue[1]
		// 	}
		// }
		// fmt.Println("TOKEN STR:", tokenStr)

		// Define a struct to hold the JSON data
		// var requestData PrepareForJuding

		// if err := json.Unmarshal(c.Body(), &requestData); err != nil {
		// 	fmt.Println("Error unmarshaling JSON:", err)
		// 	return err
		// }

		// fmt.Println("Test Input:", requestData.TestInpt)
		// fmt.Println("Test Output:", requestData.TestOutput)
		// fmt.Println("Test Code:", requestData.TestCode)
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
		fmt.Println("Output:", result.Output)
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
				// fmt.Println("Test Job ID:", job.TestId)
			default:
				// No job available, do nothing (avoid busy waiting)
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()

	app.Listen(":3001")
}

// func check_solution(c *fiber.Ctx) {

// 	str := docker_containerize.Containerize(c)

// 	fmt.Println("STR:", str)

// 	// url := "http://localhost:3001/judge/check_sol"

// 	// // testcases_url := "http://localhost:3000/add_to_queue"

// 	// submittedCode := docker_containerize.SubmittedCode{
// 	// 	Code: "def add(numbers):\n    return 6",
// 	// }
// 	// body, err := json.Marshal(submittedCode)

// 	// if err != nil {
// 	// 	fmt.Println("Error encoding JSON:", err)
// 	// 	return
// 	// }

// 	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
// 	// if err != nil {
// 	// 	fmt.Println("Error making POST request:", err)
// 	// 	return
// 	// }
// 	// defer resp.Body.Close()
// 	// fmt.Println("KI KI KI RRRR:", resp)
// 	// if resp.StatusCode != http.StatusOK {
// 	// 	fmt.Println("Received non-OK response status code:", resp.StatusCode)
// 	// 	return
// 	// }

// 	// var responseBody []byte
// 	// _, err = resp.Body.Read(responseBody)
// 	// if err != nil {
// 	// 	fmt.Println("Error reading response body:", err)
// 	// 	return
// 	// }
// 	// fmt.Println("Response:", string(responseBody))

// }
