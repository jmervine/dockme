# dockme
Simple docker wrapper for quickly spooling up containerized development environments.

## Install

> Install Docker, of course
>
> * https://docs.docker.com/installation/

```text
go get github.com/jmervine/dockme
go install github.com/jmervine/dockme
```

Or see binaries in the `builds` directory.

## Usage

```text
NAME:
    Dockme - Simple wrapper for quickly spooling up docker containers for development.

USAGE:
    Dockme [arguments...] [command]

VERSION:
    0.2.1

AUTHOR:
    Joshua Mervine

OPTIONS:
    Only custom options or options whose usage strays from dockers
    usage have help messages. All other options map directly to docker
    run options, see Docker help and documentation for details.

    --template, -T 		set docker image template, see TEMPLATES below
    --image, -i 		set docker image [required]
    --source, -s 		local source directory
    --destination, -d 		[/src] container source directory
    --publish, -p 		list of ports to publish
    --workdir, -w 		set container workdir
    --dryrun, -D		show docker command to be run
    --save, -S			save configuration to file
    --config, -C "Dockme.yml"	conifguration file path
    --expose, -E 		list of ports to expose
    --env, -e 			list of environments
    --volume, -V 		list of volume mounts
    --volumes-from 		list of containers to mount volumes from
    --name, -n
    --sudo			run Docker with sudo
    --rm, -r
    --no-rm, -k
    --interactive, -I
    --no-interactive, -x
    --tty, -t
    --no-tty, -N
    --entrypoint
    --user
    --hostname
    --domainname
    --mac-address
    --cpuset
    --memory
    --memory-swap
    --help, -h			show help
    --version, -v		print the version

TEMPLATES:
    nodebox    nodebox template w/ 'jmervine/nodebox:latest'
    ruby       ruby template w/ 'jmervine/herokudev-ruby:latest'
    rails      rails template w/ 'jmervine/herokudev-rails:latest'
    node       node template w/ 'jmervine/herokudev-node:latest'
    python2    python template w/ 'python:2-slim'
    python3    python template w/ 'python:3-slim'
    help       Shows a list of commands or help for one command

```

