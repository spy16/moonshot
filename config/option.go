package config

import (
	"strings"
)

type Option func(l *viperLoader) error

func WithEnv(prefix ...string) Option {
	return func(l *viperLoader) error {
		l.useEnv = true
		if len(prefix) > 0 {
			l.envPrefix = strings.TrimSpace(prefix[0])
		}
		return nil
	}
}

func WithName(name string) Option {
	return func(l *viperLoader) error {
		l.confName = strings.TrimSpace(name)
		return nil
	}
}

func WithFile(filePath string) Option {
	return func(l *viperLoader) error {
		l.confFile = filePath
		return nil
	}
}
