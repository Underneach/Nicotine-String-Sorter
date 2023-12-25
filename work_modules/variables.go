package work_modules

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/text/encoding"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"
)

var (
	// ColorBlue Цвета
	ColorBlue    = color.New(color.FgBlue).Add(color.Bold)
	ColorGreen   = color.New(color.FgGreen).Add(color.Bold)
	ColorRed     = color.New(color.FgRed).Add(color.Bold)
	ColorMagenta = color.New(color.FgMagenta).Add(color.Bold)
	ColorYellow  = color.New(color.FgYellow).Add(color.Bold)

	badSymbolsPattern, _          = regexp.Compile(`[^a-zA-Z0-9]+`) // Разрешенные для файла символы
	checkedLines         int64    = 0                               // Колво отработанных строк
	invalidLines         int64    = 0                               // Колво невалидных строк, 64 байта что бы блять наверняка
	checkedFiles                  = 0                               // Колво отработанных файлов
	currentFileSize      int64    = 0                               // Размер текущего файла в сорте
	currentFileLines     int64    = 0                               // Размер текущего файла в сорте
	line                 string                                     // Текущая строка для обработки
	request              string                                     // Текущий запрос для обработки
	appDir               string                                     //  Папка запуска
	readLines            []string                                   // Текущие строки для обработки
	resultFilesList      []string                                   // Список созданных файлов

	invalidPattern, _ = regexp.Compile(`.{201,}|UNKNOWN`) // Паттерн невалид строк
	requestStructMap  = make(map[string]*Work)            // Карта со структурой для каждого запроса
	fileDecoder       *encoding.Decoder                   // Декодер файла
	workerPool        *ants.MultiPool                     // Пул сортера
	sorterWG          sync.WaitGroup                      // Синхронизатор очка пула сортера
	writerPool        *ants.Pool                          // Пул записи
	writerWG          sync.WaitGroup                      // Синхронизатор очка пула записи
	dublesPool        *ants.Pool                          // Пул удаления дублей
	dublesWG          sync.WaitGroup                      // Синхронизатор очка пула дублей
	// Кеш получения доступного пула строк
	cacheMutex     sync.Mutex // Мютекс кеша метода получения колва  доступных строк
	cachedStrCount int        // Колво доступных строк | кешируется
	lastUpdate     time.Time  // Время с последней обновы cachedStrCount
	// Инициализируем аргументы сортера для глобального доступа
	filePathList   []string
	searchRequests []string
	saveType       string
)

type Work struct {
	requestPattern *regexp.Regexp      // Регулярка запроса
	resultStrings  map[string][]string // Лист с Url, Log и Pass
}

// InitVar Инициализируем аргументы сортера для глобального доступа
func InitVar(_filePathList []string, _searchRequests []string, _saveType string) {
	filePathList = _filePathList
	searchRequests = _searchRequests
	saveType = _saveType

	var err error
	workerPool, err = ants.NewMultiPool(min(runtime.NumCPU()-2, 1), 1000, ants.LeastTasks, ants.WithPreAlloc(true)) // Мультипул горутин для сорта 

	if err != nil {
		PrintErr()
		_, _ = ColorRed.Print("Невозможно запустить сортер : Ошибка пула горутин : \n\n\n", err, "\n\n\n   Нажмите Enter для выхода")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}

	writerPool, err = ants.NewPool(len(searchRequests))
	if err != nil {
		PrintErr()
		_, _ = ColorRed.Print("Невозможно запустить сортер : Ошибка пула записи : \n\n\n", err, "\n\n\n   Нажмите Enter для выхода")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}

	appDirRaw, err := os.Executable()
	if err != nil {
		PrintErr()
		fmt.Print("Не удалось определить папку запуска\n")
	}
	appDir = filepath.Dir(appDirRaw)
}
