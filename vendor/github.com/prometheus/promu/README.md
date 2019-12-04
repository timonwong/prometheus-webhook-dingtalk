# Prometheus Utility Tool [![Build Status](https://travis-ci.org/prometheus/promu.svg)][travis]

[![CircleCI](https://circleci.com/gh/prometheus/promu/tree/master.svg?style=shield)][circleci]

## Usage

```help
usage: promu [<flags>] <command> [<args> ...]

promu is the utility tool for building and releasing Prometheus projects

Flags:
  -h, --help                 Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=".promu.yml"  Path to config file
  -v, --verbose              Verbose output

Commands:
  help [<command>...]
    Show help.

  build [<flags>] [<binary-names>...]
    Build a Go project

  check licenses [<flags>] [<location>...]
    Inspect source files for each file in a given directory

  checksum [<location>...]
    Calculate the SHA256 checksum for each file in the given location

  crossbuild [<flags>] [<tarballs>]
    Crossbuild a Go project using Golang builder Docker images

  info
    Print info about current project and exit

  release [<flags>] [<location>...]
    Upload all release files to the Github release

  tarball [<flags>] [<location>...]
    Create a tarball from the built Go project

  version [<flags>]
    Print the version and exit

```

## `.promu.yml` config file

See documentation example [here](doc/examples/prometheus/.promu.yml)

## Compatibility

* Go 1.6+

## More information

* This tool is part of our reflexion about [Prometheus component Builds](https://docs.google.com/document/d/1Ql-f_aThl-2eB5v3QdKV_zgBdetLLbdxxChpy-TnWSE)
* All of the core developers are accessible via the [Prometheus Developers Mailinglist](https://groups.google.com/forum/?fromgroups#!forum/prometheus-developers) and the `#prometheus` channel on `irc.freenode.net`.

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md)

## License

Apache License 2.0, see [LICENSE](LICENSE).

[circleci]: https://circleci.com/gh/prometheus/promu
[travis]: https://travis-ci.org/prometheus/promu
