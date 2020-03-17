/*
 * Copyright (c) 2020. BlizzTrack
 */

package network

import "testing"

func TestParseCommand(t *testing.T) {
	command, _ := ParseCommand("v1/summary")

	if command.Product != "" && command.File != "" || command.Method != "summary"{
		t.Error("failed to parse version")
	}
}

func TestNewCommand(t *testing.T) {
	command := NewCommand("summary", "", "")

	if command.String() != "v1/summary" {
		t.Error("failed to stringify command got " + command.String())
	}

	command = NewCommand("products", "pro", "versions")

	if command.String() != "v1/products/pro/versions" {
		t.Error("failed to stringify command got " + command.String())
	}
}
