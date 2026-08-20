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
	"testing"
)

func TestVersionFrom(t *testing.T) {
	tests := []struct {
		name string
		info *debug.BuildInfo
		want string
	}{
		{
			name: "imported as a dependency reports the resolved version",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/app", Version: "(devel)"},
				Deps: []*debug.Module{
					{Path: "example.com/other", Version: "v1.0.0"},
					{Path: modulePath, Version: "v2.2.0"},
				},
			},
			want: "2.2.0",
		},
		{
			// `go install google.golang.org/adk/v2/cmd/adkgo@v2.2.0`: adk owns the
			// binary, so it is the main module and absent from Deps.
			name: "installed as the main module reports the released version",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: modulePath, Version: "v2.2.0"},
			},
			want: "2.2.0",
		},
		{
			name: "local build of this repo reports its pseudo-version",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: modulePath, Version: "v2.2.1-0.20260820082125-15f3438e60f6+dirty"},
			},
			want: "2.2.1-0.20260820082125-15f3438e60f6+dirty",
		},
		{
			// What `go test` and `go run` inside the repo record.
			name: "main module without a recorded version falls back",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: modulePath, Version: "(devel)"},
			},
			want: devVersion,
		},
		{
			name: "adk absent entirely falls back",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/app", Version: "v1.0.0"},
				Deps: []*debug.Module{{Path: "example.com/other", Version: "v1.0.0"}},
			},
			want: devVersion,
		},
		{
			name: "replace directive leaves an empty version and falls back",
			info: &debug.BuildInfo{
				Main: debug.Module{Path: "example.com/app", Version: "(devel)"},
				Deps: []*debug.Module{{Path: modulePath, Version: ""}},
			},
			want: devVersion,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := versionFrom(tc.info); got != tc.want {
				t.Errorf("versionFrom() = %q, want %q", got, tc.want)
			}
		})
	}
}

// Under `go test`, ADK Go is the main module and Go records its version as
// "(devel)", so the resolved value is the fallback.
func TestResolveVersionUnderTest(t *testing.T) {
	if got := resolveVersion(); got != devVersion {
		t.Errorf("resolveVersion() = %q, want fallback %q", got, devVersion)
	}
}

// Whatever the build mode, the reported version must be a non-empty string with
// no leading "v", since it is embedded directly in the google-adk/<version>
// header and reported as the OTel and MCP client version.
func TestVersionIsUsable(t *testing.T) {
	if Version == "" {
		t.Fatal("Version is empty")
	}
	if strings.HasPrefix(Version, "v") {
		t.Errorf("Version = %q, must not carry a leading %q", Version, "v")
	}
}
