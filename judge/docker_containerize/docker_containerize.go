package docker_containerize

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gofiber/fiber/v2"
)

type SubmittedCode struct {
	Code string `json:"code"`
}

// func Containerize(c *fiber.Ctx) error {

// 	// Construct the absolute path to the code.py file
// 	// codeFilePath := "./docker_containerize/code.py"

// 	// Create or truncate the code.py file
// 	// file, err := os.Create(codeFilePath)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// defer file.Close()

// 	// Write the user-submitted code from the request body to the code.py file
// 	// err = os.WriteFile(codeFilePath, []byte(submittedCode.Code), 0644)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// Call the Containerize function
// 	err = handler(c)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func Containerize(c *fiber.Ctx) error {
	fmt.Println("INSIDE OOOOOO!!!!!!")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	// Parse the JSON body
	var submittedCode SubmittedCode
	if err := c.BodyParser(&submittedCode); err != nil {
		return err
	}

	// Create a tar archive from the build context directory
	buildContextDir := "./docker_containerize/Dockerfile.unknown" // Replace with the path to the directory containing your Dockerfile
	// codeFilePath := "./docker_containerize/code.py"            // Replace with the path to the code file
	testCasesFilePath := "./docker_containerize/test_cases.py" // Replace with the path to the test cases file
	buildContextTarReader, err := createTarArchive(buildContextDir, submittedCode.Code, testCasesFilePath)
	if err != nil {
		return err
	}

	imageName := "python_execute:v1"
	functionName := "add" // Replace with the name of the function to execute (To be fetched from the database or the request body)
	buildArgs := map[string]*string{
		"FUNCTION_NAME": &functionName, // Use the FUNCTION_NAME environment variable
	}

	// Build the Docker image from the Dockerfile
	buildResponse, err := cli.ImageBuild(ctx, buildContextTarReader, types.ImageBuildOptions{Tags: []string{imageName}, BuildArgs: buildArgs})
	if err != nil {
		return err
	}
	defer buildResponse.Body.Close()

	// Print the build output
	_, err = io.Copy(os.Stdout, buildResponse.Body)
	if err != nil {
		return err
	}

	// Create a new container
	config := &container.Config{
		Image: imageName, // Replace with the name of the built image
	}
	containerResponse, err := cli.ContainerCreate(ctx, config, nil, nil, nil, "")
	if err != nil {
		return err
	}

	// Start the container
	err = cli.ContainerStart(ctx, containerResponse.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(ctx, containerResponse.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("Container exited with non-zero status code: %d", status.StatusCode)
		}
	}

	// Retrieve the container logs
	var logBuffer bytes.Buffer
	out, err := cli.ContainerLogs(ctx, containerResponse.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the container logs to the buffer
	_, err = stdcopy.StdCopy(&logBuffer, &logBuffer, out)
	if err != nil {
		return err
	}

	// Return the log output
	c.SendString(logBuffer.String())

	// Close and remove the container
	err = cli.ContainerRemove(ctx, containerResponse.ID, container.RemoveOptions{Force: true})
	if err != nil {
		return err
	}

	return nil
}

func createTarArchive(dockerfilePath string, codeContent string, testCasesFilePath string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	// Add the Dockerfile to the tar archive
	dockerfileData, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return nil, err
	}
	dockerfileHeader := &tar.Header{
		Name: "Dockerfile",
		Size: int64(len(dockerfileData)),
		Mode: 0644,
	}
	if err := tw.WriteHeader(dockerfileHeader); err != nil {
		return nil, err
	}
	if _, err := tw.Write(dockerfileData); err != nil {
		return nil, err
	}

	// Add the code.py file to the tar archive
	// codeFileData, err := os.ReadFile(codeFilePath)
	// if err != nil {
	// 	return nil, err
	// }
	codeFileHeader := &tar.Header{
		Name: "code.py",
		Size: int64(len(codeContent)),
		Mode: 0644,
	}
	if err := tw.WriteHeader(codeFileHeader); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(codeContent)); err != nil {
		return nil, err
	}

	// Add the test_cases.py file to the tar archive
	testCasesFileData, err := os.ReadFile(testCasesFilePath)
	if err != nil {
		return nil, err
	}
	testCasesFileHeader := &tar.Header{
		Name: "test_cases.py",
		Size: int64(len(testCasesFileData)),
		Mode: 0644,
	}
	if err := tw.WriteHeader(testCasesFileHeader); err != nil {
		return nil, err
	}
	if _, err := tw.Write(testCasesFileData); err != nil {
		return nil, err
	}

	return &buf, nil
}
