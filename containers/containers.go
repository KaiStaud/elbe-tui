package containers

type BuildResult int
type DownloadState int

const (
	Build_Failed BuildResult = iota // 0
	Build_Done
	Empty_Project
	Busy
	Needs_Build
	PBuild_Failed
	PBuild_Done
	PBuildBusy
	PBuildNeeds_Build                      // 4
	DownloadStarted   DownloadState = iota // 5
	DownloadPathEntered
	DownloadFinished
)

type Project struct {
	Path     string
	Name     string
	Result   BuildResult
	Progress DownloadState
}
