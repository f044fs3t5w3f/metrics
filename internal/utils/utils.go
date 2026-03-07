package utils

import (
	"fmt"
	"io"
)

func PrintBuildInfo(w io.Writer, buildVersion, buildDate, buildCommit string) {
	v := func(val string) string {
		if val != "" {
			return val
		}
		return "N/A"
	}
	fmt.Fprintf(w, "Build version: %s\n", v(buildVersion))
	fmt.Fprintf(w, "Build date: %s\n", v(buildDate))
	fmt.Fprintf(w, "Build commit: %s\n", v(buildCommit))
}
