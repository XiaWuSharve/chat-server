package zlog

import "fmt"

func Fatal(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func Error(err error) {
	if err == nil {
		return
	}
}

func Info(pattern string, value ...any) {
	// TODO: 换成zap
	fmt.Printf(pattern, value...)
}
