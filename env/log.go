package env

/**
 * 日志级别
 * Error	错误、警告
 * Info		非错误和警告，比如：信息、调试
 */
const (
	LogLvInfo  = "Info"
	LogLvError = "Error"
)

/**
 * 日志驱动（日志输出的最终落地点）
 * Console	控制台
 * File		本地文件
 * Server	服务器
 */
const (
	LogDriverConsole = "console"
	LogDriverFile    = "file"
	LogDriverServer  = "server"
)
