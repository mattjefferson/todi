// Command todoist runs the Todoist CLI.
package main

import (
	"os"

	"github.com/mattjefferson/todoist-cli/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
