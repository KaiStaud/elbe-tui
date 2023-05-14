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
