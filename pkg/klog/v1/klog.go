package klog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var logging *loggingT

// severity identifies the sort of log: info, warning etc. It also implements
// the flag.Value interface. The -LogLevel flag is of type severity and
// should be modified only through the flag.Value interface. The values match
// the corresponding constants in C++.
type severity int32 // sync/atomic int32

// These constants identify the log levels in order of increasing severity.
// A message written to a high-severity log file is also written to each
// lower-severity log file.
const (
	infoDebug0 severity = iota
	infoDebug1
	infoDebug2
	infoLog
	warningLog
	errorLog
	fatalLog
	numSeverity = 6
)

// 仅用于整齐显示等级字符串
var severityString = []string{
	infoDebug0: " DEBUG0 ",
	infoDebug1: " DEBUG1 ",
	infoDebug2: " DEBUG2 ",
	infoLog:    " INFO   ",
	warningLog: " WARN   ",
	errorLog:   " ERROR  ",
	fatalLog:   " FATAL  ",
}

// OutputStats tracks the number of output lines and bytes written.
type OutputStats struct {
	lines int64
	bytes int64
}

// Lines returns the number of lines written.
func (s *OutputStats) Lines() int64 {
	return atomic.LoadInt64(&s.lines)
}

// Bytes returns the number of bytes written.
func (s *OutputStats) Bytes() int64 {
	return atomic.LoadInt64(&s.bytes)
}

// Stats tracks the number of lines of output and number of bytes
// per severity level. Values must be read with atomic.LoadInt64.
var Stats struct {
	Info, Warning, Error OutputStats
}

var severityStats = [numSeverity]*OutputStats{
	infoLog:    &Stats.Info,
	warningLog: &Stats.Warning,
	errorLog:   &Stats.Error,
}

// flushSyncWriter is the interface satisfied by logging destinations.
type flushSyncWriter interface {
	Flush() error
	Sync() error
	io.Writer
}

// Flush flushes all pending log I/O.
func Flush() {
	logging.lockAndFlushAll()
}

/*
LogParam_t、LoggingT、loggingT 三者有联系又有区别的原因：
LogParam_t：仅有参数，外部传递
LoggingT：按实例封装输出函数
loggingT：内部使用
*/
type LogParam_t struct {
	// 日志等级
	LogLevel severity

	// 输出方式 1: std 2: file 3: std+file
	Showtype int

	// 指定文件或目录 TODO
	LogDirFile string

	// 日志文件前缀
	LogNamePrefix string

	// 单个日志文件大小，单位：字节
	LogFileMaxSize uint64

	// 输出日志每行的前缀，0：带文件行号  1：不带文件行号
	LogHeadType int

	// 写文件时，日志回写频率，单位：秒。为0表示即时刷新
	LogFlushInterval int
}

type LoggingT struct {
	loggingT
}

type loggingT struct {
	// freeList is a list of byte buffers, maintained under freeListMu.
	freeList *buffer
	// freeListMu maintains the free list. It is separate from the main mutex
	// so buffers can be grabbed and printed to without holding the main lock,
	// for better parallelization.
	freeListMu sync.Mutex

	// mu protects the remaining elements of this structure and is
	// used to synchronize logging.
	mu sync.Mutex
	// file holds writer for each of the log types.
	file flushSyncWriter

	// 当前日志日期，用于判断新日期目录
	logDate string

	// 输出到终端
	toStderr bool
	// 输出文件的同时是否输出到终端
	alsoToStderr bool
	// 指定文件（只写到该文件，优化级最高）
	logFile string

	// 日志目录
	logDir string

	// 设置定时写入文件时，首次写入，此值为true，用于在定时前也写到日志中
	hasWritten bool
	LogParam_t
}

// buffer holds a byte Buffer for reuse. The zero value is ready for use.
type buffer struct {
	bytes.Buffer
	tmp  [64]byte // temporary byte array for creating headers.
	next *buffer
}

// getBuffer returns a new, ready-to-use buffer.
func (l *loggingT) getBuffer() *buffer {
	l.freeListMu.Lock()
	b := l.freeList
	if b != nil {
		l.freeList = b.next
	}
	l.freeListMu.Unlock()
	if b == nil {
		b = new(buffer)
	} else {
		b.next = nil
		b.Reset()
	}
	return b
}

// putBuffer returns a buffer to the free list.
func (l *loggingT) putBuffer(b *buffer) {
	if b.Len() >= 256 {
		// Let big buffers die a natural death.
		return
	}
	l.freeListMu.Lock()
	b.next = l.freeList
	l.freeList = b
	l.freeListMu.Unlock()
}

/*
header formats a log header as defined by the C++ implementation.
It returns a buffer containing the formatted header and the user's file and line number.
The depth specifies how many stack frames above lives the source line to be identified in the log message.

Log lines have this form:
	Lmmdd hh:mm:ss.uuuuuu threadid file:line] msg...
where the fields are defined as follows:
	L                A single character, representing the log level (eg 'I' for INFO)
	mm               The month (zero padded; ie May is '05')
	dd               The day (zero padded)
	hh:mm:ss.uuuuuu  Time in hours, minutes and fractional seconds
	threadid         The space-padded thread ID as returned by GetTID()
	file             The file name
	line             The line number
	msg              The user-supplied message
*/
func (l *loggingT) header(s severity, depth int) (*buffer, string, int) {
	_, file, line, ok := runtime.Caller(3 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		if slash := strings.LastIndex(file, "/"); slash >= 0 {
			path := file
			file = path[slash+1:]
		}
	}
	return l.formatHeader(s, file, line), file, line
}

// formatHeader formats a log header using the provided file name and line number.
func (l *loggingT) formatHeader(s severity, file string, line int) *buffer {
	now := time.Now()
	if line < 0 {
		line = 0 // not a real line number, but acceptable to someDigits
	}
	if s > fatalLog {
		s = infoLog // for safety.
	}
	buf := l.getBuffer()

	// Avoid Fprintf, for speed. The format is so simple that we can do it quickly by hand.
	// It's worth about 3X. Fprintf is hard.
	if l.LogHeadType == 0 {
		// 格式示例：
		// [2022-10-23 11:28:15.130 rootCmd.go:90] xxx
		year, month, day := now.Date()
		hour, minute, second := now.Clock()
		buf.tmp[0] = '['
		buf.nDigits(4, 1, year, '0')
		buf.tmp[5] = '-'
		buf.twoDigits(6, int(month))
		buf.tmp[8] = '-'
		buf.twoDigits(9, day)
		buf.tmp[11] = ' '
		buf.twoDigits(12, hour)
		buf.tmp[14] = ':'
		buf.twoDigits(15, minute)
		buf.tmp[17] = ':'
		buf.twoDigits(18, second)
		buf.tmp[20] = '.'
		buf.threeDigits(21, now.Nanosecond()/1000/1000)
		buf.tmp[24] = ' '

		buf.Write(buf.tmp[:25])

		buf.WriteString(file)
		buf.tmp[0] = ':'
		n := buf.someDigits(1, line)
		buf.tmp[n+1] = ']'
		buf.tmp[n+2] = ' '
		buf.Write(buf.tmp[:n+3])

	} else if l.LogHeadType == 1 {
		// 格式示例：
		// 2022-10-24 10:12:26.790  INFO   -  xxx
		// 2022-10-24 10:12:26.790  DEBUG2 -  xxx

		year, month, day := now.Date()
		hour, minute, second := now.Clock()
		buf.nDigits(4, 0, year, '0')
		buf.tmp[4] = '-'
		buf.twoDigits(5, int(month))
		buf.tmp[7] = '-'
		buf.twoDigits(8, day)
		buf.tmp[10] = ' '
		buf.twoDigits(11, hour)
		buf.tmp[13] = ':'
		buf.twoDigits(14, minute)
		buf.tmp[16] = ':'
		buf.twoDigits(17, second)
		buf.tmp[19] = '.'
		buf.threeDigits(20, now.Nanosecond()/1000/1000)
		buf.tmp[23] = ' '
		buf.tmp[24] = ' '
		buf.Write(buf.tmp[:24])

		buf.WriteString(severityString[s])
		buf.tmp[0] = '-'
		buf.tmp[1] = ' '
		buf.tmp[2] = ' '
		buf.Write(buf.tmp[:3])

	}
	return buf
}

// Some custom tiny helper functions to print the log header efficiently.

const digits = "0123456789"

// twoDigits formats a zero-prefixed two-digit integer at buf.tmp[i].
func (buf *buffer) twoDigits(i, d int) {
	buf.tmp[i+1] = digits[d%10]
	d /= 10
	buf.tmp[i] = digits[d%10]
}

func (buf *buffer) threeDigits(i, d int) {
	buf.tmp[i+2] = digits[d%10]
	d /= 10
	buf.tmp[i+1] = digits[d%10]
	d /= 10
	buf.tmp[i] = digits[d%10]
}

// nDigits formats an n-digit integer at buf.tmp[i],
// padding with pad on the left.
// It assumes d >= 0.
func (buf *buffer) nDigits(n, i, d int, pad byte) {
	j := n - 1
	for ; j >= 0 && d > 0; j-- {
		buf.tmp[i+j] = digits[d%10]
		d /= 10
	}
	for ; j >= 0; j-- {
		buf.tmp[i+j] = pad
	}
}

// someDigits formats a zero-prefixed variable-width integer at buf.tmp[i].
func (buf *buffer) someDigits(i, d int) int {
	// Print into the top, then copy down. We know there's space for at least
	// a 10-digit number.
	j := len(buf.tmp)
	for {
		j--
		buf.tmp[j] = digits[d%10]
		d /= 10
		if d == 0 {
			break
		}
	}
	return copy(buf.tmp[i:], buf.tmp[j:])
}

// redirectBuffer is used to set an alternate destination for the logs
type redirectBuffer struct {
	w io.Writer
}

func (rb *redirectBuffer) Sync() error {
	return nil
}

func (rb *redirectBuffer) Flush() error {
	return nil
}

func (rb *redirectBuffer) Write(bytes []byte) (n int, err error) {
	return rb.w.Write(bytes)
}

// // SetOutput sets the output destination for all severities
// func SetOutput(w io.Writer) {
// 	logging.mu.Lock()
// 	defer logging.mu.Unlock()
// 	rb := &redirectBuffer{
// 		w: w,
// 	}
// 	logging.file = rb
// }

// // SetOutputBySeverity sets the output destination for specific severity
// func SetOutputBySeverity(name string, w io.Writer) {
// 	logging.mu.Lock()
// 	defer logging.mu.Unlock()

// 	rb := &redirectBuffer{
// 		w: w,
// 	}
// 	logging.file = rb
// }

var onceDaemon sync.Once

// output writes the data to the log files and releases the buffer.
func (l *loggingT) output(s severity, buf *buffer, file string, line int, alsoToStderr bool) {
	l.mu.Lock()
	data := buf.Bytes()
	if l.toStderr { // 一旦设置输出到终端，则不会再写文件，互斥
		os.Stderr.Write(data)
	} else {
		// 判断在写文件时是否也要输出到终端 最后一个是判断阈值，因为前面已经做了限制，此处注释
		if alsoToStderr || l.alsoToStderr /* || s >= l.LogLevel.get() */ {
			os.Stderr.Write(data)
		}

		// 不管按日期目录还是单独文件，都只使用一个file，可统一
		if l.file == nil {
			if err := l.createFiles(s); err != nil {
				os.Stderr.Write(data) // Make sure the message appears somewhere.
				l.exit(err)
			}
		}
		l.file.Write(data)
	}

	l.putBuffer(buf)

	if l.hasWritten == false {
		l.flushAll()
	}

	l.mu.Unlock()
	// 暂不明白用途，注释掉
	// if stats := severityStats[s]; stats != nil {
	// 	atomic.AddInt64(&stats.lines, 1)
	// 	atomic.AddInt64(&stats.bytes, int64(len(data)))
	// }
}

// timeoutFlush calls Flush and returns when it completes or after timeout
// elapses, whichever happens first.  This is needed because the hooks invoked
// by Flush may deadlock when klog.Fatal is called from a hook that holds
// a lock.
func timeoutFlush(timeout time.Duration) {
	done := make(chan bool, 1)
	go func() {
		Flush() // calls logging.lockAndFlushAll()
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		fmt.Fprintln(os.Stderr, "klog: Flush took longer than", timeout)
	}
}

// stacks is a wrapper for runtime.Stack that attempts to recover the data for all goroutines.
func stacks(all bool) []byte {
	// We don't know how big the traces are, so grow a few times if they don't fit. Start large, though.
	n := 10000
	if all {
		n = 100000
	}
	var trace []byte
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, all)
		if nbytes < len(trace) {
			return trace[:nbytes]
		}
		n *= 2
	}
	return trace
}

// logExitFunc provides a simple mechanism to override the default behavior
// of exiting on error. Used in testing and to guarantee we reach a required exit
// for fatal logs. Instead, exit could be a function rather than a method but that
// would make its use clumsier.
var logExitFunc func(error)

// exit is called if there is trouble creating or writing log files.
// It flushes the logs and exits the program; there's no point in hanging around.
// l.mu is held.
func (l *loggingT) exit(err error) {
	fmt.Fprintf(os.Stderr, "log: exiting because of error: %s\n", err)
	// If logExitFunc is set, we do that instead of exiting.
	if logExitFunc != nil {
		logExitFunc(err)
		return
	}
	l.flushAll()
	os.Exit(2)
}

// syncBuffer joins a bufio.Writer to its underlying file, providing access to the
// file's Sync method and providing a wrapper for the Write method that provides log
// file rotation. There are conflicting methods, so the file cannot be embedded.
// l.mu is held for all its methods.
type syncBuffer struct {
	logger *loggingT
	*bufio.Writer
	file *os.File
	// sev      severity
	nbytes   uint64 // The number of bytes written to this file
	maxbytes uint64 // The max number of bytes this syncBuffer.file can hold before cleaning up.
}

func (sb *syncBuffer) Sync() error {
	return sb.file.Sync()
}

/*
每次写日志时，都会判断时间，因为不清楚何时要新建日志文件
*/
func (sb *syncBuffer) Write(p []byte) (n int, err error) {
	now := time.Now()
	// 判断时间，如果到0点，立刻切换
	// 注：用日期作判断条件，因为是字符串比较，因此带上目录，方便传递到下一函数
	// sb.logger.logDate 在前面的函数已创建
	datestr := fmt.Sprintf("%v/%04d-%02d/%02d", sb.logger.logDir, now.Year(), now.Month(), now.Day())
	if datestr > sb.logger.logDate {
		// fmt.Println("-----!!!! rotate file true  may be new day", datestr, sb.logger.logDate)
		if err := sb.rotateFile(datestr, now, true); err != nil {
			sb.logger.exit(err)
		}
		sb.logger.logDate = datestr
	}
	// 如果超过大小，也切换
	if sb.nbytes+uint64(len(p)) >= sb.maxbytes {
		// fmt.Println("------------------- rotate file max...")
		if err := sb.rotateFile(datestr, now, false); err != nil {
			sb.logger.exit(err)
		}
	}
	n, err = sb.Writer.Write(p)
	sb.nbytes += uint64(n)
	if err != nil {
		sb.logger.exit(err)
	}
	return
}

var log_file_idx int = 1

// rotateFile closes the syncBuffer's file and starts a new one.
// The startup argument indicates whether this is the initial startup of klog.
// If startup is true, existing files are opened for appending instead of truncated.
func (sb *syncBuffer) rotateFile(dir string, now time.Time, startup bool) error {
	if sb.file != nil {
		sb.Flush()
		sb.file.Close()
	}

	// now := time.Now()
	// if dir == "" {
	// 	dir = fmt.Sprintf("%v/%04d-%02d/%02d", sb.logger.logDir, now.Year(), now.Month(), now.Day())

	// }
	// 创建日期目录
	mkDir(dir)

	// 不管何种方法，都会有此名称，故先设置
	prename := fmt.Sprintf("%slog.%04d-%02d-%02d",
		sb.logger.LogNamePrefix,
		now.Year(),
		now.Month(),
		now.Day())

	var err error
	tfile := ""

	// 启动时(新目录，或新启动)，读已有日志文件，得大小，当成已有的
	if startup == true {
		log_file_idx1 := 1
		files, err := getFileListByPrefix(dir, prename, true, 1)
		// 有文件且读到最新的
		if err == nil && len(files) == 1 { // 读到有文件，读其大小判断追加或新建
			tfile = files[0]
			// 获取后缀，但要去掉前面的'.'
			log_file_idx1, _ = strconv.Atoi(tfile[strings.LastIndex(tfile, ".")+1:])

			size, err := fileSize(tfile)
			// fmt.Printf("read file: [%v] [%v] [%v] [%v]\n", tfile, log_file_idx1, size, sb.logger.LogFileMaxSize)

			// 最新文件大小小于指定大小，用之，否则新建
			if err == nil && size < sb.logger.LogFileMaxSize {
				sb.nbytes = size
				log_file_idx = log_file_idx1
			} else {
				sb.nbytes = 0
				log_file_idx = log_file_idx1 + 1
			}
		} else { // 读不到文件，从头新建
			sb.nbytes = 0
			log_file_idx = 1
		}
	} else { // 为 false 表示是自身回滚，肯定是新文件
		sb.nbytes = 0
		log_file_idx++
	}

	// fmt.Println("roatefile... ", startup, tfile, log_file_idx, sb.nbytes)

	// 日志名称
	logname := fmt.Sprintf("%v.%d", prename, log_file_idx)

	// 日志文件完整路径
	alogname := filepath.Join(dir, logname)

	sb.file, _, err = create(now, startup, alogname)
	if err != nil {
		return err
	}

	sb.Writer = bufio.NewWriterSize(sb.file, bufferSize)

	return nil

}

// bufferSize sizes the buffer associated with each log file. It's large
// so that log records can accumulate without the logging thread blocking
// on disk I/O. The flushDaemon will block instead.
const bufferSize = 256 * 1024

// createFiles creates all the log files for severity from sev down to infoLog.
// l.mu is held.
func (l *loggingT) createFiles(sev severity) error {
	now := time.Now()
	dir := fmt.Sprintf("%v/%04d-%02d/%02d", l.logDir, now.Year(), now.Month(), now.Day())

	// 首次肯定会到此函数，所以会首先创建
	l.logDate = dir
	sb := &syncBuffer{
		logger:   l,
		maxbytes: l.LogFileMaxSize,
	}
	// fmt.Println("============ rotate file true")
	if err := sb.rotateFile(dir, now, true); err != nil {
		return err
	}
	l.file = sb

	return nil
}

// The startup argument indicates whether this is the initial startup of klog.
// If startup is true, existing files are opened for appending instead of truncated.
func create(t time.Time, startup bool, logname string) (f *os.File, filename string, err error) {
	mylogfile := logname // 默认用日期日志文件
	// 指定文件，则只使用该文件
	if logging.logFile != "" {
		mylogfile = logging.logFile
	}

	f, err = openOrCreate(mylogfile, startup)
	if err == nil {
		return f, mylogfile, nil
	}

	return nil, "", fmt.Errorf("log: cannot create log: %v", err)
}

// The startup argument indicates whether this is the initial startup of klog.
// If startup is true, existing files are opened for appending instead of truncated.
func openOrCreate(name string, startup bool) (*os.File, error) {
	if startup {
		// fmt.Println("need to open file ", name)
		f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		return f, err
	}
	// fmt.Println("need to create file ", name)
	f, err := os.Create(name)
	return f, err
}

///////////////////////////////////////////

// 默认刷新频率，可调用外部接口Flush手动写

// flushDaemon periodically flushes the log file buffers.
func (l *loggingT) flushDaemon() {
	for range time.NewTicker(time.Duration(l.LogFlushInterval) * time.Second).C { // flushInterval
		l.lockAndFlushAll()
		l.hasWritten = true
	}
}

// lockAndFlushAll is like flushAll but locks l.mu first.
func (l *loggingT) lockAndFlushAll() {
	l.mu.Lock()
	l.flushAll()
	l.mu.Unlock()
}

// flushAll flushes all the logs and attempts to "sync" their data to disk.
// l.mu is held.
func (l *loggingT) flushAll() {
	// Flush from fatal down, in case there's trouble flushing.
	if l.file != nil {
		l.file.Flush() // ignore error
		l.file.Sync()  // ignore error
	}
}

// // init sets up the defaults and runs flushDaemon.
// func init() {
// 	logging.LogLevel = infoLog // Default LogLevel is INFO.
// 	logging.logDir = ""
// 	logging.logFile = ""
// 	logging.LogFileMaxSize = MaxSize
// 	logging.toStderr = true
// 	logging.alsoToStderr = false                      // tmp...
// 	logging.LogNamePrefix = filepath.Base(os.Args[0]) // default:program name
// 	logging.LogHeadType = 1
// 	logging.logDate = ""
// 	logging.LogFlushInterval = 0
// }
