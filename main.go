package main

import (
	"fmt"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"os"
)

var logPath string
var logLevel int

func main() {
	Execute()
}

// Execute adds all child commands to the root command and sets the flags
// appropriately. This is called by main.main(). It only needs to happen once to
// the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "mainnet-commitments-client",
	Short: "Main command for mainnet-commitments client",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		initLog(logPath, logLevel)

		err := buildPage()
		if err != nil {
			jww.FATAL.Panic(err)
		}
	},
}

// init is the initialization function for Cobra which defines commands and
// flags.
func init() {
	rootCmd.Flags().StringVarP(&logPath, "logPath", "l", "test.log",
		"File path to save log file to.")
	rootCmd.Flags().IntVarP(&logLevel, "logLevel", "v", 2,
		"Verbosity level for log printing (2+ = Trace, 1 = Debug, 0 = Info).")
}

// initLog initializes logging thresholds and the log path. If not path is
// provided, the log output is not set. Possible values for logLevel:
//  0  = info
//  1  = debug
//  2+ = trace
func initLog(logPath string, logLevel int) {
	// Set log level to highest verbosity while setting up log files
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelTrace)

	// Set log file output
	if logPath != "" {
		logFile, err := os.OpenFile(
			logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			jww.ERROR.Printf("Could not open log file %q: %+v\n", logPath, err)
		} else {
			jww.INFO.Printf("Setting log output to %q", logPath)
			jww.SetLogOutput(logFile)
		}
	} else {
		jww.INFO.Printf("No log output set: no log path provided")
	}

	// Select the level of logs to display
	var threshold jww.Threshold
	if logLevel > 1 {
		// Turn on trace logs
		threshold = jww.LevelTrace
	} else if logLevel == 1 {
		// Turn on debugging logs
		threshold = jww.LevelDebug
	} else {
		// Turn on info logs
		threshold = jww.LevelInfo
	}

	// Set logging thresholds
	jww.SetLogThreshold(threshold)
	jww.SetStdoutThreshold(threshold)
	jww.INFO.Printf("Log level set to: %s", threshold)
}
