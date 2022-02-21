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

		buildPage()
	},
}

// init is the initialization function for Cobra which defines commands and
// flags.
func init() {
	rootCmd.Flags().StringVarP(&logPath, "logPath", "l",
		"./mainnet-commitments-client.log", "File path to save log file to.")
	rootCmd.Flags().IntVarP(&logLevel, "logLevel", "v", 0,
		"Verbosity level for log printing (2+ = Trace, 1 = Debug, 0 = Info).")
}

// initLog initializes logging thresholds and the log path.
func initLog(logPath string, logLevel int) {

	// Set log file output
	logFile, err := os.OpenFile(
		logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Could not open log file %q: %+v\n", logPath, err)
	} else {
		jww.SetLogOutput(logFile)
	}

	// Check the level of logs to display
	if logLevel > 1 {
		// Turn on trace logs
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelTrace)
		jww.INFO.Printf("Log level set to: %s", jww.LevelTrace)
	} else if logLevel == 1 {
		// Turn on debugging logs
		jww.SetLogThreshold(jww.LevelDebug)
		jww.SetStdoutThreshold(jww.LevelDebug)
		jww.INFO.Printf("Log level set to: %s", jww.LevelDebug)
	} else {
		// Turn on info logs
		jww.SetLogThreshold(jww.LevelInfo)
		jww.SetStdoutThreshold(jww.LevelInfo)
		jww.INFO.Printf("Log level set to: %s", jww.LevelInfo)
	}
}
