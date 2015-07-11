package main

import (
	"github.com/jmervine/dockme"
	"github.com/jmervine/dockme/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/jmervine/dockme/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "Dockme"
	app.Version = dockme.VERSION
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
			Value: dockme.DEFAULT_CONFIG,
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
			Name:  "env-file",
			Usage: "list of environment files",
		},
		cli.StringFlag{
			Name:  "link, l",
			Usage: "list of links",
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
			Name: "hostname, H",
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
		template := dockme.DEFAULT_TEMPLATE
		if c.String("template") != "" {
			template = c.String("template")
		}
		var dm *dockme.Dockme

		config := c.String("config")

		read, err := ioutil.ReadFile(config)
		if err == nil {
			if err = yaml.Unmarshal(read, &dm); err != nil {
				log.Fatalf("error: %v", err)
			}
		} else {
			dm = dockme.Templates[template]
		}

		dm.Config = config

		if c.String("image") != "" {
			dm.Image = c.String("image")
		}

		// ensure image at least
		if dm.Image == "" {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		if c.String("source") != "" {
			dm.Source = c.String("source")
		}
		if c.String("destination") != "" {
			dm.Destination = c.String("destination")
		}
		if c.String("workdir") != "" {
			dm.Workdir = c.String("workdir")
		}
		if c.String("name") != "" {
			dm.Name = c.String("name")
		}
		if c.String("entrypoint") != "" {
			dm.Entrypoint = c.String("entrypoint")
		}
		if c.String("user") != "" {
			dm.User = c.String("user")
		}
		if c.String("hostname") != "" {
			dm.Hostname = c.String("hostname")
		}
		if c.String("domainname") != "" {
			dm.Domainname = c.String("domainname")
		}
		if c.String("mac-address") != "" {
			dm.MacAddress = c.String("mac-address")
		}
		if c.String("cpuset") != "" {
			dm.Cpuset = c.String("cpuset")
		}
		if c.Int("memory") > 0 {
			dm.Memory = c.Int("memory")
		}
		if c.Int("memory-swap") > 0 {
			dm.MemorySwap = c.Int("memory-swap")
		}
		if c.Bool("sudo") {
			dm.Sudo = c.Bool("sudo")
		}
		if c.Bool("rm") {
			dm.Rm = c.Bool("rm")
		}
		if c.Bool("no-rm") {
			dm.Rm = false
		}
		if c.Bool("interactive") {
			dm.Interactive = c.Bool("interactive")
		}
		if c.Bool("no-interactive") {
			dm.Interactive = false
		}
		if c.Bool("tty") {
			dm.Tty = c.Bool("tty")
		}
		if c.Bool("no-tty") {
			dm.Tty = false
		}
		if c.Bool("dryrun") {
			dm.Dryrun = c.Bool("dryrun")
		}
		if c.Bool("save") {
			dm.Save = true
		}
		if c.String("expose") != "" {
			dm.Expose = split(c.String("expose"))
		}
		if c.String("env") != "" {
			dm.Env = split(c.String("env"))
		}
		if c.String("env-file") != "" {
			dm.EnvFile = split(c.String("env-file"))
		}
		if c.String("link") != "" {
			dm.Links = split(c.String("link"))
		}
		if c.String("volume") != "" {
			dm.Volume = split(c.String("volume"))
		}
		if c.String("volumes-from") != "" {
			dm.VolumesFrom = split(c.String("volumes-from"))
		}
		if c.String("publish") != "" {
			dm.Publish = split(c.String("publish"))
		}

		cmd := strings.Join(c.Args(), " ")
		if cmd != "" {
			dm.Cmd = cmd
		}

		dm.ApplyDefaults()

		if dm.Save {
			dm.SaveConfig()
		}

		dm.Exec()
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
