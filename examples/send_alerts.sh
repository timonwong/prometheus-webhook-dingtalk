#!/bin/bash

SOURCE="${BASH_SOURCE[0]}"
while [[ -h "$SOURCE" ]]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ ${SOURCE} != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"

# demoAlerts file: /web/ui/react-app/src/pages/PlaygroundDemoAlert.json
demoAlerts=$(<$DIR/../web/ui/react-app/src/pages/PlaygroundDemoAlert.json)

curl -XPOST -d"$demoAlerts" http://localhost:8060/dingtalk/webhook1/send
