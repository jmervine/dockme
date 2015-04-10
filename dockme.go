package main

// this a first pass, expect refactors

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	//"path/filepath"
	"strings"
)

const DOCKER = "docker"
const ACTION = "run"

type dockme struct {
	hostname    string
	domainname  string
	user        string
	cpuset      string
	image       string
	workdir     string
	macAddress  string
	name        string
	expose      string // straying from type
	entrypoint  string // straying from type
	cmd         string // straying from type
	source      string // custom arg
	destination string // custom arg
	template    string // custom arg

	memory     int64
	memorySwap int64
	cpuShares  int64

	interactive bool
	tty         bool
	dryrun      bool // custom arg
	save        bool // custom arg

	env         []string
	volume      []string // straying from type
	volumesFrom []string // straying from type
}

func (dm *dockme) args() []string {
	args := []string{ACTION}

	// TODO: find a more meta way?

	if dm.name != "" {
		args = append(args, fmt.Sprintf("--name=%s", dm.name))
	}

	if dm.hostname != "" {
		args = append(args, fmt.Sprintf("--hostname=%s", dm.hostname))
	}

	if dm.domainname != "" {
		args = append(args, fmt.Sprintf("--domainname=%s", dm.domainname))
	}

	if dm.user != "" {
		args = append(args, fmt.Sprintf("--user=%s", dm.user))
	}

	if dm.cpuset != "" {
		args = append(args, fmt.Sprintf("--cpuset=%s", dm.cpuset))
	}

	if dm.expose != "" {
		args = append(args, fmt.Sprintf("--expose=%s", dm.expose))
	}

	if dm.workdir != "" {
		args = append(args, fmt.Sprintf("--workdir=%s", dm.workdir))
	}

	if dm.macAddress != "" {
		args = append(args, fmt.Sprintf("--mac-address=%s", dm.macAddress))
	}

	if dm.entrypoint != "" {
		args = append(args, fmt.Sprintf("--mac-address=%s", dm.entrypoint))
	}

	if (dm.source == "" || dm.destination == "") && (dm.source != "" || dm.destination != "") {
		log.Fatal(errors.New("Can't specify source or destination, both or neither"))
	}

	if dm.source != "" && dm.destination != "" {
		dm.volume = append(dm.volume,
			fmt.Sprintf("%s:%s", dm.source, dm.destination))
	}

	if dm.memory > 0 {
		args = append(args, fmt.Sprintf("--memory=%d", dm.memory))
	}

	if dm.memorySwap > 0 {
		args = append(args, fmt.Sprintf("--memory-swap=%d", dm.memorySwap))
	}

	if dm.cpuShares > 0 {
		args = append(args, fmt.Sprintf("--cpu-shares=%d", dm.cpuShares))
	}

	if dm.tty {
		args = append(args, "--tty")
	}

	if dm.interactive {
		args = append(args, "--interactive")
	}

	if len(dm.env) > 0 {
		for _, e := range dm.env {
			args = append(args, fmt.Sprintf("--env=%s", e))
		}
	}

	if len(dm.volume) > 0 {
		for _, e := range dm.volume {
			args = append(args, fmt.Sprintf("--volume=%s", e))
		}
	}

	if len(dm.volumesFrom) > 0 {
		for _, e := range dm.volumesFrom {
			args = append(args, fmt.Sprintf("--volumes-from=%s", e))
		}
	}

	if dm.image == "" {
		log.Fatal(errors.New("Image is required"))
	}
	args = append(args, dm.image)

	if dm.cmd == "" {
		dm.cmd = "bash"
	}
	args = append(args, dm.cmd)

	return args
}

func (dm *dockme) exec() {
	args := dm.args()

	log.Printf("%s %s\n", DOCKER, strings.Join(args, " "))
	if dm.dryrun {
		log.Println("")
		os.Exit(0)
	}

	cmd := exec.Command(DOCKER, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

var templates = map[string]*dockme{
	"default": &dockme{
		image:       "jmervine/zshrc:latest",
		name:        "dockme",
		interactive: true,
		tty:         true,
		cmd:         "bash",
	},
	"ruby": &dockme{
		image:       "ruby:latest",
		name:        "rubydev",
		interactive: true,
		tty:         true,
		cmd:         "bash",
	},
	"node": &dockme{
		image:       "node:latest",
		name:        "nodedev",
		interactive: true,
		tty:         true,
		cmd:         "bash",
	},
	"nodebox": &dockme{
		image:       "jmervine/nodebox:latest",
		name:        "nodeboxdev",
		interactive: true,
		tty:         true,
		cmd:         "bash",
	},
}

func main() {
	template := "default"

	// always start with defaults
	dm := templates[template]

	if dm.source == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
		dm.source = cwd
	}

	if dm.destination == "" {
		dm.destination = "/src"
	}

	//dm.dryrun = true

	dm.exec()
}
