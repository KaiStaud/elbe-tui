/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"elbe-prj/cmd"
	"elbe-prj/containers"
	"elbe-prj/erlang"
	"elbe-prj/utils"
	"log" // TODO: use zap instead
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

func main() {
	f, err := os.OpenFile("elbe.go.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	cmd.Execute()

	log.Println("Reading in config config.env")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/elbe-tui/")
	viper.SetConfigName("config")
	err = viper.ReadInConfig()
	var elbe_bin = viper.GetString("elbe")
	var elbe_dl_dir = viper.GetString("default_dl_dir")
	log.Println("elbe bin is located at" + elbe_bin + ", default-dl-dir is " + elbe_dl_dir)
	utils.LoadConfig(elbe_bin)

	app := elbe_bin

	arg0 := "control"
	arg1 := "list_projects"
	arg2 := ""
	arg3 := ""

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		log.Printf("Couldn't get initial project list from elbe-cmd %s %s %s, maybe your config.env isnt't handled correctly", app, arg0, arg1)
		log.Printf("get_projects backtrace:%s", err.Error())
		return
	}

	var projects []containers.Project
	s := utils.SplitLines(string(stdout))

	for i, v := range s {
		log.Println(i, v)
		p := utils.ParseLine(v)
		projects = append(projects, p)
	}
	var m = erlang.InitialModel(projects)
	p := tea.NewProgram(m) //erlang.InitialModel(projects))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
