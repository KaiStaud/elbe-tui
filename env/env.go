package env

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
	"elbe-prj/containers"
	"elbe-prj/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Pack struct {
	//	src  BuildResult
	path string
}

type DefaultEnvironments struct {
	DefaultEnvironments []DefaultEnvironment `json:"default_environments"`
}

// User struct which contains a name
// a type and a list of social links
type DefaultEnvironment struct {
	Name          string        `json:"name"`
	Xml           string        `json:"xml"`
	Prj           string        `json:"prj"`
	PrjFile       string        `json:"prjfile"`
	Volume        int           `json:"volume"`
	PbuildSources PbuildSources `json:"pbuild_sources"`
}

type PbuildSources struct {
	PbuildSources []pbuild_source `json:"pbuild_src_tree"`
}

type pbuild_source struct {
	PackageName string `json:"package_name"`
	Path        string `json:"path"`
	UploadPath  string `json:"upload_from"`
}

func ReadEnvProject(path string) DefaultEnvironments {

	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var users DefaultEnvironments

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &users)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for i := 0; i < len(users.DefaultEnvironments); i++ {
		//		fmt.Println("User Type: " + users.DefaultEnvironments[i].)
		fmt.Println("User Xml: " + users.DefaultEnvironments[i].Xml)
		fmt.Println("User Name: " + users.DefaultEnvironments[i].Name)
		fmt.Println("Facebook Url: " + users.DefaultEnvironments[i].PbuildSources.PbuildSources[1].PackageName)
	}
	return users
}

func InflateEnv(path string, envs DefaultEnvironments, build_finished chan<- containers.PBuilderState, info chan<- string) {
	prj := ""
	done := false
	ch := make(chan bool)

	for i := 0; i < len(envs.DefaultEnvironments); i++ {
		build_finished <- containers.CreatePrj
		done = false
		// Get project from prj
		if envs.DefaultEnvironments[i].Prj != "" {
			prj = envs.DefaultEnvironments[i].Prj
		} else if envs.DefaultEnvironments[i].PrjFile != "" {
			prj = utils.ReadProjectString(envs.DefaultEnvironments[i].PrjFile)
		} else {
			// pbuilder create xml to create project -> save the prjvalue for later use
			build_finished <- containers.CreatePrj
			prj = utils.CreateProject(envs.DefaultEnvironments[i].Xml, envs.DefaultEnvironments[i].PrjFile, ch)
			for done == false {
				select {
				case _ = <-ch:
					done = true
				default:
					info <- fmt.Sprintf("Setting up project %s", envs.DefaultEnvironments[i].Xml)
				}

			}
			// TODO: also write the missing fields?
		}
		build_finished <- containers.UploadPkg

		// upload all packages

		for j := 0; j < len(envs.DefaultEnvironments[i].PbuildSources.PbuildSources); j++ {
			utils.UploadPackage(prj, envs.DefaultEnvironments[i].PbuildSources.PbuildSources[j].Path+"/"+envs.DefaultEnvironments[i].PbuildSources.PbuildSources[j].PackageName, ch)

			done = false
			for done == false {
				select {
				case _ = <-ch:
					done = true
				default:
					info <- envs.DefaultEnvironments[i].PbuildSources.PbuildSources[j].PackageName
				}
			}

		}
	}
	build_finished <- containers.PbuilderDone

}

func DeflateEnv() {

}
