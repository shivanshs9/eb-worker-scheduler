package cmd

import (
	"os"

	app "github.com/shivanshs9/eb-worker-scheduler/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start the SQS scheduler, that triggers jobs via SQS based on given crontab",
	Run:   runScheduler,
}

var debug bool
var options *app.AppOptions

func init() {
	options = new(app.AppOptions)
	rootCmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "verbose logging")
	rootCmd.Flags().StringVarP(&options.QueueUrl, "queueUrl", "q", "", "Provide the queue URL (required)")
	rootCmd.MarkFlagRequired("queueUrl")

	rootCmd.Flags().StringVarP(&options.YamlPath, "path", "p", "cron.yaml", "Provide the path to cron.yaml file.")
	rootCmd.Flags().IntVarP(&options.DeduplicationBufferSec, "dedup-sec", "", 5, "Provide the second range for which to keep a same deduplication ID")
}

func runScheduler(cmd *cobra.Command, args []string) {
	log := logrus.New()
	log.Info("Verbose logging enabled")
	if debug {
		log.SetLevel(logrus.DebugLevel)
	}
	options.RetryCount = 3
	app.StartApp(options, log)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
