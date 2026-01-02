package app

import (
	"fmt"
	"io"
)

func printUsage(out io.Writer) {
	fmt.Fprintln(out, `todoist - Todoist tasks CLI

USAGE:
  todoist [global flags] <command> [args]

COMMANDS:
  task    Manage tasks
  auth    Manage auth token
  config  Manage config

GLOBAL FLAGS:
  -h, --help        Show help
  --version         Show version
  -q, --quiet        Less output
  -v, --verbose      Verbose output
  --json            JSON output
  --plain           Plain output
  --no-input        Disable prompts
  --no-color        Disable color
  --config <path>   Config path override
  --api-base <url>  API base (default https://api.todoist.com)
  --label-cli       Add label 'cli' to created tasks

NOTES:
  Task identifiers accept exact task titles unless --id is set.
  Project references use exact project titles.
`)
}

func printTaskUsage(out io.Writer) {
	fmt.Fprintln(out, `todoist task - task commands

USAGE:
  todoist task list [project_title]
  todoist task get <task>
  todoist task add <content>
  todoist task update <task>
  todoist task close <task>
  todoist task reopen <task>
  todoist task delete <task>
  todoist task quick <text>

NOTES:
  <task> accepts exact task title unless --id is set.
`)
}

func printAuthUsage(out io.Writer) {
	fmt.Fprintln(out, `todoist auth - auth commands

USAGE:
  todoist auth login
  todoist auth logout
  todoist auth status
`)
}

func printConfigUsage(out io.Writer) {
	fmt.Fprintln(out, `todoist config - config commands

USAGE:
  todoist config get <key>
  todoist config set <key> <value>
  todoist config path
  todoist config view

KEYS:
  token
  api_base
  default_project
  default_labels
  label_cli
`)
}
