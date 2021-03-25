package main

import (
	"os"
	"path/filepath"

	_ "github.com/rancher/charts/testing/register"
	"github.com/rancher/charts/testing/run"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) == 1 {
		logrus.Infof("Usage: go run %s <globs> ", filepath.Base(os.Args[0]))
		return
	}
	if err := run.RunTests(os.Args[1:]); err != nil {
		logrus.Error(err)
		return
	}
}
