# dockme
Simple docker wrapper for quickly spooling up containerized development environments.

## Install

> Install Docker, of course
>
> * https://docs.docker.com/installation/

```text
mkdir -p ~/.bin
curl -sSL https://raw.githubusercontent.com/jmervine/dockme/master/dockme > ~/.bin/dockme
chmod +x ~/.bin/dockme
echo "export PATH=~/.bin:$PATH" >> ~/.bashrc
source ~/.bashrc

dockme --help
```

## Usage

```text
Usage: dockme [template] [options] [-- command]

Simple wraper for quickly spooling up docker comtainers a development
environments.

  Templates:
  - vim     jmervine/vimrc:latest
  - node    node:latest
  - ruby    ruby:latest
  - rails   rails:latest
  - python  python:2
  - golang  golang:latest

  Template configurations must be defined as the first argument and are applied
  and then overwritten by any passed arguments.

  Options:
  -s, --source         local host source directory
  -d, --destination    remote host source directory
  -i, --image          image
  -H, --hostname       see 'docker run --help' for details
  -p, --port           see 'docker run --help' for details
  -w, --workdir        see 'docker run --help' for details
  -w, --workdir        see 'docker run --help' for details
  -n, --name           see 'docker run --help' for details
  -N, --net            see 'docker run --help' for details
  -r, --rm             see 'docker run --help' for details
      --cpuset         see 'docker run --help' for details
      --memory         see 'docker run --help' for details
  -v, --volumes        see 'docker run --help' for details
      --volumes-from   see 'docker run --help' for details
  -D, --dryrun         only print what would be executed
      --save           save configuration to local .dockmerc file

  Defaults:
  - source       '/Users/jmervine/Development/dockme'
  - destination  '/src'
  - image        'jmervine/vimrc'
  - rm           'true'
  - workdir      destination directory

External Configuration

dockme can be configured via an external file. It will look for a file named
.dockmerc in your home directory and the current working directory, unless one
is passed via the configuration argument. *All arguments passed will supersede
arguments configured by dockmerc files.*

  Configuration Precedence

  > ARGUMENTS
  --> .dockmerc
  ----> ~/.dockmerc

  Configuration Sample:

  # file: .dockmerc
  _source=/path/to/source
  _destination=/path/to/dest
  _image=username/image:tag
  _rm=true
  _workdir=/path/to/workdir
  _name=container_name
  _net=host
  _cpuset=2
  _memory=1g
  _volumes=/some/vol:/some/path
  _volumes_from=volume_container
  _command=echo a command

```

## Examples

```text
# defaults
$ dockme
docker run -it --rm --workdir=/src --volume=/path/to/project:/src jmervine/vimrc

# custom examples
$ pwd
/Users/jmervine/Development/dockme
$ dockme -i ruby:latest -- irb
+ docker run -it --workdir=/src --rm --volume=/Users/jmervine/Development/dockme:/src ruby:latest irb
irb(main):001:0> `pwd`.strip
=> "/src"
irb(main):002:0> RUBY_VERSION
=> "2.2.1"
irb(main):003:0>

# template example
$ cd /path/to/project
$ dockme ruby
docker run -it --rm --workdir=/src --volume=/path/to/project:/src ruby:latest
```
