package main

import (
	"fmt"

	fwew "github.com/fwew/fwew-lib/v5"
)

type version struct {
	Major, Minor, Patch int
	Label               string
}

// Version information
var Version = version{
	1, 7, 1,
	"",
}

func (v *version) String() string {
	var label string
	if v.Label != "" {
		label = "-" + v.Label
	}
	return fmt.Sprintf("discord bot: %d.%d.%d%s\n%s", v.Major, v.Minor, v.Patch, label, fwew.Version.String())
}
