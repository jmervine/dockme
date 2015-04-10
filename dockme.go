package main

// this a first pass, expect refactors

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
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
	publish     string
	entrypoint  string // straying from type
	cmd         string // straying from type
	source      string // custom arg
	destination string // custom arg
	template    string // custom arg

	memory     int
	memorySwap int
	cpuShares  int

	rm          bool
	interactive bool
	tty         bool
	dryrun      bool // custom arg
	save        bool // custom arg

	env         []string
	expose      []string // straying from type
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
		fmt.Errorf("Can't specify source or destination, both or neither")
		os.Exit(1)
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

	if len(dm.env) > 0 {
		for _, e := range dm.env {
			args = append(args, fmt.Sprintf("--env=%s", e))
		}
	}

	if len(dm.expose) > 0 {
		for _, e := range dm.expose {
			args = append(args, fmt.Sprintf("--expose=%s", e))
		}
	}

	if len(dm.volume) > 0 {
		for _, v := range dm.volume {
			args = append(args, fmt.Sprintf("--volume=%s", v))
		}
	}

	if len(dm.volumesFrom) > 0 {
		for _, e := range dm.volumesFrom {
			args = append(args, fmt.Sprintf("--volumes-from=%s", e))
		}
	}

	if dm.image == "" {
		fmt.Errorf("Image is required")
		os.Exit(1)
	}
	args = append(args, dm.image)
	args = append(args, dm.cmd)

	return args
}

func (dm *dockme) exec() {
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

func (dm *dockme) applyDefaults() {
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

var templates = map[string]*dockme{
	"default": &dockme{
		image:       "jmervine/zshrc:latest",
		name:        "dockme",
		interactive: true,
		tty:         true,
		rm:          true,
		cmd:         "bash",
	},
	"ruby": &dockme{
		image:       "ruby:latest",
		name:        "rubydev",
		interactive: true,
		tty:         true,
		rm:          true,
		cmd:         "bash",
	},
	"node": &dockme{
		image:       "node:latest",
		name:        "nodedev",
		interactive: true,
		tty:         true,
		rm:          true,
		cmd:         "bash",
	},
	"nodebox": &dockme{
		image:       "jmervine/nodebox:latest",
		name:        "nodeboxdev",
		interactive: true,
		tty:         true,
		rm:          true,
		cmd:         "bash",
	},
}

func main() {
	template := "default"

	app := cli.NewApp()
	app.Name = "dockme"
	app.Version = "0.2.0"
	app.Author = "Joshua Mervine"
	app.Email = "joshua@mervine.net"
	app.Usage = "Simple wrapper for quickly spooling up docker containers for development."

	cli.AppHelpTemplate = `NAME:
    {{.Name}} - {{.Usage}}

USAGE:
    {{.Name}} [template] [arguments...] [command]

VERSION:
    {{.Version}}

AUTHOR:
    {{.Author}}

TEMPLATES:
    {{range .Commands}}{{ .Name }}{{ "\t" }}{{.Usage}}
    {{end}}{{if .Flags}}

OPTIONS:
    Only custom options or options whose usage strays from dockers
    usage have help messages. All other options map directly to docker
    run options, see Docker help and documentation for details.

    {{range .Flags}}{{.}}
    {{end}}{{end}}
`

	app.Flags = []cli.Flag{
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

	app.Commands = []cli.Command{}
	for k, t := range templates {
		app.Commands = append(app.Commands, cli.Command{
			Name:  k,
			Usage: fmt.Sprintf("%s template w/ '%s'", k, t.image),
			Action: func(c *cli.Context) {
				template = c.Args().First()
			},
		})
	}

	split := func(s string) []string {
		sp := strings.Split(s, ",")
		if len(sp) == 1 {
			sp = strings.Split(s, " ")
		}

		for n, i := range sp {
			sp[n] = strings.TrimSpace(i)
		}

		return sp
	}

	app.Action = func(c *cli.Context) {
		dm := templates[template]
		if c.String("image") != "" {
			dm.image = c.String("image")
		}
		if c.String("source") != "" {
			dm.source = c.String("source")
		}
		if c.String("destination") != "" {
			dm.destination = c.String("destination")
		}
		if c.String("publish") != "" {
			dm.publish = c.String("publish")
		}
		if c.String("name") != "" {
			dm.name = c.String("name")
		}
		if c.String("entrypoint") != "" {
			dm.entrypoint = c.String("entrypoint")
		}
		if c.String("user") != "" {
			dm.user = c.String("user")
		}
		if c.String("hostname") != "" {
			dm.hostname = c.String("hostname")
		}
		if c.String("domainname") != "" {
			dm.domainname = c.String("domainname")
		}
		if c.String("mac-address") != "" {
			dm.macAddress = c.String("mac-address")
		}
		if c.String("cpuset") != "" {
			dm.cpuset = c.String("cpuset")
		}
		if c.Int("memory") > 0 {
			dm.memory = c.Int("memory")
		}
		if c.Int("memory-swap") > 0 {
			dm.memorySwap = c.Int("memory-swap")
		}
		if c.Bool("rm") {
			dm.rm = c.Bool("rm")
		}
		if c.Bool("no-rm") {
			dm.rm = false
		}
		if c.Bool("interactive") {
			dm.interactive = c.Bool("interactive")
		}
		if c.Bool("no-interactive") {
			dm.interactive = false
		}
		if c.Bool("tty") {
			dm.tty = c.Bool("tty")
		}
		if c.Bool("no-tty") {
			dm.tty = false
		}
		if c.Bool("dryrun") {
			dm.dryrun = c.Bool("dryrun")
		}
		if c.Bool("save") {
			dm.save = c.Bool("save")
		}
		if c.String("expose") != "" {
			dm.expose = split(c.String("expose"))
		}
		if c.String("env") != "" {
			dm.env = split(c.String("env"))
		}
		if c.String("volume") != "" {
			dm.volume = split(c.String("volume"))
		}
		if c.String("volumes-from") != "" {
			dm.volumesFrom = split(c.String("volumes-from"))
		}

		dm.applyDefaults()
		dm.exec()
	}

	app.Run(os.Args)
}
