package workflows

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.temporal.io/sdk/workflow"

	"github.com/sirupsen/logrus"

	"github.com/Julien4218/temporal-workflow-scheduler/instrumentation"
)

const WorkflowName = "ErWorkflow"

type EventWorkflowInput struct {
	TimestampStart timestamp.Timestamp
}

func EventWorkflow(ctx workflow.Context, input *EventWorkflowInput) error {
	workflowStartTimestamp := *timestamppb.Now()
	workflowStarttime := workflowStartTimestamp.AsTime()

	nowAsMinute := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond).Truncate(time.Nanosecond)
	timeInterval := nowAsMinute.Sub(workflowStarttime)

	logrus.Infof("%s-EventWorkflow started", instrumentation.Hostname)
	defer logrus.Infof("%s-EventWorkflow completed", instrumentation.Hostname)
	txn := instrumentation.NrApp.StartTransaction("EventWorkflow")
	defer txn.End()

	// ctx = updateWorkflowContextOptions(ctx)
	logrus.Infof(fmt.Sprintf("Got input:%s", input.TimestampStart.String()))

	// eventInput := &newrelicActivities.CreateEventInput{
	// 	AccountID:     0,
	// 	EventDataJson: "",
	// }

	// var result interface{}
	// if err := workflow.ExecuteActivity(ctx, newrelicActivities.NewCreateEventActivity, input).Get(ctx, &result); err != nil {
	// 	logrus.Errorf("Activity NewCreateEventActivity failed. Error: %s", err)
	// 	return err
	// }

	logrus.Infof("time difference is %s", timeInterval)

	return nil
}

// func updateWorkflowContextOptions(ctx workflow.Context) workflow.Context {
// 	// Define the SlackMessageActivity Execution options
// 	// StartToCloseTimeout or ScheduleToCloseTimeout must be set
// 	activityOptions := workflow.ActivityOptions{
// 		StartToCloseTimeout: 10 * time.Second,
// 	}
// 	ctx = workflow.WithActivityOptions(ctx, activityOptions)
// 	return ctx
// }
