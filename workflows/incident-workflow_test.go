package workflows

import (
	"context"
	"errors"

	slackActivities "github.com/Julien4218/temporal-slack-activity/activities"
	"github.com/Julien4218/temporal-slack-activity/models"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/temporal"
)

func (s *UnitTestSuite) Test_IncidentWorkflow_PostMessageFailed() {
	s.Suite.T().Setenv("SLACK_CHANNEL", "my slack channel")

	s.env.OnActivity(slackActivities.PostMessageActivity, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, slackDate models.SlackActivityData) (models.MessageDetails, error) {
			// your mock function implementation
			return models.MessageDetails{}, errors.New("Slack PostMessage failed")
		})
	s.env.ExecuteWorkflow(IncidentWorkflow, nil)

	s.True(s.env.IsWorkflowCompleted())

	err := s.env.GetWorkflowError()
	s.Error(err)
	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("Slack PostMessage failed", applicationErr.Error())
}

func (s *UnitTestSuite) Test_IncidentWorkflow_ShouldGetOsSlackChannel() {
	s.env.ExecuteWorkflow(IncidentWorkflow, nil)

	s.True(s.env.IsWorkflowCompleted())

	err := s.env.GetWorkflowError()
	s.Error(err)
	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("required environment variable SLACK_CHANNEL is not set", applicationErr.Error())
}

func (s *UnitTestSuite) Test_IncidentWorkflow_PostMessageSuccess() {
	s.Suite.T().Setenv("SLACK_CHANNEL", "my slack channel")

	s.env.OnActivity(slackActivities.PostMessageActivity, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, slackDate models.SlackActivityData) (models.MessageDetails, error) {
			s.Equal("my slack channel", slackDate.ChannelId)
			return models.MessageDetails{}, nil
		})
	s.env.ExecuteWorkflow(IncidentWorkflow, nil)

	s.True(s.env.IsWorkflowCompleted())

	err := s.env.GetWorkflowError()
	s.Nil(err)
	var result string
	err = s.env.GetWorkflowResult(&result)
	s.Nil(err)
	s.Equal("DONE", result)
}
