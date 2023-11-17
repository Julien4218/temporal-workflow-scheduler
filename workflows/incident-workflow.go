package workflows

import (
	"errors"
	"os"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/sirupsen/logrus"

	slackActivities "github.com/Julien4218/temporal-slack-activity/activities"
	slackModels "github.com/Julien4218/temporal-slack-activity/models"

	"github.com/Julien4218/temporal-workflow-scheduler/instrumentation"
)

const IncidentWorkflowName = "IncidentWorkflow"

type IncidentWorkflowInput struct {
}

func IncidentWorkflow(ctx workflow.Context, input *IncidentWorkflowInput) (string, error) {
	logrus.Infof("%s-SlackWorkflow started:%s", instrumentation.Hostname, IncidentWorkflowName)
	defer logrus.Infof("%s-SlackWorkflow completed:%s", instrumentation.Hostname, IncidentWorkflowName)
	txn := instrumentation.NrApp.StartTransaction(IncidentWorkflowName)
	defer txn.End()

	ctx = updateIncidentWorkflowContextOptions(ctx, 10*time.Minute)
	logrus.Infof("Got incident input:%s", input)

	slackChannel := os.Getenv("SLACK_CHANNEL")
	if slackChannel == "" {
		return "", errors.New("required environment variable SLACK_CHANNEL is not set")
	}
	message := "Starting an Incident in Upboard"
	messageData := slackModels.SlackActivityData{
		ChannelId:            slackChannel,
		FirstResponseWarning: message,
	}
	var result slackModels.MessageDetails
	if err := workflow.ExecuteActivity(ctx, slackActivities.PostMessageActivity, messageData).Get(ctx, &result); err != nil {
		logrus.Errorf("Activity failed. Error: %s", err)
		return "", err
	}

	logrus.Infof("%s workflow completed.", IncidentWorkflowName)

	return "DONE", nil
}

func updateIncidentWorkflowContextOptions(ctx workflow.Context, startToCloseTimeout time.Duration) workflow.Context {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: startToCloseTimeout,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	return ctx
}
