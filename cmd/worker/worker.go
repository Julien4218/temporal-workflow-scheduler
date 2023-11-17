package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Julien4218/temporal-workflow-scheduler/instrumentation"
	"github.com/Julien4218/temporal-workflow-scheduler/workflows"

	slackActivities "github.com/Julien4218/temporal-slack-activity/activities"
	slackInstrumentation "github.com/Julien4218/temporal-slack-activity/instrumentation"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var (
	QueueName string
)

func init() {
	workerCmd.Flags().StringVar(&QueueName, "queue", workflows.DefaultQueueName, "Queue")
}

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run worker",
	Run: func(cmd *cobra.Command, args []string) {

		err := godotenv.Load()
		if err != nil {
			fmt.Printf("Error loading .env file")
			os.Exit(1)
		}

		instrumentation.Init()
		logrus.Infof("%s-Worker started on queue %s", instrumentation.Hostname, QueueName)

		client, err := client.Dial(client.Options{
			HostPort:  os.Getenv("TEMPORAL_HOSTPORT"),
			Namespace: "default",
		})
		if err == nil {
			defer client.Close()
			workerInstance := worker.New(client, QueueName, worker.Options{})

			workerInstance.RegisterWorkflow(workflows.ErWorkflow)
			workerInstance.RegisterWorkflow(workflows.IncidentWorkflow)
			workerInstance.RegisterWorkflow(workflows.RetrospectiveWorkflow)

			slackInstrumentation.AddLogger(func(message string) { logrus.Info(message) })
			workerInstance.RegisterActivity(slackActivities.PostMessageActivity)
			workerInstance.RegisterActivity(slackActivities.GetMessageReactions)
			workerInstance.RegisterActivity(slackActivities.AddReaction)

			err = workerInstance.Run(worker.InterruptCh())
		}

		if err != nil {
			logrus.Errorf("%s-Worker exited with error: %v", instrumentation.Hostname, err)
		}
		logrus.Infof("%s-Worker exited", instrumentation.Hostname)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	// workerCmd.Use = appName

	// Silence Cobra's internal handling of command usage help text.
	// Note, the help text is still displayed if any command arg or
	// flag validation fails.
	workerCmd.SilenceUsage = true

	// Silence Cobra's internal handling of error messaging
	// since we have a custom error handler in main.go
	workerCmd.SilenceErrors = true

	err := workerCmd.Execute()
	return err
}
