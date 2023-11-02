package main

import (
	"fmt"
	"log/slog"
)

type ControlLog struct {
	Group     string `vici:"group"`
	Level     string `vici:"level"`
	IkeSAName string `vici:"ikesa-name"`
	Message   string `vici:"msg"`
}

func EmitControlLog(cl ControlLog) {
	slog.Info(cl.Message, "group", cl.Group, "level", cl.Level)
}

func Errorf(format string, args ...any) {
	slog.Default().Error(fmt.Sprintf(format, args...))
}
