## about
elbe-tui (or as preferred shortenened etui) abstracts away elbe's cumbersome cli interface and (at least) aims to provide
instead a sructured,keyed and user-friendly "do-it-all in one" application to build debian based rfs images.
When started, etui list all projects in a color coded way:

![startscreen](https://github.com/KaiStaud/elbe-tui/blob/main/doc/01-etui_startscreen.png)

As you might already noticed, it also provides keyboard bindings to access the most common tasks. 
For example r resets a failed project, t deletes it and g downloads files:

![keybindings](https://github.com/KaiStaud/elbe-tui/blob/main/doc/03_reset_prj.webm)

When started etui lists all stored initvm projects, and provides direct access to them via its keybindings.
The user is recommended to provide a personized config.json file to configure etui, which will provide the working directory whith its structure.

Thanks to 
https://github.com/charmbracelet/bubbletea for providing an excellent cli-library
https://github.com/spf13/viper to provide extensible configuration functionality  
https://github.com/spf13/cobra for providing scriptable cli options and
https://github.com/esiqveland/notify to generate desktop notifications

## features
- projects are color keyed on their build result
- delete, reset and file-download on projects
- .swu-image generation
- pre and post build scripts
- debianization of userspace,kernel and bootloader package source
- pbuilder configuration:
    - auto generating project
    - uploading/ building packages
## getting started
clone this repo and build the binary with go build .
modify config.json and copy it and the directories scripts and templates to /etc/elbe-tui

run with ./elbe-tui, provide cli options to access advanced features.
help is provided with --help and manpage debian/elbe-tui8.man

## default keymapping
- t: delete project
- r: reset project
- arrow keys / aswd: navigation
- q / esc : quit view
- enter / return : submit input
- p: create .swu image
- g: get files


## clioptions

## work in progress

- virtualization support
    - docker / podman
    - k8s
    - hawkbit
- integration with custom mirrors and ppas
- povide etui as snap package / github releases
- integration with libdbus
    - notification on finished or failed builds
    

## future work
- integration w/ qemu for automized image testing
- web interface
- drop into C10shell on failure
- filtered make.log and log.txt's

## known bugs