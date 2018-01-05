package tools

import (
	log "github.com/cihub/seelog"
)

/**
 * [Init_log 日志初始化配置]
 */
func InitLog() {
	logger, err := log.LoggerFromConfigAsFile("./config/log.xml")

	if err != nil {
		log.Critical("err parsing config log file", err)
		return
	}

	log.ReplaceLogger(logger)
}

/**
 * [LogInfo 记录info级别日志]
 */
func LogInfo(v ...interface{}) {
	log.Info(v)
}

/**
 * [LogError 记录error级别日志]
 */
func LogError(v ...interface{}) {
	log.Error(v)
}

/**
 * [LogFlush 写入日志]
 * @return {[type]} [将缓存写入日志里]
 */
func LogFlush(){
	defer log.Flush()
}
