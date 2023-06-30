/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"elbe-prj/containers"
	"elbe-prj/erlang"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// debianizeCmd represents the debianize command
var debianizeCmd = &cobra.Command{
	Use:   "debianize",
	Short: "debianize source folder",
	Long:  `Creates a customized debian folder for given source directory`,
	Run: func(cmd *cobra.Command, args []string) {
		/*
			fmt.Println("debianize called")
			debianize_sh := "/etc/elbe-tui/scripts/debianize.sh"
			template_dir := "/etc/elbe-tui/templates"
			output_dir := cmd.Flags().Lookup("output-dir").Value.String()
			package_type := cmd.Flags().Lookup("package-type").Value.String()
			stdout, err := exec.Command(debianize_sh, template_dir, output_dir, package_type).Output()
			if err != nil {
				log.Printf("Couldnt debianize source:%v", string(stdout))
			}
		*/
		var projects []containers.Project

		var m = erlang.InitialModel(projects)
		p := tea.NewProgram(m) //erlang.InitialModel(projects))
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(debianizeCmd)
	debianizeCmd.PersistentFlags().String("package-type", "", "bootloader, kernel, [dkms-]module, application")
	debianizeCmd.PersistentFlags().String("output-dir", "", "")
}
