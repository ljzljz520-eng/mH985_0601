package command

import (
	"errors"
	"strconv"
	"strings"
)

type Request struct {
	Name   string
	Params map[string]string
}

func Parse(line string) (Request, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return Request{}, errors.New("command is empty")
	}
	req := Request{Name: parts[0], Params: make(map[string]string)}
	for _, part := range parts[1:] {
		pair := strings.SplitN(part, "=", 2)
		if len(pair) != 2 || pair[0] == "" {
			return Request{}, errors.New("command parameter must be key=value")
		}
		req.Params[pair[0]] = pair[1]
	}
	return req, nil
}

func (r Request) String(key string) string {
	return r.Params[key]
}

func (r Request) Int(key string, fallback int) int {
	value := r.Params[key]
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func (r Request) Has(key string) bool {
	_, ok := r.Params[key]
	return ok
}
