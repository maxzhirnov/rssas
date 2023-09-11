package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	Log     *logrus.Logger
	LogFile *os.File
}

func NewLogger() (*Logger, error) {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logFile, err := setupFile()
	if err != nil {
		return nil, err
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	return &Logger{
		Log:     log,
		LogFile: logFile,
	}, nil
}

func setupFile() (*os.File, error) {
	// Проверяем наличие папки logs и создаем, если ее нет
	logsDir := "logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		err := os.Mkdir(logsDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	// Создаем или открываем файл для логирования
	logFile, err := os.OpenFile(filepath.Join(logsDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open or create log file: %v\\n", err)
		return nil, err
	}
	return logFile, nil
}
