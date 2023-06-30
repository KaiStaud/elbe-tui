% elbe-tui(1) Version 1.0 |  "elbe-tui" Documentation

NAME
====

**elbe-tui** — interacts with elbe initvm

SYNOPSIS
========

| **elbe-tui** 
| **elbe-tui** \[**COMMAND**]

DESCRIPTION
===========
Provides a cli-interface to elbe initvm.
Permanent options are provided in the appliations config file /etc/elbe-tui/config.json

Options
-------

-h, --help

:   Prints brief usage information.

-o, --output

:   Prints log output to stdout

-v, --version

:   Prints the current version number.

debianize

: Debianizes source-folder

genupdate

: builds a .swu-image

headless 

: Shortcuts for common tasks

Available options are:

1. -d, --delete delete projects with build state done,failed,Needs_Build
2. -r, --reset  reset busy project
3. -a --all  the previous commands are applied command to all applicable projects
4. -e --exit to skip the visualisation part of the application.

headless/cli-Examples
---------------------
| **elbe-tui headless -d all -r all** cleans initvm of any project, then runs the textinterface
| **elbe-tui headless -d all -e** cleans initvm of failed/done projects, then **skips** textinterface
| **elbe-tui debianize** creates an debian folder in current directory. The user needs to input additional data.
| **elbe-tui genupdate -n update-image-v1.swu** builds an .swu-update

keymap for text-interface
-------------------------
| **arrow-keys,a/s/w/d/f**  navigates up down
| **esc, ctrl+c** quit window application
| **enter/return** confirm input
| **p** spawns debianization dialogue. Navigate through input-fields with tab/previous-tab.
| **r** reset selected project
| **t** delete selected project
| **g** download projects' files from initvm 

FILES
=====
*/etc/elbe-tui/scripts*

: shell helper scripts

*/etc/elbe-tui/templates/*

:  Template files for .deb and .swu-packaging

*/etc/elbe-tui/config.json*

:   Global default dedication file.

BUGS
====

No bugs so far :)