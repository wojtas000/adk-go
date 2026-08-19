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
	"strings"
	"testing"
)

// Under `go test`, ADK Go is the main module rather than a recorded dependency,
// so no module version is available and resolveVersion falls back to devVersion.
func TestVersionFallsBackWhenNoModuleVersion(t *testing.T) {
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
