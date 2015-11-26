package main

import (
	"errors"
)

type CommandHandler struct {
	cmdHandler map[string]func(args []string) error
}

func NewCommandHander() *CommandHandler {
	return &CommandHandler{
		cmdHandler: make(map[string](func(args []string) error)),
	}
}

func (ch *CommandHandler) SetHandler(cmd string, handler func(args []string) error) {
	ch.cmdHandler[cmd] = handler
}

func (ch *CommandHandler) GetHandler(cmd string) (func(args []string) error, error) {
	handler, exist := ch.cmdHandler[cmd]
	if !exist {
		return nil, errors.New("command not found: " + cmd)
	}
	return handler, nil
}
