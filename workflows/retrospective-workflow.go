package workflows

import (
	"errors"
	"os"
	"time"

	slackActivities "github.com/Julien4218/temporal-slack-activity/activities"
	slackModels "github.com/Julien4218/temporal-slack-activity/models"
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/workflow"

	"github.com/Julien4218/temporal-workflow-scheduler/instrumentation"
)

const RetrospectiveWorkflowName = "RetrospectiveWorkflow"

type RetrospectiveWorkflowInput struct {
}

func RetrospectiveWorkflow(ctx workflow.Context, input *RetrospectiveWorkflowInput) (string, error) {
	logrus.Infof("%s-SlackWorkflow started:%s", instrumentation.Hostname, RetrospectiveWorkflowName)
	defer logrus.Infof("%s-SlackWorkflow completed:%s", instrumentation.Hostname, RetrospectiveWorkflowName)
	txn := instrumentation.NrApp.StartTransaction(RetrospectiveWorkflowName)
	defer txn.End()

	ctx = updateIncidentWorkflowContextOptions(ctx, 10*time.Minute)
	logrus.Infof("Got retrospective input:%s", input)

	slackChannel := os.Getenv("SLACK_CHANNEL")
	if slackChannel == "" {
		return "", errors.New("required environment variable SLACK_CHANNEL is not set")
	}
	message := "Starting retrospective"
	messageData := slackModels.SlackActivityData{
		ChannelId:            slackChannel,
		FirstResponseWarning: message,
	}
	var result slackModels.MessageDetails
	if err := workflow.ExecuteActivity(ctx, slackActivities.PostMessageActivity, messageData).Get(ctx, &result); err != nil {
		logrus.Errorf("Activity failed. Error: %s", err)
		return "", err
	}
	logrus.Infof("%s workflow completed.", RetrospectiveWorkflowName)

	return "DONE", nil
}
