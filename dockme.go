package dockme

// this a first pass, expect refactors

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const VERSION = "0.3.1"
const DOCKER = "docker"
const ACTION = "run"
const DEFAULT_TEMPLATE = "default"
const DEFAULT_CONFIG = "Dockme.yml"

type Dockme struct {
	Hostname    string `yaml:"hostname,omitempty"`
	Domainname  string `yaml:"domainname,omitempty"`
	User        string `yaml:"user,omitempty"`
	Cpuset      string `yaml:"cpuset,omitempty"`
	Image       string `yaml:"image,omitempty"`
	Workdir     string `yaml:"workdir,omitempty"`
	MacAddress  string `yaml:"mac_address,omitempty"`
	Name        string `yaml:"name,omitempty"`
	Entrypoint  string `yaml:"entrypoint,omitempty"`
	Cmd         string `yaml:"command,omitempty"`
	Source      string `yaml:"source,omitempty"`
	Destination string `yaml:"destination,omitempty"`
	Config      string `yaml:"omit0,omitempty"`

	Memory     int `yaml:"memory,omitempty"`
	MemorySwap int `yaml:"memory_swap,omitempty"`
	CpuShares  int `yaml:"cpu_shares,omitempty"`

	Rm          bool `yaml:"rm,omitempty"`
	Interactive bool `yaml:"interactive,omitempty"`
	Tty         bool `yaml:"tty,omitempty"`
	Sudo        bool `yaml:"sudo,omitempty"`
	Dryrun      bool `yaml:"omit1,omitempty"`
	Save        bool `yaml:"omit2,omitempty"`
	SrcFromCwd  bool `yaml:"omit3,omitempty"` // don't save source

	Expose      []string `yaml:"expose,omitempty,flow"`
	Env         []string `yaml:"env,omitempty,flow"`
	Links       []string `yaml:"link,omitempty,flow"`
	Volume      []string `yaml:"volume,omitempty,flow"`
	VolumesFrom []string `yaml:"volumes_from,omitempty,flow"`
	Publish     []string `yaml:"publish,omitempty,flow"`
}

func (dm *Dockme) Args() []string {
	args := []string{ACTION}

	// TODO: find a more meta way?

	if dm.Name != "" {
		args = append(args, fmt.Sprintf("--name=%s", dm.Name))
	}

	if dm.Hostname != "" {
		args = append(args, fmt.Sprintf("--hostname=%s", dm.Hostname))
	}

	if dm.Domainname != "" {
		args = append(args, fmt.Sprintf("--domainname=%s", dm.Domainname))
	}

	if dm.User != "" {
		args = append(args, fmt.Sprintf("--user=%s", dm.User))
	}

	if dm.Cpuset != "" {
		args = append(args, fmt.Sprintf("--cpuset=%s", dm.Cpuset))
	}

	if dm.Workdir != "" {
		args = append(args, fmt.Sprintf("--workdir=%s", dm.Workdir))
	}

	if dm.MacAddress != "" {
		args = append(args, fmt.Sprintf("--mac-address=%s", dm.MacAddress))
	}

	if dm.Entrypoint != "" {
		args = append(args, fmt.Sprintf("--mac-address=%s", dm.Entrypoint))
	}

	if (dm.Source == "" || dm.Destination == "") && (dm.Source != "" || dm.Destination != "") {
		log.Fatal(errors.New("Can't specify source or destination, both or neither"))
	}

	if dm.Source != "" && dm.Destination != "" {
		dm.Volume = append(dm.Volume,
			fmt.Sprintf("%s:%s", dm.Source, dm.Destination))
	}

	if dm.Memory > 0 {
		args = append(args, fmt.Sprintf("--memory=%d", dm.Memory))
	}

	if dm.MemorySwap > 0 {
		args = append(args, fmt.Sprintf("--memory-swap=%d", dm.MemorySwap))
	}

	if dm.CpuShares > 0 {
		args = append(args, fmt.Sprintf("--cpu-shares=%d", dm.CpuShares))
	}

	if dm.Rm {
		args = append(args, "--rm")
	}

	if dm.Tty {
		args = append(args, "--tty")
	}

	if dm.Interactive {
		args = append(args, "--interactive")
	}

	if len(dm.Expose) > 0 {
		for _, e := range dm.Expose {
			args = append(args, fmt.Sprintf("--expose=%s", e))
		}
	}

	if len(dm.Env) > 0 {
		for _, e := range dm.Env {
			args = append(args, fmt.Sprintf("--env=%s", e))
		}
	}

	if len(dm.Links) > 0 {
		for _, e := range dm.Links {
			args = append(args, fmt.Sprintf("--link=%s", e))
		}
	}

	if len(dm.Volume) > 0 {
		for _, e := range dm.Volume {
			args = append(args, fmt.Sprintf("--volume=%s", e))
		}
	}

	if len(dm.VolumesFrom) > 0 {
		for _, e := range dm.VolumesFrom {
			args = append(args, fmt.Sprintf("--volumes-from=%s", e))
		}
	}

	if len(dm.Publish) > 0 {
		for _, e := range dm.Publish {
			args = append(args, fmt.Sprintf("--publish=%s", e))
		}
	}

	if dm.Image == "" {
		log.Fatal(errors.New("Image is required"))
	}

	args = append(args, dm.Image)
	args = append(args, strings.Split(dm.Cmd, " ")...)

	return args
}

func (dm *Dockme) Exec() {
	var exe string
	var args []string
	if dm.Sudo {
		exe = "sudo"
		args = []string{DOCKER}
		args = append(args, dm.Args()...)
	} else {
		exe = DOCKER
		args = dm.Args()
	}

	fmt.Printf("+ %s %s\n", exe, strings.Join(args, " "))
	if dm.Dryrun {
		fmt.Println("")
		os.Exit(0)
	}

	cmd := exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func (dm *Dockme) ApplyDefaults() {
	// source is local dir if not
	if dm.Source == "" {
		dm.SrcFromCwd = true
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
		dm.Source = cwd
	}

	// hostname is image basename
	if dm.Hostname == "" {
		imageSansTag := strings.Split(dm.Image, ":")[0]
		imageSplit := strings.Split(imageSansTag, "/")
		dm.Hostname = imageSplit[len(imageSplit)-1]
	}

	if dm.Destination == "" {
		dm.Destination = "/src"
	}

	if dm.Workdir == "" {
		dm.Workdir = dm.Destination
	}

	if dm.Cmd == "" {
		dm.Cmd = "bash"
	}

}

func (dm *Dockme) SaveConfig() {
	// remove those that shouldn't be saved
	config := dm.Config
	dm.Config = ""

	srcFromCwd := dm.SrcFromCwd
	dm.SrcFromCwd = false

	dryrun := dm.Dryrun
	dm.Dryrun = false

	save := dm.Save
	dm.Save = false

	// don't save Source if from Cwd
	var savedSrc string
	if srcFromCwd {
		savedSrc = dm.Source
		dm.Source = ""
	}

	data, err := yaml.Marshal(&dm)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// add leading meta
	meta := append([]byte("# Generated by dockme v"), []byte(VERSION)...)
	meta = append(meta, []byte(" (github.com/jmervine/dockme)\n#\ton ")...)
	meta = append(meta, []byte(time.Now().Format("2006-01-02 15:04:05 -0700"))...)
	meta = append(meta, []byte("\n---\n")...)
	data = append(meta, data...)

	if err = ioutil.WriteFile(config, data, 0644); err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Wrote %s\n", config)

	if savedSrc != "" {
		dm.Source = savedSrc
	}

	// reset those that shouldn't be saved
	dm.Config = config
	dm.SrcFromCwd = srcFromCwd
	dm.Dryrun = dryrun
	dm.Save = save
}

var Templates = map[string]*Dockme{
	"default": &Dockme{
		Interactive: true,
		Tty:         true,
		Rm:          true,
		Cmd:         "bash",
	},
	"ruby": &Dockme{
		Image:       "jmervine/herokudev-ruby:latest",
		Name:        "rubydev",
		Interactive: true,
		Tty:         true,
		Rm:          true,
		Cmd:         "bash",
	},
	"rails": &Dockme{
		Image:       "jmervine/herokudev-rails:latest",
		Name:        "railsdev",
		Interactive: true,
		Tty:         true,
		Rm:          true,
		Cmd:         "bash",
	},
	"node": &Dockme{
		Image:       "jmervine/herokudev-node:latest",
		Name:        "nodedev",
		Interactive: true,
		Tty:         true,
		Rm:          true,
		Cmd:         "bash",
	},
	"nodebox": &Dockme{
		Image:       "jmervine/nodebox:latest",
		Name:        "nodeboxdev",
		Interactive: true,
		Tty:         true,
		Rm:          true,
		Cmd:         "bash",
	},
	"python2": &Dockme{
		Image:       "python:2-slim",
		Name:        "python",
		Interactive: true,
		Tty:         true,
		Rm:          true,
		Cmd:         "bash",
	},
	"python3": &Dockme{
		Image:       "python:3-slim",
		Name:        "python",
		Interactive: true,
		Tty:         true,
		Rm:          true,
		Cmd:         "bash",
	},
}
