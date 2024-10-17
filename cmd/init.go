package cmd

import "github.com/xpfo-go/logs"

func initLogs() {
	conf := logs.GetLogConf()
	conf.FileName = "log"
	conf.MaxAge = 21
	conf.Level = "debug"
	logs.InitLogSetting(conf)
}
