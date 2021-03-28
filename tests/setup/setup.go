package setup

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	valuesFilesDir = "values"
)

var (
	chdirOnce sync.Once
	chdirErr  error

	testsDir string
	repoRoot string
)

func Setup() error {
	flag.Parse()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	if err := chdirToRepositoryRoot(); err != nil {
		return err
	}
	return nil
}

func chdirToRepositoryRoot() error {
	chdirOnce.Do(func() {
		var err error
		testsDir, err = os.Getwd()
		if err != nil {
			chdirErr = err
		}
		repoRoot = filepath.Dir(filepath.Dir(filepath.Dir(testsDir)))
		if err = os.Chdir(repoRoot); err != nil {
			chdirErr = fmt.Errorf("Could not change working directory from %s to repository root at %s: %s", testsDir, repoRoot, err)
		}
	})
	return chdirErr
}
