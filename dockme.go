package main

// this a first pass, expect refactors

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/exec"

	//"path/filepath"
	"strings"
)

const DOCKER = "docker"
const ACTION = "run"
const DEFAULT_TEMPLATE = "default"

type Dockme struct {
	hostname    string
	domainname  string
	user        string
	cpuset      string
	image       string
	workdir     string
	macAddress  string
	name        string
	publish     string
	entrypoint  string // straying from type
	cmd         string // straying from type
	source      string // custom arg
	destination string // custom arg

	memory     int
	memorySwap int
	cpuShares  int

	rm          bool
	interactive bool
	tty         bool
	dryrun      bool // custom arg
	save        bool // custom arg

	expose      []string
	env         []string
	volume      []string // straying from type
	volumesFrom []string // straying from type
}

func (dm *Dockme) args() []string {
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

	if dm.rm {
		args = append(args, "--rm")
	}

	if dm.tty {
		args = append(args, "--tty")
	}

	if dm.interactive {
		args = append(args, "--interactive")
	}

	if len(dm.expose) > 0 {
		for _, e := range dm.expose {
			args = append(args, fmt.Sprintf("--expose=%s", e))
		}
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

func (dm *Dockme) exec() {
	args := dm.args()

	fmt.Printf("+ %s %s\n", DOCKER, strings.Join(args, " "))
	if dm.dryrun {
		fmt.Println("")
		os.Exit(0)
	}

	cmd := exec.Command(DOCKER, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func (dm *Dockme) applyDefaults() {
	// source is local dir if not
	if dm.source == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Errorf("%s", err)
			os.Exit(1)
		}
		dm.source = cwd
	}

	// desination is /src if not set
	if dm.destination == "" {
		dm.destination = "/src"
	}

	// hostname is image basename
	if dm.hostname == "" {
		imageSansTag := strings.Split(dm.image, ":")[0]
		imageSplit := strings.Split(imageSansTag, "/")
		dm.hostname = imageSplit[len(imageSplit)-1]
	}

	// cmd is bash
	if dm.cmd == "" {
		dm.cmd = "bash"
	}
}

var templates = map[string]*Dockme{
	"default": &Dockme{
		image:       "jmervine/zshrc:latest",
		name:        "dockme",
		interactive: true,
		tty:         true,
		rm:          true,
		cmd:         "bash",
	},
	"ruby": &Dockme{
		image:       "jmervine/herokudev-ruby:latest",
		name:        "rubydev",
		interactive: true,
		tty:         true,
		rm:          true,
		cmd:         "bash",
	},
	"rails": &Dockme{
		image:       "jmervine/herokudev-rails:latest",
		name:        "railsdev",
		interactive: true,
		tty:         true,
		rm:          true,
		cmd:         "bash",
	},
	"node": &Dockme{
		image:       "jmervine/herokudev-node:latest",
		name:        "nodedev",
		interactive: true,
		tty:         true,
		rm:          true,
		cmd:         "bash",
	},
	"nodebox": &Dockme{
		image:       "jmervine/nodebox:latest",
		name:        "nodeboxdev",
		interactive: true,
		tty:         true,
		rm:          true,
		cmd:         "bash",
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "Dockme"
	app.Version = "0.2.0"
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
    nodebox     nodebox template w/ 'jmervine/nodebox:latest'
    default     default template w/ 'jmervine/zshrc:latest'
    ruby        ruby template w/ 'jmervine/herokudev-ruby:latest'
    rails       rails template w/ 'jmervine/herokudev-rails:latest'
    node        node template w/ 'jmervine/herokudev-node:latest'
    help        Shows a list of commands or help for one command

`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "template, T",
			Usage: "set docker image template, see TEMPLATES below",
		},
		cli.StringFlag{
			Name:  "image, I",
			Usage: "set docker image",
		},
		cli.StringFlag{
			Name:  "source, s",
			Usage: "local source directory",
		},
		cli.StringFlag{
			Name:  "destination, d",
			Usage: "container source directory",
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
			Name:  "dryrun",
			Usage: "show docker command to be run",
		},
		cli.BoolFlag{
			Name:  "save, S",
			Usage: "save configuration to .dockmerc file",
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
			Name: "rm, r",
		},
		cli.BoolFlag{
			Name: "no-rm, k",
		},
		cli.BoolFlag{
			Name: "interactive, i",
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
		var template = DEFAULT_TEMPLATE
		if c.String("template") != "" {
			template = c.String("template")
		}

		dockme := templates[template]

		if c.String("image") != "" {
			dockme.image = c.String("image")
		}
		if c.String("source") != "" {
			dockme.source = c.String("source")
		}
		if c.String("destination") != "" {
			dockme.destination = c.String("destination")
		}
		if c.String("publish") != "" {
			dockme.publish = c.String("publish")
		}
		if c.String("workdir") != "" {
			dockme.workdir = c.String("workdir")
		}
		if c.String("name") != "" {
			dockme.name = c.String("name")
		}
		if c.String("entrypoint") != "" {
			dockme.entrypoint = c.String("entrypoint")
		}
		if c.String("user") != "" {
			dockme.user = c.String("user")
		}
		if c.String("hostname") != "" {
			dockme.hostname = c.String("hostname")
		}
		if c.String("domainname") != "" {
			dockme.domainname = c.String("domainname")
		}
		if c.String("mac-address") != "" {
			dockme.macAddress = c.String("mac-address")
		}
		if c.String("cpuset") != "" {
			dockme.cpuset = c.String("cpuset")
		}
		if c.Int("memory") > 0 {
			dockme.memory = c.Int("memory")
		}
		if c.Int("memory-swap") > 0 {
			dockme.memorySwap = c.Int("memory-swap")
		}
		if c.Bool("rm") {
			dockme.rm = c.Bool("rm")
		}
		if c.Bool("no-rm") {
			dockme.rm = false
		}
		if c.Bool("interactive") {
			dockme.interactive = c.Bool("interactive")
		}
		if c.Bool("no-interactive") {
			dockme.interactive = false
		}
		if c.Bool("tty") {
			dockme.tty = c.Bool("tty")
		}
		if c.Bool("no-tty") {
			dockme.tty = false
		}
		if c.Bool("dryrun") {
			dockme.dryrun = c.Bool("dryrun")
		}
		if c.Bool("save") {
			dockme.save = c.Bool("save")
		}
		if c.String("expose") != "" {
			dockme.expose = split(c.String("expose"))
		}
		if c.String("env") != "" {
			dockme.env = split(c.String("env"))
		}
		if c.String("volume") != "" {
			dockme.volume = split(c.String("volume"))
		}
		if c.String("volumes-from") != "" {
			dockme.volumesFrom = split(c.String("volumes-from"))
		}

		dockme.cmd = strings.Join(c.Args(), " ")

		dockme.applyDefaults()
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
