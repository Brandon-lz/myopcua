package utils

import (
    "fmt"
    "log/slog"
    "os"
    "runtime/debug"
)

// recover and log， when create goroutine panic, it will be recovered and log
// 推荐在windows下使用，linux可以直接查看控制台输出查找问题
func RecoverAndLog(callbacks ...func(err error)) {
    if r := recover(); r != nil {
        panicFile, err := os.OpenFile("panic.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            slog.Error("写入panic日志文件失败：%v", err)
            return
        }
        defer panicFile.Close()

        switch e:=r.(type) {
        case error:
            panicFile.WriteString(fmt.Sprintf("未知异常：%v\n%s\n", r, string(debug.Stack())))
            for _, callback := range callbacks {
                callback(e)
            }
        default:
            panicFile.WriteString(fmt.Sprintf("未知panic：%v\n%s\n", r, string(debug.Stack())))
            for _, callback := range callbacks {
                callback(fmt.Errorf("%v", r))
            }
        }
    }
}

// usage
// defer RecoverAndLog()
// defer RecoverAndLog(func(err error){quit <- true})
