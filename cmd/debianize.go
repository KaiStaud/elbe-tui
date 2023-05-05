/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// debianizeCmd represents the debianize command
var debianizeCmd = &cobra.Command{
	Use:   "debianize",
	Short: "debianize source folder",
	Long:  `Creates a customized source folder for given source directory`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("debianize called")
	},
}

func init() {
	rootCmd.AddCommand(debianizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	debianizeCmd.PersistentFlags().String("bootloader", "", "debianizes u-boot and tf-a bootloader-source")
	debianizeCmd.PersistentFlags().String("kernel", "", "debianizes linux kernel source")
	debianizeCmd.PersistentFlags().String("module", "", "creates a debianized kernel module")
	debianizeCmd.PersistentFlags().String("dkms", "", "creates a debianized dkms-kernel module")
	debianizeCmd.PersistentFlags().String("application", "", "debianizes application code")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// debianizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
