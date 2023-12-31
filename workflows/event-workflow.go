package workflows

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.temporal.io/sdk/workflow"

	"github.com/sirupsen/logrus"

	newrelicActivities "github.com/Julien4218/temporal-newrelic-activity/activities"
	"github.com/Julien4218/temporal-workflow-scheduler/instrumentation"
)

const WorkflowName = "ErWorkflow"

type EventWorkflowInput struct {
	TimestampStart timestamp.Timestamp
}

func EventWorkflow(ctx workflow.Context, input *EventWorkflowInput) (string, error) {
	workflowStarttime := timestamppb.Now().AsTime()
	nowAsMinute := timestamppb.New(time.Date(
		workflowStarttime.Year(),
		workflowStarttime.Month(),
		workflowStarttime.Day(),
		workflowStarttime.Hour(),
		workflowStarttime.Minute(),
		0, 0, time.UTC))
	timeIntervalMs := workflowStarttime.Sub(nowAsMinute.AsTime()).Milliseconds()

	ctx = updateWorkflowContextOptions(ctx)
	logrus.Infof("%s-EventWorkflow started", instrumentation.Hostname)
	defer logrus.Infof("%s-EventWorkflow completed", instrumentation.Hostname)
	txn := instrumentation.NrApp.StartTransaction("EventWorkflow")
	defer txn.End()

	logrus.Infof(fmt.Sprintf("Got input:%s", input.TimestampStart.String()))

	accountID, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if err != nil {
		logrus.Errorf(fmt.Sprintf("error could not find environment variable NEW_RELIC_ACCOUNT_ID"))
		return "", err
	}
	eventInput := &newrelicActivities.CreateEventInput{
		AccountID:     accountID,
		EventDataJson: fmt.Sprintf("{\"eventType\":\"temporalWorkflowScheduler\", \"timestamp\":%d, \"timeIntervalMs\":%d}", workflowStarttime.Unix(), timeIntervalMs),
	}

	var result interface{}
	if err := workflow.ExecuteActivity(ctx, newrelicActivities.NewCreateEventActivity, eventInput).Get(ctx, &result); err != nil {
		logrus.Errorf("Activity NewCreateEventActivity failed. Error: %s", err)
		return "", err
	}

	message := fmt.Sprintf("time difference is %d", timeIntervalMs)
	logrus.Infof(message)

	info := workflow.GetInfo(ctx)
	runID := info.WorkflowExecution.RunID
	output := fmt.Sprintf("%s,%s,%d,%d", runID, workflowStarttime, workflowStarttime.Unix(), timeIntervalMs)
	directory := fmt.Sprintf("%d/%d/%d/%d/%d",
		workflowStarttime.Year(), workflowStarttime.Month(), workflowStarttime.Day(),
		workflowStarttime.Hour(), workflowStarttime.Minute())
	appendFile(directory, "output.log", output)
	return message, nil
}

func appendFile(directory string, filename string, message string) {
	os.MkdirAll(directory, 0700)
	filepath := fmt.Sprintf("%s/%s", directory, filename)
	f, err := os.OpenFile(filepath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errorMessage := fmt.Sprintf("error while opening file:%s", err.Error())
		logrus.Infof(errorMessage)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintln(message)); err != nil {
		errorMessage := fmt.Sprintf("error while appending to file:%s", err.Error())
		logrus.Infof(errorMessage)
	}
}

func updateWorkflowContextOptions(ctx workflow.Context) workflow.Context {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	return ctx
}
