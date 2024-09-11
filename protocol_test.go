package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCommand(t *testing.T) {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	r := bytes.NewReader(cmd.Bytes())

	pcmd := ParseCommand(r)
	assert.Equal(t, cmd, pcmd)
}
