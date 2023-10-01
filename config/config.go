package config

/*
elbe-prj collects its configuration by
- the initial config (s.c. "initconfig") provided by config.yaml
	- initconfig contains:
		- path to elbe binary
		- fallback directories
		- messaging subsystem
			- integration with dbus
			- notification over mail-client
			- push notification to mobile devices app
			- daemon mode
		-
- individual project-configs (s.c. prjconfig) stored and created in a database
	- prjconfigs consist of:
		- initvm projects
		- associated xmls
		- storage paths
		- user build pbuild packages

- individual project-environments (s.c prjenvs) provided with boardname.yaml.
	- prjenvs consist of:
		- board-xml
		- pbuilding sources ( either in git or on disk)
		- package directories ( either on disk,git or server)
		- optional output directory
		- optional post and prebuild commands
- once fully set up, utils enables the user to convert a prjconfig into prjenvs and vice-versa.
	Additionally, if specified the backing database can also be exported and imported.
*/
import (
	"github.com/spf13/viper"
)

type Env struct {
	ElbeBin      string
	DefaultDlDir string
	WorkDir      string
	SubmitFlags  string
	HighlightPrj string
}

func ReadEnv() Env {
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/elbe-tui/")
	viper.SetConfigName("config")
	viper.ReadInConfig()

	return Env{
		ElbeBin:      viper.GetString("elbe"),
		DefaultDlDir: viper.GetString("default_dl_dir"),
		WorkDir:      viper.GetString("work_dir"),
		HighlightPrj: viper.GetString("highlight_pbuilds"),
	}

}
