/*
Copyright © 2024 Chad
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/chadsmith12/filefinder.git/packages/filescanner"
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

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().IntP("workers", "w", 3, "The number of workers to use to find files. Defauls to 3")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Invalid usage. Expects starting directory and pattern")
		os.Exit(1)
	}

	pattern, err := regexp.Compile(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	numberWorkers, err := cmd.Flags().GetInt("workers")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get number of workers, using %d by default...\n", 3)
		numberWorkers = 3
	}
	worker := filescanner.NewFileWorker(numberWorkers)
	worker.StartWorkers(args[0], pattern)
	for result := range worker.Read() {
		fmt.Printf("%s #%d: %s\n", result.File, result.LineNumber, result.Text)
	}
}
