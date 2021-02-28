package logger

import "go.uber.org/zap"

func Logger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func SugaredLogger() *zap.SugaredLogger {
	defer Logger().Sync()

	return Logger().Sugar()
}

func TruncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}
