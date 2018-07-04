package main

import (
	"path/filepath"
	"strings"
)

func starts(search string, line string) bool {
	if len(line) < len(search) {
		return false
	}
	return line[:len(search)] == search
}

func after(search string, line string) string {
	if len(line) < len(search) {
		return ""
	}
	line = line[len(search):]
	line = strings.Replace(line, "\"", "", 2)
	line = strings.Replace(line, ";", "", 2)
	return filepath.Clean(line)
}

func trueFile(callFrom string, importName string) string {
	importName = strings.Replace(importName, "\"", "", 2)
	importName = strings.Replace(importName, ";", "", 2)
	pos := strings.LastIndex(callFrom, "/")
	ret := callFrom[:pos+1] + importName
	return filepath.Clean(ret)
}
