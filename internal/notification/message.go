package notification

import (
	"fmt"
	"strings"
)

const divider = "=============================="

func section(sb *strings.Builder, heading string, body ...string) {
	fmt.Fprintf(sb, "%s\n%s\n%s\n\n", heading, divider, strings.Join(body, "\n"))
}
