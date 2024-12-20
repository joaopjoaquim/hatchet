package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/hatchet-dev/hatchet/pkg/client"
	"github.com/hatchet-dev/hatchet/pkg/client/compute"
	"github.com/hatchet-dev/hatchet/pkg/cmdutils"
	"github.com/hatchet-dev/hatchet/pkg/worker"
)

type userCreateEvent struct {
	Username string            `json:"username"`
	UserID   string            `json:"user_id"`
	Data     map[string]string `json:"data"`
}

type stepOneOutput struct {
	Message string `json:"message"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	events := make(chan string, 50)
	interrupt := cmdutils.InterruptChan()

	cleanup, err := run(events)
	if err != nil {
		panic(err)
	}

	<-interrupt

	if err := cleanup(); err != nil {
		panic(fmt.Errorf("error cleaning up: %w", err))
	}
}

func run(events chan<- string) (func() error, error) {
	c, err := client.New()

	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	w, err := worker.NewWorker(
		worker.WithClient(
			c,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating worker: %w", err)
	}

	pool := "test-pool"
	basicCompute := compute.Compute{
		Pool:        &pool,
		NumReplicas: 1,
		CPUs:        1,
		MemoryMB:    1024,
		CPUKind:     compute.ComputeKindSharedCPU,
		Regions:     []compute.Region{compute.Region("ewr")},
	}

	performancePool := "performance-pool"
	performanceCompute := compute.Compute{
		Pool:        &performancePool,
		NumReplicas: 1,
		CPUs:        2,
		MemoryMB:    1024,
		CPUKind:     compute.ComputeKindPerformanceCPU,
		Regions:     []compute.Region{compute.Region("ewr")},
	}

	err = w.RegisterWorkflow(
		&worker.WorkflowJob{
			On:          worker.Events("user:create:simple"),
			Name:        "simple",
			Description: "This runs after an update to the user model.",
			Concurrency: worker.Expression("input.user_id"),
			Steps: []*worker.WorkflowStep{
				worker.Fn(func(ctx worker.HatchetContext) (result *stepOneOutput, err error) {
					input := &userCreateEvent{}

					err = ctx.WorkflowInput(input)

					if err != nil {
						return nil, err
					}

					log.Printf("step-one")
					events <- "step-one"

					return &stepOneOutput{
						Message: "Username is: " + input.Username,
					}, nil
				},
				).SetName("step-one").SetCompute(&basicCompute),
				worker.Fn(func(ctx worker.HatchetContext) (result *stepOneOutput, err error) {
					input := &stepOneOutput{}
					err = ctx.StepOutput("step-one", input)

					if err != nil {
						return nil, err
					}

					log.Printf("step-two")
					events <- "step-two"

					return &stepOneOutput{
						Message: "Above message is: " + input.Message,
					}, nil
				}).SetName("step-two").AddParents("step-one").SetCompute(&performanceCompute),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error registering workflow: %w", err)
	}

	go func() {
		testEvent := userCreateEvent{
			Username: "echo-test",
			UserID:   "1234",
			Data: map[string]string{
				"test": "test",
			},
		}

		log.Printf("pushing event user:create:simple")
		// push an event
		err := c.Event().Push(
			context.Background(),
			"user:create:simple",
			testEvent,
			client.WithEventMetadata(map[string]string{
				"hello": "world",
			}),
		)
		if err != nil {
			panic(fmt.Errorf("error pushing event: %w", err))
		}
	}()

	cleanup, err := w.Start()
	if err != nil {
		panic(err)
	}

	return cleanup, nil
}
