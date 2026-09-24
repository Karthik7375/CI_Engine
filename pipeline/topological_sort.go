package pipeline

import (
	"cicd/model"
	"cicd/runner"
	"context"
	"errors"
	"fmt"
)

func Topological_sort_jobs(pipeline map[string]model.Job) ([]string, error) {
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize inDegree for all jobs
	for name := range pipeline {
		inDegree[name] = 0
	}

	for name, job := range pipeline {
		for _, dep := range job.Needs {
			graph[dep] = append(graph[dep], name)
			inDegree[name]++
		}
	}

	// for name, job := range graph {
	// 	fmt.Println(name, " => ", job)
	// }

	// Queue for jobs with no dependencies
	var queue []string
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}
	// Result slice for the sorted order
	var result []string

	for len(queue) > 0 {
		// Dequeue a job
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Reduce in-degree for all jobs that depend on current
		for _, neighbor := range graph[current] {
			inDegree[neighbor]--
			fmt.Printf("Neighbour to %s => %s\n`", current, neighbor)
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check for cycles (if result doesn't include all jobs)
	if len(result) != len(pipeline) {
		return nil, fmt.Errorf("cycle detected in job dependencies")
	}

	return result, nil
}

func PipelineExecution(ctx context.Context, r runner.Runner, pipeline map[string]model.Job) error {
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize inDegree for all jobs
	for name := range pipeline {
		inDegree[name] = 0
	}

	for name, job := range pipeline {
		for _, dep := range job.Needs {
			graph[dep] = append(graph[dep], name)
			inDegree[name]++
		}
	}

	for name, job := range graph {
		fmt.Println(name, " => ", job)
	}

	finished := make(chan model.JobResult)

	totalJobs := len(pipeline)
	completed := 0

	startJob := func(name string) {
		go func() {
			result, err := r.Run(ctx, name, pipeline[name])
			if err != nil {
				fmt.Printf("result: %v\n", result)
			}
			finished <- result
		}()
	}
	// Start all root jobs (parallel)
	for name, degree := range inDegree {
		if degree == 0 {
			startJob(name)
		}
	}
	for completed < totalJobs {

		select {
		case <-ctx.Done():
			return ctx.Err()

		case result := <-finished:

			if result.Err != "" {
				return errors.New(result.Err)
			}

			completed++

			// Unblock dependent jobs
			for _, dependent := range graph[result.Name] {
				inDegree[dependent]--

				if inDegree[dependent] == 0 {
					startJob(dependent)
				}
			}
		}
	}
	return nil
}
