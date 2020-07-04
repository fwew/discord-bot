package main

import (
	"fmt"
	fwew "github.com/fwew/fwew-lib/v5"
)

type version struct {
	Major, Minor, Patch int
	Label               string
}

var Version = version{
	1, 1, 3,
	"dev",
}

func (v *version) String() string {
	return fmt.Sprintf("discord bot: %d.%d.%d-%s\n%s", v.Major, v.Minor, v.Patch, v.Label, fwew.Version.String())
}
