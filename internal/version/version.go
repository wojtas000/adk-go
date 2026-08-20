// Copyright 2025 Google LLC
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

package version

import (
	"runtime/debug"
	"strings"
)

const modulePath = "google.golang.org/adk/v2"
const devVersion = "2.x-dev"

var Version = resolveVersion()

func resolveVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return devVersion
	}
	return versionFrom(info)
}

// versionFrom reads the ADK Go version out of build info, without the leading
// "v", or returns devVersion when the build recorded none.
func versionFrom(info *debug.BuildInfo) string {
	// Imported as a library: the consuming build records our version here.
	for _, dep := range info.Deps {
		if dep.Path == modulePath && dep.Version != "" {
			return strings.TrimPrefix(dep.Version, "v")
		}
	}
	// Built as the main module, so we are not one of our own dependencies. This
	// covers `go install google.golang.org/adk/v2/cmd/adkgo@vX.Y.Z`, which
	// reports the released version, and a local build of this repo, which
	// reports a pseudo-version. Go writes "(devel)" when it recorded none.
	if info.Main.Path == modulePath && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return strings.TrimPrefix(info.Main.Version, "v")
	}
	return devVersion
}
