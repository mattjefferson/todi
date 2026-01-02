package app

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func runAuth(ctx context.Context, state *state, args []string) int {
	if len(args) == 0 {
		printAuthUsage(state.Out)
		return 2
	}
	switch args[0] {
	case "login":
		return runAuthLogin(ctx, state, args[1:])
	case "logout":
		return runAuthLogout(state)
	case "status":
		return runAuthStatus(state)
	case "-h", "--help", "help":
		printAuthUsage(state.Out)
		return 0
	default:
		fmt.Fprintln(state.Err, "error: unknown auth command:", args[0])
		printAuthUsage(state.Err)
		return 2
	}
}

func runAuthLogin(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist auth login", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(state.Err, "error:", err)
		return 2
	}
	if help {
		printAuthUsage(state.Out)
		return 0
	}
	if state.NoInput || !isTTY(os.Stdin) {
		fmt.Fprintln(state.Err, "error: login requires TTY (disable --no-input)")
		return 2
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprint(state.Err, "Todoist token: ")
	token, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(state.Err, "error:", err)
		return 1
	}
	token = strings.TrimSpace(token)
	if token == "" {
		fmt.Fprintln(state.Err, "error: token required")
		return 2
	}

	state.Config.Token = token
	if err := state.Config.Save(state.ConfigPath); err != nil {
		fmt.Fprintln(state.Err, "error:", err)
		return 1
	}
	fmt.Fprintln(state.Out, "token saved")
	_ = ctx
	return 0
}

func runAuthLogout(state *state) int {
	state.Config.Token = ""
	if err := state.Config.Save(state.ConfigPath); err != nil {
		fmt.Fprintln(state.Err, "error:", err)
		return 1
	}
	fmt.Fprintln(state.Out, "token cleared")
	return 0
}

func runAuthStatus(state *state) int {
	envToken := os.Getenv("TODOIST_TOKEN")
	if envToken != "" {
		fmt.Fprintln(state.Out, "token set (TODOIST_TOKEN)")
		return 0
	}
	if state.Config.Token != "" {
		fmt.Fprintln(state.Out, "token set (config)")
		return 0
	}
	fmt.Fprintln(state.Out, "token missing")
	return 3
}
