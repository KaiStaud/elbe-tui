package utils

import (
	"elbe-prj/containers"
	"fmt"
	"os/exec"
)

func CreateDebian(c containers.Debianize) {
	cpDebian(c)
	substituteArch(c)
	substituteRelease(c)
	SubstituteDefconfig(c)
}

func cpDebian(c containers.Debianize) {
	cpCmd := exec.Command("cp", "-rf", c.TemplateDir, c.DestDir)
	err := cpCmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func substituteArch(c containers.Debianize) {
	var replace = "s/dummy_arch/" + c.Arch + "/g"
	//	sed -i 's/foo/bar/g' *
	cpCmd := exec.Command("sed", "-i", replace, "*")
	err := cpCmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}

}

func substituteRelease(c containers.Debianize) {
	var replace = "s/dummy_release/" + c.Release + "/g"
	//	sed -i 's/foo/bar/g' *
	cpCmd := exec.Command("sed", "-i", replace, "*")
	err := cpCmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}

}

func SubstituteDefconfig(c containers.Debianize) {
	var replace = "s/dummy_config/" + c.Config + "/g"
	//	sed -i 's/foo/bar/g' *
	cpCmd := exec.Command("sed", "-i", replace, "*")
	err := cpCmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}

}
