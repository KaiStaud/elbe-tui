/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Execute additional build steps before/after elbe submit",
	Long:  `Execute pre/postbuild script `,
	Run: func(cmd *cobra.Command, args []string) {

		if cmd.Flags().Lookup("pre-script").Changed {
			fmt.Println("w /pre called")
			stdout, err := exec.Command("/bin/sh", cmd.Flags().Lookup("pre-script").Value.String()).Output()
			if err != nil {
				log.Printf("Couldnt run pre-build script, recovered error:%v", string(stdout))
			}
		}
		if cmd.Flags().Lookup("post-script").Changed {
			stdout, err := exec.Command("/bin/sh", cmd.Flags().Lookup("post-script").Value.String()).Output()
			if err != nil {
				log.Printf("Couldnt run post-build script, recovered error:%v", string(stdout))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pipelineCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	pipelineCmd.PersistentFlags().String("pre-script", "", "prebuild pipeline-script")
	pipelineCmd.PersistentFlags().String("post-script", "", "postbuild pipeline-script")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pipelineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
