package user_modules

import (
	"bufio"
	"github.com/fatih/color"
	"os"
	"runtime"
	"strings"
	"sync"
)

var (

	// ColorBlue Цвета
	ColorBlue    = color.New(color.FgBlue).Add(color.Bold)
	ColorGreen   = color.New(color.FgGreen).Add(color.Bold)
	ColorRed     = color.New(color.FgRed).Add(color.Bold)
	ColorMagenta = color.New(color.FgMagenta).Add(color.Bold)
	ColorYellow  = color.New(color.FgYellow).Add(color.Bold)

	filesSize       int64                       // Размер всех входных файлов
	userInputReader = bufio.NewReader(os.Stdin) // Альтернативный ридер инпута с поддержкой пробелов
	userOs          = runtime.GOOS              // ОС юзера
	updateWG        sync.WaitGroup              // ВГ обновы
	isLogoPrinted   = false                     // Напечатано ли лого

	workMode       string
	filePathList   []string
	searchRequests []string
	saveType       string
	cleanType      string
	delimetr       = strings.TrimSpace(":")
)
