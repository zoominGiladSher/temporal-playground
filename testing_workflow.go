package gotemporalplayground

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

type WorkflowParam struct {
	X string
	Y int
}

type WorkflowResult struct {
	X string `json:"x"`
	Y int    `json:"y"`
}

type ActivityParam struct {
	X string
	Y int
}

type ActivityResult struct {
	X string `json:"x"`
	Y int    `json:"y"`
}

type TestingSignal struct {
	Message string `json:"message"`
}

func TestingWorkflow(ctx workflow.Context, param WorkflowParam) (*WorkflowResult, error) {
	var signal TestingSignal
	currentState := STATE_RUNNING
	signalChan := workflow.GetSignalChannel(ctx, SignalName)
	err := workflow.SetQueryHandler(ctx, QueryTypeCurrentState, func() (string, error) {
		return currentState, nil
	})
	if err != nil {
		currentState = STATE_FAILED_TO_REGISTER_QUERY_HANDLER
		return nil, err
	}
	// Run some stuff while waiting for signal
	// workflow.Go(ctx, func(ctx workflow.Context) {
	// 	for {
	// 		selector := workflow.NewSelector(ctx)
	// 		selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
	// 			c.Receive(ctx, &signal)
	// 		})
	// 		selector.Select(ctx)
	// 	}
	// })
	// block until signal received
	currentState = STATE_WAITING_SIGNAL
	signalChan.Receive(ctx, &signal)
	if len(signal.Message) > 0 && signal.Message == SignalAbort {
		currentState = STATE_ABORTED
		return nil, errors.New("workflow aborted")
	}

	var activityResult *ActivityResult
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	})
	currentState = STATE_WAITING_ACTIVITY
	workflow.ExecuteActivity(ctx, TestingActivity, ActivityParam{X: param.X, Y: param.Y}).Get(ctx, &activityResult)

	workflowResult := &WorkflowResult{X: activityResult.X, Y: activityResult.Y}
	currentState = STATE_WAITING_TIMER
	workflow.Sleep(ctx, 3*time.Second)
	currentState = STATE_DONE
	return workflowResult, nil
}

func TestingActivity(ctx context.Context, param ActivityParam) (*ActivityResult, error) {
	time.Sleep(2 * time.Second)
	_, err := http.Get(fmt.Sprintf("http://localhost:8091/resume?taskToken=%s", base64.URLEncoding.EncodeToString(activity.GetInfo(ctx).TaskToken)))
	if err != nil {
		return nil, err
	}

	return nil, activity.ErrResultPending
}

const (
	WorkflowId                             = "TestingWorkflow"
	WorkflowTaskQueue                      = "Testing Queue"
	SignalName                             = "Testing Signal"
	SignalAbort                            = "abort"
	QueryTypeCurrentState                  = "current_state"
	STATE_RUNNING                          = "running"
	STATE_ABORTED                          = "aborted"
	STATE_FAILED_TO_REGISTER_QUERY_HANDLER = "failed to register query handler"
	STATE_WAITING_SIGNAL                   = "waiting for signal"
	STATE_WAITING_TIMER                    = "waiting for timer"
	STATE_WAITING_ACTIVITY                 = "waiting for activity"
	STATE_DONE                             = "done"
	STATE_NOT_RUNNING                      = "not running"
	STATE_UNKNOWN                          = "unknown"
)
