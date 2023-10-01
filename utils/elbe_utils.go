package utils

import (
	"elbe-prj/config"
	"elbe-prj/containers"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

var (
	elbe_bin   = ""
	pbuild_prj = "/var/cache/elbe/a79a01ed-9091-4f8f-9f20-1ed6a7060634+"
)

func LoadConfig(elbe string) {
	elbe_bin = elbe
}

func GetProjects() []containers.Project {
	var c = config.ReadEnv()
	elbe_bin = c.ElbeBin
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
	}

	var projects []containers.Project
	s := SplitLines(string(stdout))

	for i, v := range s {
		log.Println(i, v)
		p := ParseLine(v)
		projects = append(projects, p)
	}
	return projects
}

func DeleteProject(path string, needs_reset bool) {
	app := elbe_bin
	arg0 := "control"
	arg1 := "del_project"
	arg2 := path
	log.Printf(" %s %s %s %s", app, arg0, arg1, arg2)
	cmd := exec.Command(elbe_bin, "control", "del_project", path)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(stdout)
		log.Println(err.Error())
		fmt.Println(err.Error())
		return //colorizeErrorMessage(stdout)
	}

}

func GetFiles(path string, target_dir string) {
	app := elbe_bin
	arg0 := "control"
	arg1 := "get_files"
	arg2 := "--output"
	arg3 := target_dir
	arg4 := path

	log.Printf(" %s %s %s %s %s %s", app, arg0, arg1, arg2, arg3, arg4)
	cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(stdout)
		log.Println(err.Error())
		fmt.Println(err.Error())
		return // colorizeErrorMessage(stdout)
	}

}

func ResetProject(path string) {
	cmd := exec.Command(elbe_bin, "control", "reset_project", path)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(stdout)
		log.Println(err.Error())
		fmt.Println(err.Error())
		return // colorizeErrorMessage(stdout)
	}
}

func StartInitvm() {
	cmd := exec.Command(elbe_bin, "initvm", "start")
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(stdout)
		log.Println(err.Error())
		fmt.Println(err.Error())
		return // colorizeErrorMessage(stdout)
	}
}

func StopInitvm() {
	cmd := exec.Command(elbe_bin, "initvm", "stop")
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(stdout)
		log.Println(err.Error())
		fmt.Println(err.Error())
		return // colorizeErrorMessage(stdout)
	}
}

func CreatePackage(project string) {
	cmd := exec.Command(elbe_bin, "pbuilder", "build", project, "--out=../out")
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(stdout)
		log.Println(err.Error())
		fmt.Println(err.Error())
		return // colorizeErrorMessage(stdout)
	}
}

func UploadPackage(prj string, path string, ch chan<- bool) {
	var c = config.ReadEnv()
	elbe_bin = c.ElbeBin

	go func() {
		// cmd := exec.Command(elbe_bin, "pbuilder", "create", "--xmlfile="+path, "--writeproject=/tmp/test.prj")
		cmd := exec.Command("sleep", "3")
		stdout, err := cmd.Output()
		ch <- true
		if err != nil {
			log.Println(stdout)
			log.Println(err.Error())
			fmt.Println(err.Error())

		}
	}()

}

func BuildImage(project string, xml string) {
	cmd := exec.Command(elbe_bin, "control", "build", project)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(stdout)
		log.Println(err.Error())
		fmt.Println(err.Error())
		return // colorizeErrorMessage(stdout)
	}
}

func CreateProject(path string, prj string, ch chan<- bool) string {
	var c = config.ReadEnv()
	elbe_bin = c.ElbeBin

	go func() {
		// cmd := exec.Command(elbe_bin, "pbuilder", "create", "--xmlfile="+path, "--writeproject=/tmp/test.prj")
		cmd := exec.Command("sleep", "5")
		stdout, err := cmd.Output()
		ch <- true
		if err != nil {
			log.Println(stdout)
			log.Println(err.Error())
			fmt.Println(err.Error())

		}
	}()
	buf, e := ioutil.ReadFile("/tmp/test.prj")
	if e == nil {
		print("couldnt read in tmp-file %e", e)
	}
	return string(buf)
}

func ReadProjectString(path string) string {
	return "not implemented"
}

// Get all projects wich Buildresult matches the filter
func FilterProjects(p []containers.Project, filter containers.BuildResult) []containers.Project {
	var list []containers.Project

	for _, v := range p {
		if v.Result == filter {
			list = append(list, v)
		}
	}
	return list
}
