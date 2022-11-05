// Copyright (c) 2015-2022 The Decred developers
// Copyright (c) 2022 Dave Collins
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
)

const (
	// semanticAlphabet defines the allowed characters for the pre-release and
	// build metadata portions of a semantic version string.
	semanticAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-."
)

// semverRE is a regular expression used to parse a semantic version string into
// its constituent parts.
var semverRE = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)` +
	`(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*` +
	`[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

// These variables define the application version and follow the semantic
// versioning 2.0.0 spec (https://semver.org/).
var (
	// Version is the application version per the semantic versioning 2.0.0 spec
	// (https://semver.org/).
	//
	// It is defined as a variable so it can be overridden during the build
	// process with:
	// '-ldflags "-X main.Version=fullsemver"'
	// if needed.
	//
	// It MUST be a full semantic version per the semantic versioning spec or
	// the app will panic at runtime.  Of particular note is the pre-release
	// and build metadata portions MUST only contain characters from
	// semanticAlphabet.
	Version = "1.0.1-pre"
)

// parseUint32 converts the passed string to an unsigned integer or returns an
// error if it is invalid.
func parseUint32(s string, fieldName string) (uint32, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("malformed semver %s: %w", fieldName, err)
	}
	return uint32(val), err
}

// checkSemString returns an error if the passed string contains characters that
// are not in the provided alphabet.
func checkSemString(s, alphabet, fieldName string) error {
	for _, r := range s {
		if !strings.ContainsRune(alphabet, r) {
			return fmt.Errorf("malformed semver %s: %q invalid", fieldName, r)
		}
	}
	return nil
}

// parseSemVer parses various semver components from the provided string.
func parseSemVer(s string) (uint32, uint32, uint32, string, string, error) {
	// Parse the various semver component from the version string via a regular
	// expression.
	m := semverRE.FindStringSubmatch(s)
	if m == nil {
		err := fmt.Errorf("malformed version string %q: does not conform to "+
			"semver specification", s)
		return 0, 0, 0, "", "", err
	}

	major, err := parseUint32(m[1], "major")
	if err != nil {
		return 0, 0, 0, "", "", err
	}

	minor, err := parseUint32(m[2], "minor")
	if err != nil {
		return 0, 0, 0, "", "", err
	}

	patch, err := parseUint32(m[3], "patch")
	if err != nil {
		return 0, 0, 0, "", "", err
	}

	preRel := m[4]
	err = checkSemString(preRel, semanticAlphabet, "pre-release")
	if err != nil {
		return 0, 0, 0, "", "", err
	}

	build := m[5]
	err = checkSemString(build, semanticAlphabet, "buildmetadata")
	if err != nil {
		return 0, 0, 0, "", "", err
	}

	return major, minor, patch, preRel, build, nil
}

// vcsInfo attempts to return the version control system short commit hash
// that was used to build the binary and whether or not the working tree is
// modified.  It currently only detects git commits.
func vcsInfo() (string, bool) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	var vcs, revision string
	var isModified bool
	for _, bs := range bi.Settings {
		switch bs.Key {
		case "vcs":
			vcs = bs.Value
		case "vcs.revision":
			revision = bs.Value
		case "vcs.modified":
			isModified = bs.Value == "true"
		}
	}
	if vcs == "" {
		return "", isModified
	}
	if vcs == "git" && len(revision) > 9 {
		revision = revision[:9]
	}
	return revision, isModified
}

func init() {
	major, minor, patch, preRel, buildMetadata, err := parseSemVer(Version)
	if err != nil {
		panic(err)
	}
	if buildMetadata == "" {
		revision, isModified := vcsInfo()
		if revision != "" {
			Version = fmt.Sprintf("%d.%d.%d", major, minor, patch)
			if preRel != "" {
				Version += "-" + preRel
			}
			Version += "+" + revision
			if isModified {
				Version += ".modified"
			}
		}
	}
}
