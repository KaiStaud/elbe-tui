/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"elbe-prj/containers"
	"elbe-prj/utils"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// headlessCmd represents the headless command
var headlessCmd = &cobra.Command{
	Use:   "headless",
	Short: "Enables command line functionality",
	Long:  `Allows further customized commands, enables pipelined commands and early exit`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("headless called")
		reset := cmd.Flags().Lookup("reset").Changed
		delete := cmd.Flags().Lookup("delete").Changed
		all := cmd.Flags().Lookup("all").Changed
		exit := cmd.Flags().Lookup("exit").Changed

		log.Printf("Running headless. Passed args translate to reset=%v delete=%v all=%v", reset, delete, all)

		var prjs = utils.GetProjects()
		var failed_prjs = utils.FilterProjects(prjs, containers.Build_Failed)
		var busy_prjs = utils.FilterProjects(prjs, containers.Build_Failed)
		var util_calls = 0
		//var to_reset_projects = utils.FilterProjects()
		if delete {
			log.Printf("Deleting all failed projects...")
			for _, v := range failed_prjs {
				utils.DeleteProject(v.Path, false)
				util_calls++
			}
			log.Printf("Deleted %d projects from initvm", util_calls)
		}
		if reset {
			log.Printf("Reseting all busy projects...")
			for _, v := range busy_prjs {
				utils.ResetProject(v.Path)
			}
		}
		if exit {
			os.Exit(0)
		}

	},
}

func init() {
	rootCmd.AddCommand(headlessCmd)
	headlessCmd.PersistentFlags().StringP("reset", "r", "", "reset busy a failed project")
	headlessCmd.PersistentFlags().StringP("delete", "d", "", "delete a project")
	headlessCmd.PersistentFlags().StringP("all", "a", "", "apply previous flags to all listed projects")
	headlessCmd.PersistentFlags().BoolP("exit", "e", false, "exit after executing commands ")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// headlessCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// headlessCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
