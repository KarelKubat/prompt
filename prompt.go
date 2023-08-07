package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/KarelKubat/flagnames"
	"github.com/fatih/color"
)

var (
	hostFlag        = flag.String("host", "", "hostname, default: don't include")
	dateFlag        = flag.Bool("date", true, "include date/time")
	gitFlag         = flag.Bool("git", true, "include git clean/unclean")
	cwdFlag         = flag.Bool("cwd", true, "include last part of current working directory")
	g4clientFlag    = flag.Bool("g4client", true, "include g4 client id")
	appendFlag      = flag.String("append", "\\n", "closer for prompt, use \\\\n for newlines")
	prependFlag     = flag.String("prepend", "\\n", "opener for prompt, use \\\\n for newlines")
	alwaysColorFlag = flag.Bool("alwayscolor", false, "when true, colorize even if `stdout` is a pipe")
)

func main() {
	flagnames.Patch()
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `
Usage: prompt - outputs nicely readable prompt to stdout.
Use in .bashrc as follows:
  PS1=""
  PROMPT_COMMAND="prompt -host MYHOST -alwayscolor"
Use in .fish as follows:q
  function fish_prompt
    prompt -host MYHOST -alwayscolor
  end

Useful flags:`)
		flag.PrintDefaults()
	}
	flag.Parse()

	// Always colorize?
	if *alwaysColorFlag {
		color.NoColor = false
	}

	prompt := escape(*prependFlag)

	// Hostname
	prompt = add(prompt, color.New(color.FgHiCyan, color.Underline), *hostFlag)

	// Date/time
	if *dateFlag {
		prompt = add(prompt, color.New(color.FgHiYellow), time.Now().Format("2006-01-02 15:04:05"))
	}

	// Git status
	if *gitFlag {
		prompt = add(prompt, color.New(color.FgRed), gitStatus())
	}

	// CWD
	if *cwdFlag {
		prompt = add(prompt, color.New(color.FgHiGreen), cwd())
	}

	// g4 client
	if *g4clientFlag {
		prompt = add(prompt, color.New(color.FgHiBlue), g4client())
	}

	// Closer
	prompt = add(prompt, color.New(color.FgWhite), escape(*appendFlag))

	fmt.Print(prompt)
}

func add(prompt string, col *color.Color, addition string) string {
	if addition == "" {
		return prompt
	}
	if prompt != "" && prompt[len(prompt)-1] != '\n' {
		prompt += " "
	}
	prompt += col.Sprint(addition)
	return prompt
}

func runCmd(c []string) (string, error) {
	cmd := exec.Command(c[0], c[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func gitStatus() string {
	if _, err := os.Stat(".git"); err != nil {
		return ""
	}
	gitStatus, err := runCmd([]string{"git", "status", "-s"})
	if err != nil {
		return ""
	}
	out := "git"
	if len(strings.Split(gitStatus, "\n")) > 1 {
		out += " unclean"
	}
	return out
}

func cwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("cwd:%v", err)
	}
	parts := strings.Split(dir, "/")
	return parts[len(parts)-1]
}

func g4client() string {
	client, err := runCmd([]string{"g4clients", "--current"})
	if err != nil {
		return ""
	}
	out := strings.Split(client, "\n")[0]
	if out != "" {
		out = fmt.Sprintf("[ %v ]", out)
	}
	return out
}

func escape(s string) string {
	out := ""
	for len(s) > 0 {
		if strings.HasPrefix(s, "\\n") {
			out += "\n"
			s = s[2:]
		} else {
			out += s[0:1]
			s = s[1:]
		}
	}
	return out
}
