/*
Copyright © 2024 Chad
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/chadsmith12/filefinder.git/packages/worker"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "filefinder",
	Short: "A simple utility program like grep, that finds a pattern in files.",
	Long: `A simple utility program, like grep, that finds a pattern in files. 
	This will walk through the directory given and go through all files to find a pattern given.
	Will walk the files in a concurrent fashion and you can manage how many worker threads it should use.`,
	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Invalid usage. Expects starting directory and pattern")
		os.Exit(1)
	}

	worker := worker.NewFileWorker(args[0], args[1])
	worker.Start()

	for result := range worker.Result() {
		fmt.Printf("#%d: %s - %s\n", result.LineNumber, result.File, result.Text)
	}
}
