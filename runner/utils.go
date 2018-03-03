package runner

import (
	"os"
	"path/filepath"
	"strings"
)

func initFolders() {
	runnerLog("InitFolders")
	path := setting.TmpPath()
	runnerLog("mkdir %s", path)
	err := os.Mkdir(path, 0755)
	if err != nil {
		runnerLog(err.Error())
	}
}

func isTmpDir(path string) bool {
	absolutePath, _ := filepath.Abs(path)
	absoluteTmpPath, _ := filepath.Abs(setting.TmpPath())

	return absolutePath == absoluteTmpPath
}

func isIgnoredFolder(path string) bool {
	paths := strings.Split(path, "/")
	if len(paths) <= 0 {
		return false
	}

	for _, e := range strings.Split(setting.Ignored(), ",") {
		if strings.Contains(path, strings.TrimSpace(e)) {
			return true
		}
	}
	return false
}

func isWatchedFile(path string) bool {
	if strings.HasSuffix(path, "_test.go") {
		return false
	}

	absolutePath, _ := filepath.Abs(path)
	absoluteTmpPath, _ := filepath.Abs(setting.TmpPath())

	if strings.HasPrefix(absolutePath, absoluteTmpPath) {
		return false
	}

	ext := filepath.Ext(path)

	for _, e := range strings.Split(setting.ValidExt(), ",") {
		if strings.TrimSpace(e) == ext {
			return true
		}
	}

	return false
}

func shouldRebuild(eventName string) bool {
	for _, e := range strings.Split(setting.NoRebuildExt(), ",") {
		e = strings.TrimSpace(e)
		fileName := strings.Replace(strings.Split(eventName, ":")[0], `"`, "", -1)
		if strings.HasSuffix(fileName, e) {
			return false
		}
	}

	return true
}

func createBuildErrorsLog(message string) bool {
	logFilePath := setting.BuildErrorsFilePath()
	var f *os.File
	if _, err := os.Stat(logFilePath); err != nil {
		f, err = os.Create(logFilePath)
		if err != nil {
			return false
		}
		defer f.Close()
	} else {
		f, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return false
		}
		defer f.Close()
	}

	_, err := f.WriteString(message)
	if err != nil {
		return false
	}

	return true
}

func removeBuildErrorsLog() error {
	logFilePath := setting.BuildErrorsFilePath()
	if _, err := os.Stat(logFilePath); err != nil {
		return nil
	}
	return os.Remove(logFilePath)
}
