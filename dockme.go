package main

// this a first pass, expect refactors

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const VERSION = "0.2.1"
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
	config      string

	Memory     int `yaml:"memory,omitempty"`
	MemorySwap int `yaml:"memory_swap,omitempty"`
	CpuShares  int `yaml:"cpu_shares,omitempty"`

	Rm          bool `yaml:"rm,omitempty"`
	Interactive bool `yaml:"interactive,omitempty"`
	Tty         bool `yaml:"tty,omitempty"`
	Sudo        bool `yaml:"sudo,omitempty"`
	dryrun      bool
	save        bool
	srcFromCwd  bool // don't save source

	Expose      []string `yaml:"expose,omitempty,flow"`
	Env         []string `yaml:"env,omitempty,flow"`
	Volume      []string `yaml:"volume,omitempty,flow"`
	VolumesFrom []string `yaml:"volumes_from,omitempty,flow"`
	Publish     []string `yaml:"publish,omitempty,flow"`
}

func (dm *Dockme) args() []string {
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

func (dm *Dockme) exec() {
	var exe string
	var args []string
	if dm.Sudo {
		exe = "sudo"
		args = []string{DOCKER}
		args = append(args, dm.args()...)
	} else {
		exe = DOCKER
		args = dm.args()
	}

	fmt.Printf("+ %s %s\n", exe, strings.Join(args, " "))
	if dm.dryrun {
		fmt.Println("")
		os.Exit(0)
	}

	cmd := exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func (dm *Dockme) applyDefaults() {
	// source is local dir if not
	if dm.Source == "" {
		dm.srcFromCwd = true
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

	if dm.Cmd == "" {
		dm.Cmd = "bash"
	}

}

func (dm *Dockme) saveConfig() {

	// don't save Source if from Cwd
	var savedSrc string
	if dm.srcFromCwd {
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

	if err = ioutil.WriteFile(dm.config, data, 0644); err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Wrote %s\n", dm.config)

	if savedSrc != "" {
		dm.Source = savedSrc
	}
}

var templates = map[string]*Dockme{
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

func main() {
	app := cli.NewApp()
	app.Name = "Dockme"
	app.Version = VERSION
	app.Author = "Joshua Mervine"
	app.Email = "joshua@mervine.net"
	app.Usage = "Simple wrapper for quickly spooling up docker containers for development."

	cli.AppHelpTemplate = `NAME:
    {{.Name}} - {{.Usage}}

USAGE:
    {{.Name}} [arguments...] [command]

VERSION:
    {{.Version}}

AUTHOR:
    {{.Author}}

OPTIONS:
    Only custom options or options whose usage strays from dockers
    usage have help messages. All other options map directly to docker
    run options, see Docker help and documentation for details.

    {{range .Flags}}{{.}}
    {{end}}
TEMPLATES:
    nodebox    nodebox template w/ 'jmervine/nodebox:latest'
    ruby       ruby template w/ 'jmervine/herokudev-ruby:latest'
    rails      rails template w/ 'jmervine/herokudev-rails:latest'
    node       node template w/ 'jmervine/herokudev-node:latest'
    python2    python template w/ 'python:2-slim'
    python3    python template w/ 'python:3-slim'
    help       Shows a list of commands or help for one command

`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "template, T",
			Usage: "set docker image template, see TEMPLATES below",
		},
		cli.StringFlag{
			Name:  "image, i",
			Usage: "set docker image [required]",
		},
		cli.StringFlag{
			Name:  "source, s",
			Usage: "local source directory",
		},
		cli.StringFlag{
			Name:  "destination, d",
			Usage: "[/src] container source directory",
		},
		cli.StringFlag{
			Name:  "publish, p",
			Usage: "list of ports to publish",
		},
		cli.StringFlag{
			Name:  "workdir, w",
			Usage: "set container workdir",
		},
		cli.BoolFlag{
			Name:  "dryrun, D",
			Usage: "show docker command to be run",
		},
		cli.BoolFlag{
			Name:  "save, S",
			Usage: "save configuration to file",
		},
		cli.StringFlag{
			Name:  "config, C",
			Usage: "conifguration file path",
			Value: DEFAULT_CONFIG,
		},
		cli.StringFlag{
			Name:  "expose, E",
			Usage: "list of ports to expose",
		},
		cli.StringFlag{
			Name:  "env, e",
			Usage: "list of environments",
		},
		cli.StringFlag{
			Name:  "volume, V",
			Usage: "list of volume mounts",
		},
		cli.StringFlag{
			Name:  "volumes-from",
			Usage: "list of containers to mount volumes from",
		},
		cli.StringFlag{
			Name: "name, n",
		},
		cli.BoolFlag{
			Name:  "sudo",
			Usage: "run Docker with sudo",
		},
		cli.BoolFlag{
			Name: "rm, r",
		},
		cli.BoolFlag{
			Name: "no-rm, k",
		},
		cli.BoolFlag{
			Name: "interactive, I",
		},
		cli.BoolFlag{
			Name: "no-interactive, x",
		},
		cli.BoolFlag{
			Name: "tty, t",
		},
		cli.BoolFlag{
			Name: "no-tty, N",
		},
		cli.StringFlag{
			Name: "entrypoint",
		},
		cli.StringFlag{
			Name: "user",
		},
		cli.StringFlag{
			Name: "hostname",
		},
		cli.StringFlag{
			Name: "domainname",
		},
		cli.StringFlag{
			Name: "mac-address",
		},
		cli.StringFlag{
			Name: "cpuset",
		},
		cli.StringFlag{
			Name: "memory",
		},
		cli.StringFlag{
			Name: "memory-swap",
		},
	}

	app.Action = func(c *cli.Context) {
		template := DEFAULT_TEMPLATE
		if c.String("template") != "" {
			template = c.String("template")
		}
		var dockme *Dockme

		config := c.String("config")

		read, err := ioutil.ReadFile(config)
		if err == nil {
			if err = yaml.Unmarshal(read, &dockme); err != nil {
				log.Fatalf("error: %v", err)
			}
		} else {
			dockme = templates[template]
		}

		dockme.config = config

		if c.String("image") != "" {
			dockme.Image = c.String("image")
		}

		// ensure image at least
		if dockme.Image == "" {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		if c.String("source") != "" {
			dockme.Source = c.String("source")
		}
		if c.String("destination") != "" {
			dockme.Destination = c.String("destination")
		}
		if c.String("workdir") != "" {
			dockme.Workdir = c.String("workdir")
		}
		if c.String("name") != "" {
			dockme.Name = c.String("name")
		}
		if c.String("entrypoint") != "" {
			dockme.Entrypoint = c.String("entrypoint")
		}
		if c.String("user") != "" {
			dockme.User = c.String("user")
		}
		if c.String("hostname") != "" {
			dockme.Hostname = c.String("hostname")
		}
		if c.String("domainname") != "" {
			dockme.Domainname = c.String("domainname")
		}
		if c.String("mac-address") != "" {
			dockme.MacAddress = c.String("mac-address")
		}
		if c.String("cpuset") != "" {
			dockme.Cpuset = c.String("cpuset")
		}
		if c.Int("memory") > 0 {
			dockme.Memory = c.Int("memory")
		}
		if c.Int("memory-swap") > 0 {
			dockme.MemorySwap = c.Int("memory-swap")
		}
		if c.Bool("sudo") {
			dockme.Sudo = c.Bool("sudo")
		}
		if c.Bool("rm") {
			dockme.Rm = c.Bool("rm")
		}
		if c.Bool("no-rm") {
			dockme.Rm = false
		}
		if c.Bool("interactive") {
			dockme.Interactive = c.Bool("interactive")
		}
		if c.Bool("no-interactive") {
			dockme.Interactive = false
		}
		if c.Bool("tty") {
			dockme.Tty = c.Bool("tty")
		}
		if c.Bool("no-tty") {
			dockme.Tty = false
		}
		if c.Bool("dryrun") {
			dockme.dryrun = c.Bool("dryrun")
		}
		if c.Bool("save") {
			dockme.save = true
		}
		if c.String("expose") != "" {
			dockme.Expose = split(c.String("expose"))
		}
		if c.String("env") != "" {
			dockme.Env = split(c.String("env"))
		}
		if c.String("volume") != "" {
			dockme.Volume = split(c.String("volume"))
		}
		if c.String("volumes-from") != "" {
			dockme.VolumesFrom = split(c.String("volumes-from"))
		}
		if c.String("publish") != "" {
			dockme.Publish = split(c.String("publish"))
		}

		cmd := strings.Join(c.Args(), " ")
		if cmd != "" {
			dockme.Cmd = cmd
		}

		dockme.applyDefaults()

		if dockme.save {
			dockme.saveConfig()
		}

		dockme.exec()
	}
	app.Run(os.Args)
}

func split(s string) []string {
	sp := strings.Split(s, ",")
	if len(sp) == 1 {
		sp = strings.Split(s, " ")
	}

	for n, i := range sp {
		sp[n] = strings.TrimSpace(i)
	}

	return sp
}
