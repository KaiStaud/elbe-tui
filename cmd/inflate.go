/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"elbe-prj/containers"
	"elbe-prj/env"
	"elbe-prj/erlang"

	"github.com/spf13/cobra"
)

// ./elbe-prj inflate --env-file /home/sta/projects/elbe-tui/test2.json
var inflateCmd = &cobra.Command{
	Use:   "inflate",
	Short: "inflate environment from json",
	Long:  `Creates and populates projects by reading specified json`,
	Run: func(cmd *cobra.Command, args []string) {

		if cmd.Flags().Lookup("env-file").Changed {
			//

			wait_busy := make(chan containers.PBuilderState)
			info := make(chan string)

			envs := env.ReadEnvProject(cmd.Flags().Lookup("env-file").Value.String())

			// todo: wrap this !
			go func() {
				env.InflateEnv("", envs, wait_busy, info)
			}()
			go func() {
				erlang.InitWorker(wait_busy, info)
			}()
		}
	},
}

func init() {
	rootCmd.AddCommand(inflateCmd)
	inflateCmd.PersistentFlags().String("env-file", "", "")
}
