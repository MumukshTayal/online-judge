package docker_containerize

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
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
	Code        string `json:"test_code"`
	TestsInput  string `json:"test_input"`
	TestsOutput string `json:"test_output"`
	Language    string `json:"language"`
}

func Containerize(c *fiber.Ctx) (string, error) {
	// fmt.Println("INSIDE OOOOOO!!!!!!")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	// Parse the JSON body
	var submittedCode SubmittedCode
	if err := c.BodyParser(&submittedCode); err != nil {
		return "", err
	}

	// fmt.Println("CODE:", submittedCode.Code)
	// fmt.Println("TESTS INPUT:", submittedCode.TestsInput)
	// fmt.Println("TESTS OUTPUT:", submittedCode.TestsOutput)
	// fmt.Println("LANGUAGE:", submittedCode.Language)

	// Create a tar archive from the build context directory
	var buildContextDir string
	switch submittedCode.Language {
	case "py":
		buildContextDir = "./docker_containerize/Dockerfile.python"
	case "cpp":
		buildContextDir = "./docker_containerize/Dockerfile.cplpl"
	case "c":
		buildContextDir = "./docker_containerize/Dockerfile.cplpl"
	case "java":
		buildContextDir = "./docker_containerize/Dockerfile.jv"
	default:
		return "", errors.New("unsupported language here")
	}
	// codeFilePath := "./docker_containerize/code.py"            // Replace with the path to the code file
	testCasesFilePath := "./docker_containerize/test_cases.py" // Replace with the path to the test cases file
	language := submittedCode.Language                         // Add a Language field to the SubmittedCode struct
	buildContextTarReader, err := createTarArchive(buildContextDir, submittedCode.Code, testCasesFilePath, submittedCode.TestsInput, submittedCode.TestsOutput, submittedCode.Language)
	if err != nil {
		return "", err
	}

	imageName := "python_execute:v1"
	// functionName := "add" // Replace with the name of the function to execute (To be fetched from the database or the request body)
	buildArgs := map[string]*string{
		"LANG": &language, // Use the FUNCTION_NAME environment variable
	}

	// Build the Docker image from the Dockerfile
	buildResponse, err := cli.ImageBuild(ctx, buildContextTarReader, types.ImageBuildOptions{Tags: []string{imageName}, BuildArgs: buildArgs})
	if err != nil {
		return "", err
	}
	defer buildResponse.Body.Close()

	// Print the build output
	_, err = io.Copy(os.Stdout, buildResponse.Body)
	if err != nil {
		return "", err
	}

	// Create a new container
	config := &container.Config{
		Image: imageName, // Replace with the name of the built image
		// Env: []string{
		// 	"LANGUAGE=" + language, // Set the LANGUAGE environment variable
		// },
	}
	containerResponse, err := cli.ContainerCreate(ctx, config, nil, nil, nil, "")
	if err != nil {
		return "", err
	}
	// fmt.Println(containerResponse)
	// Start the container
	err = cli.ContainerStart(ctx, containerResponse.ID, container.StartOptions{})
	if err != nil {
		return "", err
	}

	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(ctx, containerResponse.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return "", fmt.Errorf("Container exited with non-zero status code: %d", status.StatusCode)
		}
	}

	// Retrieve the container logs
	var logBuffer bytes.Buffer
	out, err := cli.ContainerLogs(ctx, containerResponse.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy the container logs to the buffer
	_, err = stdcopy.StdCopy(&logBuffer, &logBuffer, out)
	if err != nil {
		return "", err
	}

	// Return the log output
	// c.SendString(logBuffer.String())

	// Close and remove the container
	err = cli.ContainerRemove(ctx, containerResponse.ID, container.RemoveOptions{Force: true})
	if err != nil {
		return "", err
	}

	return logBuffer.String(), nil
}

func createTarArchive(dockerfilePath string, codeContent string, testCasesFilePath string, testcasesInput string, testcasesOutput string, language string) (io.Reader, error) {
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
	// Create a file for the code
	var filename string
	if language == "py" {
		filename = "pytcode.py"
	} else if language == "cpp" || language == "c" {
		filename = "ccode.txt"
	} else if language == "java" {
		filename = "javacode.java"
	} else {
		return nil, errors.New("unsupported language")
	}

	codeFileHeader := &tar.Header{
		Name: filename,
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
	if _, err := tw.Write([]byte(testCasesFileData)); err != nil {
		return nil, err
	}

	testCasesInputFileHeader := &tar.Header{
		Name: "input.txt",
		Size: int64(len(testcasesInput)),
		Mode: 0644,
	}
	if err := tw.WriteHeader(testCasesInputFileHeader); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(testcasesInput)); err != nil {
		return nil, err
	}

	testCasesOutputFileHeader := &tar.Header{
		Name: "output.txt",
		Size: int64(len(testcasesOutput)),
		Mode: 0644,
	}
	if err := tw.WriteHeader(testCasesOutputFileHeader); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(testcasesOutput)); err != nil {
		return nil, err
	}

	return &buf, nil
}
