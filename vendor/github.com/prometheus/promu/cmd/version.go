// Copyright © 2016 Prometheus Team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/prometheus/common/version"
)

var (
	versioncmd = app.Command("version", "Print the version and exit")
	short      = versioncmd.Flag("short", "Print shorter version").Short('s').Bool()
)

func runVersion() {
	if *short {
		fmt.Printf(version.Version)
	} else {
		fmt.Println(version.Print("promu"))
	}
}
