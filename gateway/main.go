package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	gtp "go-temporal-playground"

	"go.temporal.io/sdk/client"
)

func main() {
	temporalClient, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer temporalClient.Close()

	http.HandleFunc("/start", CreateStartHandler(temporalClient))
	http.HandleFunc("/resume", CreateResumeHandler(temporalClient))
	http.HandleFunc("/signal", CreateSignalHandler(temporalClient))
	http.HandleFunc("/abort", CreateAbortHandler(temporalClient))
	http.HandleFunc("/status", CreateStatusHandler(temporalClient))
	err = http.ListenAndServe(":8091", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func CreateStartHandler(temporalClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		workflowOptions := client.StartWorkflowOptions{
			TaskQueue: gtp.WorkflowTaskQueue,
			ID:        gtp.WorkflowId,
		}
		workflowResult, err := temporalClient.ExecuteWorkflow(ctx, workflowOptions, gtp.TestingWorkflow, gtp.WorkflowParam{
			X: "Hello",
			Y: 42,
		})
		if err != nil {
			log.Println("Error starting workflow", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var result gtp.WorkflowResult
		workflowResult.Get(ctx, &result)
		resp, err := json.Marshal(result)
		if err != nil {
			log.Println("Error marshalling workflow result", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func CreateResumeHandler(temporalClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskToken, err := base64.URLEncoding.DecodeString(r.URL.Query().Get("taskToken"))
		if err != nil || len(taskToken) == 0 {
			log.Println("Error decoding task token", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		activityResult := gtp.ActivityResult{
			X: "WOW",
			Y: 69,
		}
		err = temporalClient.CompleteActivity(ctx, taskToken, activityResult, err)
		if err != nil {
			log.Println("Error completing activity", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(activityResult)
		if err != nil {
			log.Println("Error marshalling activity result", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func CreateSignalHandler(temporalClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signal := gtp.TestingSignal{
			Message: "Testing Signal",
		}

		ctx := r.Context()
		err := temporalClient.SignalWorkflow(ctx, gtp.WorkflowId, "", gtp.SignalName, signal)
		if err != nil {
			log.Println("Error signaling workflow", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal(signal)
		if err != nil {
			log.Println("Error marshalling signal", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	}
}

func CreateAbortHandler(temporalClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := temporalClient.CancelWorkflow(ctx, gtp.WorkflowId, "")
		if err != nil {
			log.Println("Error aborting workflow", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func CreateStatusHandler(temporalClient client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		status, err := temporalClient.QueryWorkflow(ctx, gtp.WorkflowId, "", gtp.QueryTypeCurrentState, nil)
		if err != nil {
			log.Println("Error querying workflow", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal(status)
		if err != nil {
			log.Println("Error marshalling status", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	}
}
