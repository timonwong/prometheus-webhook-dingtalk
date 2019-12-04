## v0.5.0 / 2019-06-21

* [CHANGE] Remove --broken from git describe. #145
* [FEATURE] Add support for aix/ppc64. #151
* [ENHANCEMENT] cmd/release: add --timeout option. #142
* [ENHANCEMENT] cmd/release: create release in GitHub if none exists. #148

## 0.4.0 / 2019-05-03

* [FEATURE] Fallback to `git describe` output if no VERSION. #130
* [BUGFIX] cmd/tarball: restore --prefix flag. #133
* [BUGFIX] cmd/release: don't leak credentials in case of error. #136

## 0.3.0 / 2019-02-18

* [FEATURE] Make extldflags extensible by configuration. #125
* [ENHANCEMENT] Avoid bind-mounting to allow building with a remote docker engine #95

## 0.2.0 / 2018-11-07

* [FEATURE] Adding changes to support s390x
* [FEATURE] Add option to disable static linking
* [FEATURE] Add support for 32bit MIPS.
* [FEATURE] Added check_licenses Command to Promu
* [ENHANCEMENT] Allow to customize nested options via env variables
* [ENHANCEMENT] Bump Go version to 1.11
* [ENHANCEMENT] Add warning if promu info is unable to determine repo info
* [BUGFIX] Fix build on SmartOS by not setting gcc's -static flag
* [BUGFIX] Fix git repository url parsing

## 0.1.0 / 2017-09-22

Initial release
