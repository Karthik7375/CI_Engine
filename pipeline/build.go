package pipeline

import (
	"bufio"
	"cicd/model"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func buildCommand(buildType string) []string {
	switch buildType {
	case "rust":
		return []string{"cargo", "build"}

	case "go":
		return []string{"go", "build"}

		// case "node":
		// 	return []string{"pnpm", "run", "build"}
	}

	return nil
}

type CIResult struct {
	Status   string
	Step     string
	ExitCode int
	Logs     string
}

func RunStep(step string, command string, args ...string) CIResult {
	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()

	result := CIResult{
		Step: step,
		Logs: strings.TrimSpace(string(output)),
	}

	if err != nil {
		result.Status = "failed"

		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}

		return result
	}

	result.Status = "success"
	result.ExitCode = 0

	return result
}

func resolveBuild(p model.Job) (string, []string, error) {
	switch p.Language {

	case "go":
		return "go",
			[]string{"build", "."},
			nil

	case "rust":
		return "cargo",
			[]string{"build"},
			nil

	// case "node":
	// 	return "pnpm",
	// 		[]string{"run", "build"},
	// 		nil

	default:
		return "", nil,
			fmt.Errorf(
				"unknown build type: %s",
				p.Language,
			)
	}
}

// func DisplayCIResult(r CIResult) {
// 	if r.Status == "success" {
// 		fmt.Println("✅ BUILD SUCCESS")
// 	} else {
// 		fmt.Println("❌ BUILD FAILED")
// 	}

// 	fmt.Printf("Step      : %s\n", r.Step)
// 	fmt.Printf("Exit Code : %d\n", r.ExitCode)

// 	fmt.Println("----------------------------------------")

// 	if r.Logs == "" {
// 		fmt.Println("(no logs)")
// 		return
// 	}

// 	fmt.Println("Logs:")
// 	fmt.Println(r.Logs)
// }

var count = 0

func StreamBuild(ctx context.Context, dir string, lang string, args ...string) (<-chan model.LogEvent, <-chan error) {
	fmt.Println()
	logs := make(chan model.LogEvent, 100)
	done := make(chan error, 1)

	go func() {
		defer close(logs)
		defer close(done)

		cmd := exec.CommandContext(ctx, lang, args...)
		// cmd.Dir = dir

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			done <- err
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			done <- err
			return
		}

		if err := cmd.Start(); err != nil {
			done <- err
			return
		}

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			scanner := bufio.NewScanner(stdout)
			println("inside scanner of stdout")
			for scanner.Scan() {
				logs <- model.LogEvent{
					Stream: "stdout",
					Line:   scanner.Text(),
				}
			}
			if scanner.Err() != nil {
				fmt.Println()
			}
		}()

		go func() {
			defer wg.Done()

			scanner := bufio.NewScanner(stderr)
			println("inside scanner of stderr")

			for scanner.Scan() {
				logs <- model.LogEvent{
					Stream: "stderr",
					Line:   scanner.Text(),
				}
			}
			if scanner.Err() != nil {
				fmt.Println()
			}
		}()

		wg.Wait()

		done <- cmd.Wait()
	}()

	return logs, done

}
