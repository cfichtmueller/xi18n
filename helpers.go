package main

import (
	"fmt"
	"os"
)

func createComment(key string) string {
	return fmt.Sprintf("/* %s */\n", key)
}

func createEntry(key string, content string) string {
	return fmt.Sprintf("\"%s\" = \"%s\";\n", key, content)
}

func write(f *os.File, content string) {
	_, err := f.Write([]byte(content))
	if err != nil {
		fail("error: %v", err)
	}
}

func warn(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "WARN: "+format+"\n", args...)
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
