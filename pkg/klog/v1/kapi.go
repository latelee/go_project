/*
按大小切换：可指定日志文件大小，达到时自动转存一下文件，文件以数字后缀递增
按日期切割：当0点时，即切换，新目录日志文件数字从1开始

注：调试中可能会改日期，如先用10号测试，生成日志后，再切换到9号，时间流逝到10号，此时，从10号目录找最新的文件，写之
*/
package klog

import (
	"fmt"
	"os"
	"path/filepath"
)

func (l *loggingT) println(s severity, args ...interface{}) {
	// 用等级限制，下同
	if s < l.LogLevel {
		return
	}

	buf, file, line := l.header(s, 0)
	fmt.Fprintln(buf, args...)
	l.output(s, buf, file, line, false)
}

func (l *loggingT) print(s severity, args ...interface{}) {
	l.printDepth(s, 1, args...)
}

func (l *loggingT) printDepth(s severity, depth int, args ...interface{}) {
	if s < l.LogLevel {
		return
	}
	buf, file, line := l.header(s, depth)
	fmt.Fprint(buf, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf, file, line, false)
}

func (l *loggingT) printf(s severity, format string, args ...interface{}) {
	if s < l.LogLevel {
		return
	}
	buf, file, line := l.header(s, 0)
	fmt.Fprintf(buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf, file, line, false)
}

// printWithFileLine behaves like print but uses the provided file and line number.  If
// alsoLogToStderr is true, the log message always appears on standard error; it
// will also appear in the log file unless --logtostderr is set.
func (l *loggingT) printWithFileLine(s severity, file string, line int, alsoToStderr bool, args ...interface{}) {
	if s < l.LogLevel {
		return
	}
	buf := l.formatHeader(s, file, line)
	fmt.Fprint(buf, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf, file, line, alsoToStderr)
}

///////// 对外接口
func DebugL(args ...interface{}) {
	logging.print(infoDebug0, args...)
}

func DebugLln(args ...interface{}) {
	logging.println(infoDebug0, args...)
}

func DebugLf(format string, args ...interface{}) {
	logging.printf(infoDebug0, format, args...)
}

func DebugM(args ...interface{}) {
	logging.print(infoDebug1, args...)
}

func DebugMln(args ...interface{}) {
	logging.println(infoDebug1, args...)
}

func DebugMf(format string, args ...interface{}) {
	logging.printf(infoDebug1, format, args...)
}

func DebugH(args ...interface{}) {
	logging.print(infoDebug2, args...)
}

func DebugHln(args ...interface{}) {
	logging.println(infoDebug2, args...)
}

func DebugHf(format string, args ...interface{}) {
	logging.printf(infoDebug2, format, args...)
}

// Info logs to the INFO level.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Info(args ...interface{}) {
	logging.print(infoLog, args...)
}

// Infoln logs to the INFO level.
// Arguments are handled in the manner of fmt.Println; a newline is always appended.
func Infoln(args ...interface{}) {
	logging.println(infoLog, args...)
}

// Infof logs to the INFO level.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Infof(format string, args ...interface{}) {
	logging.printf(infoLog, format, args...)
}

// Warn logs to the WARN level
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Warn(args ...interface{}) {
	logging.print(warningLog, args...)
}

// Warnln logs to the WARN level
// Arguments are handled in the manner of fmt.Println; a newline is always appended.
func Warnln(args ...interface{}) {
	logging.println(warningLog, args...)
}

// Warnf logs to the WARN level
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Warnf(format string, args ...interface{}) {
	logging.printf(warningLog, format, args...)
}

// Error logs to the ERROR level
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Error(args ...interface{}) {
	logging.print(errorLog, args...)
}

// Errorln logs to the ERROR level
// Arguments are handled in the manner of fmt.Println; a newline is always appended.
func Errorln(args ...interface{}) {
	logging.println(errorLog, args...)
}

// Errorf logs to the ERROR level
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Errorf(format string, args ...interface{}) {
	logging.printf(errorLog, format, args...)
}

// Fatal logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Fatal(args ...interface{}) {
	logging.print(fatalLog, args...)
}

// Fatalln logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Println; a newline is always appended.
func Fatalln(args ...interface{}) {
	logging.println(fatalLog, args...)
}

// Fatalf logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Fatalf(format string, args ...interface{}) {
	logging.printf(fatalLog, format, args...)
}

// 用Info等级
func Print(args ...interface{}) {
	logging.print(infoLog, args...)
}

func Println(args ...interface{}) {
	logging.println(infoLog, args...)
}

func Printf(format string, args ...interface{}) {
	logging.printf(infoLog, format, args...)
}

// 强制刷新日志
func FlushLog() {
	logging.flushAll()
}

// 暂时不用的，先注释掉
// // 用Info等级
// func (l *LoggingT) Print(args ...interface{}) {
// 	l.print(infoLog, args...)
// }

// func (l *LoggingT) Println(args ...interface{}) {
// 	l.println(infoLog, args...)
// }

// func (l *LoggingT) Printf(format string, args ...interface{}) {
// 	l.printf(infoLog, format, args...)
// }

/*
封装初始化函数
dirname - 目录
prefix - 日志文件前缀
level - 等级数值（0开始，依次是DEBUG0、DEBUG1、DEBUG2、INFO、WARN、ERROR）
logsize - 日志大小
showtype - 输出方式 1：只到终端，2：只到文件 3：终端+文件
loginterval - 日志写文件间隔 为0即时写

针对全局的初始化，即整个工程只有一个日志实例
*/
func Init_normal(dirname, prefix string, level int, logsize int, showtype int, loginterval int) {
	logging = &loggingT{}

	logging.LogLevel = severity(level)

	if (showtype & 0x01) == 1 {
		logging.toStderr = true
		logging.alsoToStderr = true
	}
	if (showtype & 0x02) == 2 {
		logging.toStderr = false
	}

	if isFile(dirname) {
		// flag.Set("logFile", dirname)
		logging.logFile = dirname
	} else {
		LogDir := "./log"
		if err := mkDir(dirname); err == nil {
			LogDir = dirname
		} else {
			mkDir(LogDir)
		}
		logging.logDir = dirname
	}

	// 如果是有效的值，则设置，否则用默认的
	if logsize > 0 {
		logging.LogFileMaxSize = uint64(logsize)
	}
	if prefix != "" {
		logging.LogNamePrefix = prefix+"."
	}
	logging.LogHeadType = 1
	logging.LogFlushInterval = loginterval

	logging.hasWritten = false
	if logging.LogFlushInterval > 0 {
		go logging.flushDaemon()
	}

	return
}

// 对外提供的初始化结构体的函数

/*
初始化函数
t - 结构体

TODO 目标是创建不同的日志实例，在同一工程可以有多个实例分别记录不同模块的日志。
目前未完成。
*/
func NewKlog(t LogParam_t) *LoggingT {
	foo := LoggingT{}

	foo.LogLevel = t.LogLevel

	if isFile(t.LogDirFile) {
		foo.logFile = t.LogDirFile
	} else {
		LogDir := "./log"
		if err := mkDir(t.LogDirFile); err == nil {
			LogDir = t.LogDirFile
		} else {
			mkDir(LogDir)
		}
		foo.logDir = LogDir

	}
	if (t.Showtype & 0x01) == 1 {
		foo.toStderr = true
		foo.alsoToStderr = true
	}
	if (t.Showtype & 0x02) == 2 {
		foo.toStderr = false
	}

	// 如果是有效的值，则设置，否则用默认的
	foo.LogFileMaxSize = t.LogFileMaxSize
	if foo.LogFileMaxSize == 0 {
		foo.LogFileMaxSize = MaxSize
	}
	foo.LogNamePrefix = t.LogNamePrefix
	if foo.LogNamePrefix == "" {
		foo.LogNamePrefix = filepath.Base(os.Args[0])+"."
	}
	if t.LogHeadType == 0 {
		foo.LogHeadType = 1
	}
	foo.LogFlushInterval = t.LogFlushInterval

	return &foo

}
