// log provides log api.
// thanks to the helpful log tool logrus(https://github.com/Sirupsen/logrus)
package server

import (
	"github.com/Sirupsen/logrus"
)

var Log *logrus.Entry

func initLog(name string, level string) error {
	if Log == nil {
		// Log as JSON instead of the default ASCII formatter.
		logrus.SetFormatter(&logrus.JSONFormatter{})

		// Output to stderr instead of stdout, could also be a file.
		// logrus.SetOutput(os.Stderr)

		// logging level
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			return err
		}

		logrus.SetLevel(lvl)

		// default fields
		Log = logrus.WithFields(logrus.Fields{
			"service": name,
			"ip":      InternalIP,
		})
	}

	return nil
}
