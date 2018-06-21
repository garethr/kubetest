package cmd

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/garethr/kubetest/kubetest"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the the command to run when kubetest is run
var RootCmd = &cobra.Command{
	Use:   "kubetest <file> [file...]",
	Short: "Run tests against a Kubernetes YAML file",
	Long:  `Run tests against a Kubernetes YAML file`,
	Run: func(cmd *cobra.Command, args []string) {
		initLogging()
		success := true
		windowsStdinIssue := false
		testsDir := viper.GetString("testsDir")
		stat, err := os.Stdin.Stat()
		if err != nil {
			// Stat() will return an error on Windows in both Powershell and
			// console until go1.9 when nothing is passed on stdin.
			// See https://github.com/golang/go/issues/14853.
			if runtime.GOOS != "windows" {
				log.Error(err)
				os.Exit(1)
			} else {
				windowsStdinIssue = true
			}
		}
		// We detect whether we have anything on stdin to process
		if !windowsStdinIssue && ((stat.Mode() & os.ModeCharDevice) == 0) {
			var buffer bytes.Buffer
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				buffer.WriteString(scanner.Text() + "\n")
			}
			runSuccess := kubetest.Runs(buffer.Bytes(), testsDir, "stdin")
			if success {
				success = runSuccess
			}
		} else {
			if len(args) < 1 {
				log.Fatal("You must pass at least one file as an argument")
			}
			for _, fileName := range args {
				filePath, _ := filepath.Abs(fileName)
				fileContents, err := ioutil.ReadFile(filePath)
				if err != nil {
					log.Fatal("Could not open file ", fileName)
				}
				runSuccess := kubetest.Runs(fileContents, testsDir, fileName)
				if success {
					success = runSuccess
				}
			}
		}
		if !success {
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}

func initLogging() {
	log.SetOutput(os.Stdout)
	if !viper.GetBool("verbose") {
		log.SetLevel(log.WarnLevel)
	}
	if viper.GetBool("useJson") {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		formatter := &log.TextFormatter{
			DisableTimestamp: true,
		}
		log.SetFormatter(formatter)
	}
}

func init() {
	viper.SetEnvPrefix("KUBETEST")
	viper.AutomaticEnv()

	RootCmd.PersistentFlags().StringP("tests", "t", "tests", "Test directory")
	viper.BindPFlag("testsDir", RootCmd.PersistentFlags().Lookup("tests"))

	RootCmd.PersistentFlags().Bool("json", false, "Output results as JSON")
	viper.BindPFlag("useJson", RootCmd.PersistentFlags().Lookup("json"))

	RootCmd.PersistentFlags().Bool("verbose", false, "Output passes as well as failures")
	viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))

	RootCmd.PersistentFlags().Bool("ignoreemptyfiles", false, "Ignore empty files")
	viper.BindPFlag("ignoreemptyfiles", RootCmd.PersistentFlags().Lookup("ignoreemptyfiles"))

}
