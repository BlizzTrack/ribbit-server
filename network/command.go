/*
 * Copyright (c) 2020. BlizzTrack
 */

package network

import (
	"errors"
	"fmt"
	"strings"
)

type Command struct {
	Method  string
	Product string
	File    string
}

func NewCommand(method, product, file string) Command {
	return Command{
		Method:  method,
		Product: product,
		File:    file,
	}
}

func ParseCommand(command string) (Command, error) {
	commands := strings.Split(command, "/")

	version, options := commands[0], commands[1:]

	if version != "v1" {
		return Command{}, errors.New("invalid protocol version, only v1 allowed")
	}

	if len(options) == 1 {
		return Command{
			Method:  options[0],
			Product: "",
			File:    "",
		}, nil
	}

	if len(options) == 3 {
		return Command{
			Method:  options[0],
			Product: options[1],
			File:    options[2],
		}, nil
	}

	return Command{}, errors.New("failed to create command")
}

func (c Command) String() string {
	if c.File == "" && c.Product =="" {
		return fmt.Sprintf("v1/%s", c.Method)
	}

	return fmt.Sprintf("v1/%s/%s/%s", c.Method, c.Product, c.File)
}
