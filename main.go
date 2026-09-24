package main

import (
	"cicd/model"
	"cicd/pipeline"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Minute,
	)
	defer cancel()

	data, err := os.ReadFile("yaml_test_files/pipeline.yaml")
	if err != nil {
		panic(err)
	}
	var p model.Pipeline

	if err := yaml.Unmarshal(data, &p); err != nil {
		panic(err)
	}
	fmt.Println(p)
	// result := pipeline.ValidatePipeline(p)
	// fmt.Println(result)

	sortedJobs, err := pipeline.Topological_sort_jobs(p.Pipeline)
	// pipeline.PipelineExecution(ctx, r, p.Pipeline)
	// if err != nil {
	// 	r
	// 	log.Fatalf("Error sorting jobs: %v", err)
	// }
	fmt.Println("Execution order:", sortedJobs)
	// for key, value := range sortedJobs {
	// 	fmt.Println(key, value)
	// }

	var wg sync.WaitGroup
	finished := make(chan model.JobResult)
	for name, job := range p.Pipeline {
		wg.Add(1)

		go func(name string, job model.Job) {
			defer wg.Done()

			pipeline.ExecuteJob(ctx, name, job, finished)
		}(name, job)
	}
	go func() {
		wg.Wait()
		close(finished)
	}()

	for result := range finished {
		fmt.Println(result)
	}
}
