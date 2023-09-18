package log

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/mizhidajiangyou/go-linux/cmd"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	file   *os.File
)

const (
	infoLogFile    = "log/info.log"
	errorLogFile   = "log/error.log"
	commonLogFile  = "log/run.log"
	maxLogFileSize = 104857600
	defaultLevel   = zapcore.InfoLevel
)

func init() {

	// 创建日志文件
	file, _ = openFile(commonLogFile)

	// 创建日志记录器
	logger = makeLog(file, defaultLevel)

}

func openFile(file string) (*os.File, error) {
	e := cmd.Touch(file)
	if e != nil {
		panic(fmt.Sprintf("打开日志文件失败：%s\n", e))
	}
	f, e := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		panic(fmt.Sprintf("打开日志文件失败：%s\n", e))
	}
	ff, e := f.Stat()
	if e != nil {
		panic(fmt.Sprintf("打开日志文件失败：%s\n", e))
	}
	if ff.Size() > maxLogFileSize {
		e = compressLogFile(file)
		if e != nil {
			panic(fmt.Sprintf("压缩日志文件失败：%s\n", e))
		}
		f.Close()
		f, e = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if e != nil {
			panic(fmt.Sprintf("打开日志文件失败：%s\n", e))
		}
	}

	return f, e

}

func makeLog(file *os.File, l zapcore.Level) *zap.Logger {
	// 创建 Encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 创建信息日志 Core
	//Level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl >= zap.NewProductionConfig().Level.Level()
	//})
	Level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= l
	})
	Core := zapcore.NewCore(encoder, zapcore.AddSync(file), Level)

	// 创建控制台和文件输出的 MultiWriteSyncer
	consoleOutput := zapcore.Lock(os.Stdout)
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(file), consoleOutput)

	// 创建 Tee Core，将日志同时输出到文件和控制台
	Core = zapcore.NewCore(encoder, multiWriteSyncer, Level)

	return zap.New(Core, zap.ErrorOutput(multiWriteSyncer))
}

// Debugf 输出 Debug 级别的日志
func Debugf(format string, a ...any) {
	logger.Debug(fmt.Sprintf(format, a...))
}

// Infof 输出 Info 级别的日志
func Infof(format string, a ...any) {
	logger.Info(fmt.Sprintf(format, a...))
}

// Errorf 输出 Error 级别的日志
func Errorf(format string, a ...any) {
	logger.Error(fmt.Sprintf(format, a...))
}

// Warnf 输出 Error 级别的日志
func Warnf(format string, a ...any) {
	logger.Warn(fmt.Sprintf(format, a...))
}

func Fatalf(format string, a ...any) {
	logger.Fatal(fmt.Sprintf(format, a...))

}

func Fatal(err error) {
	logger.Fatal(fmt.Sprintf("%s", err))
}

// 压缩日志文件
func compressLogFile(logFile string) error {
	// 打开日志文件
	file, err := os.Open(logFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取当前日期作为压缩文件名的一部分
	dateStr := time.Now().Format("2006-01-02-15-04-05")
	compressedFileName := logFile + "." + dateStr + ".tar.gz"

	// 创建压缩文件
	compressedFile, err := os.Create(compressedFileName)
	if err != nil {
		return err
	}
	defer compressedFile.Close()

	// 创建 gzip.Writer
	gzipWriter := gzip.NewWriter(compressedFile)
	defer gzipWriter.Close()

	// 创建 tar.Writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// 读取日志文件并写入 tar.Writer
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header := &tar.Header{
		Name:    filepath.Base(logFile),
		Size:    info.Size(),
		Mode:    int64(info.Mode()),
		ModTime: info.ModTime(),
	}
	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return err
	}

	// 删除原日志文件
	err = os.Remove(logFile)
	if err != nil {
		return err
	}

	return nil
}

func SetLogLevel(lev string) {
	var logLevel zapcore.Level
	switch lev {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	logger = makeLog(file, logLevel)

}
