package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

const modulePath = "go.pantheon.tech/stonework"

func init() {
	formatter := &logrus.TextFormatter{
		EnvironmentOverrideColors: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			call := strings.TrimPrefix(frame.Function, modulePath)
			function = fmt.Sprintf("%s()", strings.TrimPrefix(call, "/"))
			_, file = filepath.Split(frame.File)
			file = fmt.Sprintf("%s:%d", file, frame.Line)
			return color.Debug.Sprint(function), color.Secondary.Sprint(file)
		},
	}
	logrus.SetFormatter(formatter)
	logrus.AddHook(&callerHook{})
}

type callerHook struct {
}

func (c *callerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (c *callerHook) Fire(entry *logrus.Entry) error {
	if fn, ok := entry.Data[logrus.FieldKeyFunc]; ok && entry.Caller != nil {
		fmt.Printf("LOG ENTRY (fn: %v): %+v\n", fn, entry)
	}
	return nil
}
