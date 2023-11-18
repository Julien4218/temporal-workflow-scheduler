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

func EventWorkflow(ctx workflow.Context, input *EventWorkflowInput) error {
	workflowStartTimestamp := timestamppb.Now()
	workflowStarttime := workflowStartTimestamp.AsTime()
	nowAsMinute := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond).Truncate(time.Nanosecond)
	timeInterval := nowAsMinute.Sub(workflowStarttime)

	ctx = updateWorkflowContextOptions(ctx)
	logrus.Infof("%s-EventWorkflow started", instrumentation.Hostname)
	defer logrus.Infof("%s-EventWorkflow completed", instrumentation.Hostname)
	txn := instrumentation.NrApp.StartTransaction("EventWorkflow")
	defer txn.End()

	logrus.Infof(fmt.Sprintf("Got input:%s", input.TimestampStart.String()))

	accountID, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if err != nil {
		logrus.Errorf(fmt.Sprintf("error could not find environment variable NEW_RELIC_ACCOUNT_ID"))
		return err
	}
	eventInput := &newrelicActivities.CreateEventInput{
		AccountID:     accountID,
		EventDataJson: fmt.Sprintf("{\"eventType\":\"temporalWorkflowScheduler\", \"timestamp\":\"%s\", \"timeIntervalMs\":\"%s\"}", workflowStarttime.Unix(), timeInterval.Milliseconds()),
	}

	var result interface{}
	if err := workflow.ExecuteActivity(ctx, newrelicActivities.NewCreateEventActivity, eventInput).Get(ctx, &result); err != nil {
		logrus.Errorf("Activity NewCreateEventActivity failed. Error: %s", err)
		return err
	}

	logrus.Infof("time difference is %s", timeInterval)

	return nil
}

func updateWorkflowContextOptions(ctx workflow.Context) workflow.Context {
	// Define the SlackMessageActivity Execution options
	// StartToCloseTimeout or ScheduleToCloseTimeout must be set
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	return ctx
}
