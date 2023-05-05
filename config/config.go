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
	viper.SetConfigFile("config.env")
	viper.ReadInConfig()

	return Env{
		ElbeBin:      viper.Get("ELBE").(string),
		DefaultDlDir: viper.Get("DEFAULT_DOWNLOAD_DIR").(string),
		WorkDir:      viper.GetString("WORK_DIR"),
		SubmitFlags:  viper.GetString("work_dir"),
		HighlightPrj: viper.GetString("HIGHLIGHT_PBUILDS"),
	}

}
