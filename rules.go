package main

import (
	//std necessities
	"runtime"
)

type Rule struct {
	Action   string `json:"action"`
	Features *struct {
		HasCustomResolution *bool `json:"has_custom_resolution"`
		IsDemoUser          *bool `json:"is_demo_user"`
	} `json:"features"`
	OS *struct {
		Name    *string `json:"name"`
		Arch    *string `json:"arch"`
		Version *string `json:"version"` //TODO: Implement for Windows 10
	} `json:"os"`
}

func (rule *Rule) Valid() bool {
	if rule.OS != nil {
		if rule.OS.Name != nil {
			switch *rule.OS.Name {
			case "windows":
				if runtime.GOOS == "windows" {
					return rule.Match(true)
				}
			case "linux":
				if runtime.GOOS == "linux" {
					return rule.Match(true)
				}
			case "osx", "macos":
				if runtime.GOOS == "darwin" {
					return rule.Match(true)
				}
			}
		}

		if rule.OS.Arch != nil {
			switch *rule.OS.Arch {
			case "x86":
				if runtime.GOARCH == "386" {
					return rule.Match(true)
				}
			case "x64":
				if runtime.GOARCH == "amd64" {
					return rule.Match(true)
				}
			}
		}
	}

	if rule.Features == nil && rule.OS == nil {
		return rule.Match(true)
	}

	return rule.Match(false)
}

func (rule *Rule) Match(matches bool) bool {
	if rule.Action == "allow" {
		return matches
	}
	return !matches
}
