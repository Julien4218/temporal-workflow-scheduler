package instrumentation

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"

	"github.com/Julien4218/temporal-workflow-scheduler/utils"
)

var NrApp *newrelic.Application
var Hostname string

func Init() {
	if NrApp != nil {
		log.Fatal(errors.New("NewRelic instrumentation is already setup."))
	}
	app, err := instrumentation()
	if err != nil {
		log.Fatalf("Could not setup newrelic instrumentation, detail:%s", err.Error())
	}
	NrApp = app

	hostname, _ := os.Hostname()
	Hostname = hostname
}

func instrumentation() (*newrelic.Application, error) {
	var appName = os.Getenv("NEW_RELIC_APP_NAME")
	if appName == "" {
		appName = "workflow-er-poc"
	}

	var licenseKey = os.Getenv("NEW_RELIC_LICENSE_KEY")
	if licenseKey == "" {
		return nil, errors.New("License key is missing from environment variable with key NEW_RELIC_LICENSE_KEY")
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigEnabled(true),
		newrelic.ConfigFromEnvironment(),
		newrelic.ConfigAppLogEnabled(true),
	)
	if nil != err {
		return nil, err
	}
	nrlogrusFormatter := nrlogrus.NewFormatter(app, &logrus.TextFormatter{})
	logrus.SetFormatter(nrlogrusFormatter)
	level := logrus.DebugLevel
	if l := os.Getenv("NEW_RELIC_LOG_LEVEL"); l != "" {
		level = utils.GetLevelFromString(l, logrus.DebugLevel)
	}
	logrus.SetLevel(level)
	logrus.Debug("instrumentation is setup")

	// Wait for the application to connect.
	if err = app.WaitForConnection(15 * time.Second); nil != err {
		return nil, err
	}

	return app, nil
}
