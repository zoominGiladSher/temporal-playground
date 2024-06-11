package main

import (
	"log"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	gtp "go-temporal-playground"
)

func main() {
	temporalClient, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer temporalClient.Close()
	temporalWorker := worker.New(temporalClient, gtp.WorkflowTaskQueue, worker.Options{})
	registerActivityOptions := activity.RegisterOptions{
		Name: "Testing Activity",
	}
	temporalWorker.RegisterActivityWithOptions(gtp.TestingActivity, registerActivityOptions)
	registerWorkflowOptions := workflow.RegisterOptions{
		Name: gtp.WorkflowId,
	}
	temporalWorker.RegisterWorkflowWithOptions(gtp.TestingWorkflow, registerWorkflowOptions)
	err = temporalWorker.Run(worker.InterruptCh())
	if err != nil {
		panic(err)
	}
}
