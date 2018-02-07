package runner

import (
	"fmt"
	gobuild "go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const envSettingsPrefix = "RUNNER_"

type settingKey string

func (s settingKey) String() string {
	return string(s)
}

const (
	rootPath        settingKey = "root_path"
	tmpPath                    = "tmp_path"
	buildPath                  = "build_path"
	buildLogName               = "build_log_name"
	validExt                   = "valid_ext"
	noRebuildExt               = "no_rebuild_ext"
	ignored                    = "ignored"
	buildDelay                 = "build_delay"
	colors                     = "colors"
	logColorMain               = "log_color_main"
	logColorBuild              = "log_color_build"
	logColorRunner             = "log_color_runner"
	logColorWatcher            = "log_color_watcher"
	logColorApp                = "log_color_app"
)

type Setting map[settingKey]string

func (s Setting) RootPath() string     { return filepath.Join(gobuild.Default.GOPATH, "src", s[rootPath]) }
func (s Setting) BuildPath() string    { return filepath.Join(s.RootPath(), s[buildPath]) }
func (s Setting) TmpPath() string      { return filepath.Join(s.RootPath(), s[tmpPath]) }
func (s Setting) BuildLogName() string { return s[buildLogName] }
func (s Setting) ValidExt() string     { return s[validExt] }
func (s Setting) NoRebuildExt() string { return s[noRebuildExt] }
func (s Setting) Ignored() string      { return s[ignored] }
func (s Setting) BuildDelay() time.Duration {
	value, _ := strconv.Atoi(s[buildDelay])
	return time.Duration(value) * time.Millisecond
}
func (s Setting) Colors() string          { return s[colors] }
func (s Setting) LogColorMain() string    { return s[logColorMain] }
func (s Setting) LogColorBuild() string   { return s[logColorBuild] }
func (s Setting) LogColorRunner() string  { return s[logColorRunner] }
func (s Setting) LogColorWatcher() string { return s[logColorWatcher] }
func (s Setting) LogColorApp() string     { return s[logColorApp] }
func (s Setting) BuildErrorsFilePath() string {
	return filepath.Join(s.TmpPath(), s.BuildLogName())
}
func (s Setting) LogColor(role string) string {
	switch role {
	case "main":
		return s.LogColorMain()
	case "build":
		return s.LogColorBuild()
	case "runner":
		return s.LogColorRunner()
	case "watcher":
		return s.LogColorWatcher()
	case "app":
		return s.LogColorApp()
	}
	return ""
}
func (s Setting) Set(key settingKey, value string) { s[key] = value }

var defaultSetting = Setting{
	rootPath:        ".",
	buildPath:       ".",
	tmpPath:         "./tmp",
	buildLogName:    "runner-build-errors.log",
	validExt:        ".go, .tpl, .tmpl, .html",
	noRebuildExt:    ".tpl, .tmpl, .html",
	ignored:         "assets, tmp",
	buildDelay:      "600",
	colors:          "1",
	logColorMain:    "cyan",
	logColorBuild:   "yellow",
	logColorRunner:  "green",
	logColorWatcher: "magenta",
	logColorApp:     "",
}

var setting = defaultSetting

type LogColors map[logColorsKey]string

type logColorsKey string

const (
	reset         logColorsKey = "reset"
	black                      = "black"
	red                        = "red"
	green                      = "green"
	yellow                     = "yellow"
	blue                       = "blue"
	magenta                    = "magenta"
	cyan                       = "cyan"
	white                      = "white"
	boldBlack                  = "bold_black"
	boldRed                    = "bold_red"
	boldGreen                  = "bold_green"
	boldYellow                 = "bold_yellow"
	boldBlue                   = "bold_blue"
	boldMagenta                = "bold_magenta"
	boldCyan                   = "bold_cyan"
	boldWhite                  = "bold_white"
	brightBlack                = "bright_black"
	brightRed                  = "bright_red"
	brightGreen                = "bright_green"
	brightYellow               = "bright_yellow"
	brightBlue                 = "bright_blue"
	brightMagenta              = "bright_magenta"
	brightCyan                 = "bright_cyan"
	brightWhite                = "bright_white"
)

func (l LogColors) Reset() string            { return l[reset] }
func (l LogColors) Black() string            { return l[black] }
func (l LogColors) Red() string              { return l[red] }
func (l LogColors) Green() string            { return l[green] }
func (l LogColors) Yellow() string           { return l[yellow] }
func (l LogColors) Blue() string             { return l[blue] }
func (l LogColors) Magenta() string          { return l[magenta] }
func (l LogColors) Cyan() string             { return l[cyan] }
func (l LogColors) White() string            { return l[white] }
func (l LogColors) BoldBlack() string        { return l[boldBlack] }
func (l LogColors) BoldRed() string          { return l[boldRed] }
func (l LogColors) BoldGreen() string        { return l[boldGreen] }
func (l LogColors) BoldYellow() string       { return l[boldYellow] }
func (l LogColors) BoldBlue() string         { return l[boldBlue] }
func (l LogColors) BoldMagenta() string      { return l[boldMagenta] }
func (l LogColors) BoldCyan() string         { return l[boldCyan] }
func (l LogColors) BoldWhite() string        { return l[boldWhite] }
func (l LogColors) BrightBlack() string      { return l[brightBlack] }
func (l LogColors) BrightRed() string        { return l[brightRed] }
func (l LogColors) BrightGreen() string      { return l[brightGreen] }
func (l LogColors) BrightYellow() string     { return l[brightYellow] }
func (l LogColors) BrightBlue() string       { return l[brightBlue] }
func (l LogColors) BrightMagenta() string    { return l[brightMagenta] }
func (l LogColors) BrightCyan() string       { return l[brightCyan] }
func (l LogColors) BrightWhite() string      { return l[brightWhite] }
func (l LogColors) Select(key string) string { return l[logColorsKey(key)] }

var logColors = LogColors{
	reset:         "0",
	black:         "30",
	red:           "31",
	green:         "32",
	yellow:        "33",
	blue:          "34",
	magenta:       "35",
	cyan:          "36",
	white:         "37",
	boldBlack:     "30;1",
	boldRed:       "31;1",
	boldGreen:     "32;1",
	boldYellow:    "33;1",
	boldBlue:      "34;1",
	boldMagenta:   "35;1",
	boldCyan:      "36;1",
	boldWhite:     "37;1",
	brightBlack:   "30;2",
	brightRed:     "31;2",
	brightGreen:   "32;2",
	brightYellow:  "33;2",
	brightBlue:    "34;2",
	brightMagenta: "35;2",
	brightCyan:    "36;2",
	brightWhite:   "37;2",
}

func logColor(logName string) string {
	return logColors.Select(setting.LogColor(logName))
}

func loadEnvSettings() {
	for key, _ := range setting {
		envKey := envSettingsPrefix + strings.ToUpper(key.String())
		if value := os.Getenv(envKey); value != "" {
			setting.Set(settingKey(key), value)
		}
	}
}

func loadRunnerConfigSettings(configPath string) {
	if _, err := os.Stat(configPath); err != nil {
		return
	}

	logger.Printf("Loading settings from %s", configPath)
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Printf("Loading settings from %s, error=%s", configPath, err)
		return
	}
	if err := yaml.Unmarshal(b, &setting); err != nil {
		logger.Printf("Loading settings from %s, error=%s", configPath, err)
	}
}

func initSettings(configPath string) {
	if configPath != "" {
		if _, err := os.Stat(configPath); err != nil {
			fmt.Printf("Can't find config file `%s`\n", configPath)
			os.Exit(1)
		} else {
			os.Setenv("RUNNER_CONFIG_PATH", configPath)
		}
	}

	loadEnvSettings()
	loadRunnerConfigSettings(configPath)
}
