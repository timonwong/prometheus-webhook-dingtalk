#!/bin/bash
set -e

_init() {
    # Save release LDFLAGS
    # LDFLAGS=$(go run scripts/gen-ldflags.go)
    LDFLAGS=

    # List of supported architectures
    SUPPORTED_OSARCH='linux/amd64 darwin/amd64'
}

go_build() {
    local package=$1
    local osarch=$2
    os=$(echo $osarch | cut -f1 -d'/')
    arch=$(echo $osarch | cut -f2 -d'/')
    echo -n "-->"
    printf "%15s:%s\n" "${osarch}" "${package}"

    # Release binary name
    release_bin=".build/$os-$arch/$(basename $package)"

    # Go build to build the binary.
    GOOS=$os GOARCH=$arch go build --ldflags "${LDFLAGS}" -o $release_bin "$package"
}

main() {
    # Build releases.
    echo "Executing builds for OS: ${SUPPORTED_OSARCH}"

    for package in ./cmd/*; do
        for each_osarch in ${SUPPORTED_OSARCH}; do
            go_build "$package" ${each_osarch}
        done
    done
}

_init && main
