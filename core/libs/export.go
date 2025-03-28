package libs

import (
	"runtime"
	"strings"

	"github.com/yicaoyimuys/GoGameServer/core/libs/common"
	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"github.com/yicaoyimuys/GoGameServer/core/libs/stack"
	"github.com/yicaoyimuys/GoGameServer/core/libs/system"
	"go.uber.org/zap"
)

func init() {
}

// 获取调用者的文件名和行号
func getCallerInfo() (string, int) {
	// 获取调用栈
	pc, file, line, ok := runtime.Caller(2) // 使用2来跳过当前函数和包装函数
	if !ok {
		return "unknown", 0
	}

	// 获取函数名
	funcName := runtime.FuncForPC(pc).Name()
	// 提取包名
	parts := strings.Split(funcName, ".")
	pkgPath := strings.Join(parts[:len(parts)-1], ".")

	// 简化文件路径，只保留最后的包名和文件名
	fileParts := strings.Split(file, "/")
	fileName := fileParts[len(fileParts)-1]

	// 如果包名中包含GoGameServer，只保留之后的部分
	if idx := strings.Index(pkgPath, "GoGameServer"); idx != -1 {
		pkgPath = pkgPath[idx+len("GoGameServer")+1:]
	}

	return pkgPath + "/" + fileName, line
}

func ERR(msg string, fields ...zap.Field) {
	file, line := getCallerInfo()
	logger.Error(msg, append(fields, zap.String("caller", file), zap.Int("line", line))...)
}

func WARN(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func INFO(msg string, fields ...zap.Field) {
	file, line := getCallerInfo()
	logger.Info(msg, append(fields, zap.String("caller", file), zap.Int("line", line))...)
}

func DEBUG(msg string, fields ...zap.Field) {
	file, line := getCallerInfo()
	logger.Debug(msg, append(fields, zap.String("caller", file), zap.Int("line", line))...)
}

func Run() {
	system.Run()
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	return common.If(condition, trueVal, falseVal)
}

func CheckError(err error) {
	stack.CheckError(err)
}
