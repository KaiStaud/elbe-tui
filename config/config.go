package config

import "github.com/spf13/viper"

type Env struct {
	ElbeBin      string
	DefaultDlDir string
	WorkDir      string
	SubmitFlags  string
	HighlightPrj string
}

func ReadEnv() Env {
	viper.AddConfigPath("/home/sta/projects/go/elbe-prj")
	viper.SetConfigType("json")
	viper.SetConfigFile("config.json")
	viper.ReadInConfig()

	return Env{
		ElbeBin:      viper.GetString("elbe"),
		DefaultDlDir: viper.GetString("default_dl_dir"),
		WorkDir:      viper.GetString("work_dir"),
		HighlightPrj: viper.GetString("highlight_pbuilds"),
	}

}
