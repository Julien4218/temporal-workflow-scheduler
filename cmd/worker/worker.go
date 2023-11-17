package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Julien4218/temporal-workflow-scheduler/instrumentation"
	"github.com/Julien4218/temporal-workflow-scheduler/workflows"

	newrelicActivities "github.com/Julien4218/temporal-newrelic-activity/activities"
	newrelicInstrumentation "github.com/Julien4218/temporal-newrelic-activity/instrumentation"

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

		c, err := client.Dial(client.Options{
			HostPort:  os.Getenv("TEMPORAL_HOSTPORT"),
			Namespace: "default",
		})
		if err == nil {
			defer c.Close()
			workerInstance := worker.New(c, QueueName, worker.Options{})

			workerInstance.RegisterWorkflow(workflows.EventWorkflow)

			newrelicInstrumentation.AddLogger(func(message string) { logrus.Info(message) })
			workerInstance.RegisterActivity(newrelicActivities.NewCreateEventActivity)

			workflowOptions := client.StartWorkflowOptions{
				CronSchedule: "* * * * *",
				TaskQueue:    QueueName,
				// ...
			}

			input := &workflows.EventWorkflowInput{}
			workflowRun, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.EventWorkflow, input)
			if err != nil {
				logrus.Infof(fmt.Sprintf("Unable to execute workflow detail:%s", err))
				os.Exit(2)
				// ...
			}
			logrus.Infof(fmt.Sprintf("Started workflow WorkflowID:%s RunID:%s", workflowRun.GetID(), workflowRun.GetRunID()))

			// var cronScheduleStr string
			// cronScheduleStr = "* * * * *"
			// workflowOptions := client.StartWorkflowOptions{
			// 	ID:           workflowID,
			// 	TaskQueue:    QueueName,
			// 	CronSchedule: cronScheduleStr,
			// }

			// we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cron.SampleCronWorkflow,sid)
			// if err != nil {
			// 	logrus.Infof("Unable to execute workflow", err)
			// }
			// logrus.Infof("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

			_ = workerInstance.Run(worker.InterruptCh())
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
