package pipeline

import (
	"cicd/model"
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

func ExecuteJob(ctx context.Context, name string, job model.Job, finished chan<- model.JobResult) error {
	fmt.Println("execute job")
	lang, args, err := resolveBuild(job)
	if err != nil {
		return nil
	}
	fmt.Println(time.Now())
	curr, err := os.Getwd()
	if err != nil {
		return nil
	}

	logs, done := StreamBuild(ctx, curr, lang, args...)
	// StreamBuild(ctx, ".", "echo", "hello")

	for log := range logs {
		fmt.Printf("[%s][%s] %s\n", job.Name, log.Stream, log.Line)
	}

	if err := <-done; err != nil {
		fmt.Printf("[%s] Failed: %v\n", job.Name, err)
	} else {
		fmt.Printf("[%s] Succeeded\n", job.Name)
	}
	finished <- model.JobResult{
		Name: name,
		Err:  "",
	}

	fmt.Println("hello from the end")
	time.Sleep(3 * time.Second)
	return nil
}

var wg sync.WaitGroup
var mu sync.Mutex

func goroutine_main(number *int) {
	defer wg.Done()

	mu.Lock()
	defer mu.Unlock()
	fmt.Println(*number)
	*number++
}

func worker(ch chan string) {
	ch <- "Hello!"
}
func main() {
	// wg.Add(2)
	// n := 1
	// go goroutine_main(&n)
	// go goroutine_main(&n)

	// wg.Wait()

	// ch := make(chan string)

	// go worker(ch)

	// msg := <-ch

	// fmt.Println(msg)
	// start()
}
