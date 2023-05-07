package containers

type BuildResult int
type DownloadState int

const (
	Build_Failed BuildResult = iota // 0
	Build_Done
	Empty_Project
	Busy
	Needs_Build
	DownloadStarted DownloadState = iota // 5
	DownloadPathEntered
	DownloadFinished
)

type Project struct {
	Path     string
	Name     string
	Result   BuildResult
	Progress DownloadState
}

type Debianize struct {
	PrjName     string
	TemplateDir string
	DestDir     string
	Arch        string
	Config      string
	Release     string
}
