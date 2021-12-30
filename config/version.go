package config

import "strings"

const Version = "v0.6.0-beta-20211229"

func GetVersion() (v string) {
	if strings.HasPrefix(Version, "v") {
		return Version
	}
	return "v" + Version
}
