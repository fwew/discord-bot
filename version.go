package main

import (
	"fmt"
	fwew "github.com/fwew/fwew-lib"
)

type version struct {
	Major, Minor, Patch int
	Label               string
	Name                string
}

var Version = version{
	0, 0, 1,
	"dev",
	"",
}

func (v *version) String() string {
	return fmt.Sprintf("discord bot: %d.%d.%d-%s \"%s\"\n%s", v.Major, v.Minor, v.Patch, v.Label, v.Name, fwew.Version.String())
}
