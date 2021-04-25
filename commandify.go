package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type CommandifyOptions struct {
	RunAsCommand string `yaml:"run-as-command,omitempty"`
}

type DockerComposerData struct {
	CommandifyOptions *CommandifyOptions `yaml:"x-commandify,omitempty"`
}

func printUseAndExit() {
	fmt.Println("commandify [options] <PATH_TO_DOCKER-COMPOSE.YML>")
	os.Exit(1)
}

func getAdditionalArgs(o *CommandifyOptions) []string {
	if o != nil {
		return []string{
			"run",
			o.RunAsCommand,
		}
	}

	return []string{
		"up",
	}
}

func runDockerComposeFromFilePath(p *string) {
	b, err := ioutil.ReadFile(*p)

	if err != nil {
		fmt.Println("Error opening file: " + *p)
	}

	var d DockerComposerData
	yaml.Unmarshal(b, &d)
	o := d.CommandifyOptions

	additionalArgs := getAdditionalArgs(o)
	cmd := &exec.Cmd{
		Path: "/usr/bin/env",
		Args: append(
			[]string{
				"/usr/bin/env",
				"docker",
				"compose",
				"-f",
				"/dev/stdin",
			},
			additionalArgs...,
		),
		Stdout: os.Stdout,
		Stdin:  bytes.NewReader(b),
		Stderr: os.Stderr,
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}

	cmd.Wait()
}

func main() {
	if len(os.Args) == 1 {
		printUseAndExit()
	}

	args := os.Args
	scriptPath := args[1]

	runDockerComposeFromFilePath(&scriptPath)
}
