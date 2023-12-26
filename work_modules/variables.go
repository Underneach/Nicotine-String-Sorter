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

	badSymbolsPattern, _                                     = regexp.Compile(`[^a-zA-Z0-9]+`) // Разрешенные для файла символы
	checkedLines         int64                               = 0                               // Колво отработанных строк
	invalidLines         int64                               = 0                               // Колво невалидных строк, 64 байта что бы блять наверняка
	checkedFiles                                             = 0                               // Колво отработанных файлов
	currentFileSize      int64                               = 0                               // Размер текущего файла в сорте
	currentFileLines     int64                               = 0                               // Размер текущего файла в сорте
	line                 string                                                                // Текущая строка для обработки
	request              string                                                                // Текущий запрос для обработки
	appDir               string                                                                //  Папка запуска
	result               string                                                                // Результат регулярки
	readLines            []string                                                              // Текущие строки для обработки
	tempResultLines      []string                                                              // Временный слайс найденных строк
	resultFilesList      []string                                                              // Список созданных файлов
	invalidPattern, _    = regexp.Compile(`.{201,}|UNKNOWN`)                                   // Паттерн невалид строк
	requestStructMap     = make(map[string]*Work)                                              // Карта со структурой для каждого запроса
	fileDecoder          *encoding.Decoder                                                     // Декодер файла
	workerPool           *ants.MultiPool                                                       // Пул сортера
	sorterWG             sync.WaitGroup                                                        // Синхронизатор очка пула сортера
	writerPool           *ants.Pool                                                            // Пул записи
	writerWG             sync.WaitGroup                                                        // Синхронизатор очка пула записи
	dublesPool           *ants.Pool                                                            // Пул удаления дублей
	dublesWG             sync.WaitGroup                                                        // Синхронизатор очка пула дублей
	cacheMutex           sync.Mutex                                                            // Мютекс кеша метода получения колва  доступных строк
	cachedStrCount       int                                                                   // Колво доступных строк | кешируется
	lastUpdate           time.Time                                                             // Время с последней обновы cachedStrCount
	filePathList         []string
	searchRequests       []string
	saveType             string
)

type Work struct {
	requestPattern *regexp.Regexp // Регулярка запроса
	resultStrings  []string       // Найденые строки
}

// InitVar Инициализируем аргументы сортера для глобального доступа
func InitVar(_filePathList []string, _searchRequests []string, _saveType string) {
	filePathList = _filePathList
	searchRequests = _searchRequests
	saveType = _saveType

	var err error
	workerPool, err = ants.NewMultiPool(runtime.NumCPU(), 1000, ants.RoundRobin, ants.WithPreAlloc(true)) // Мультипул горутин для сорта 

	if err != nil {
		PrintErr()
		ColorRed.Print("Невозможно запустить сортер : Ошибка пула сортера : \n\n\n		", err, "\n\n\n   Нажмите Enter для выхода")
		fmt.Scanln()
		os.Exit(1)
	}

	writerPool, err = ants.NewPool(len(searchRequests))
	if err != nil {
		PrintErr()
		ColorRed.Print("Невозможно запустить сортер : Ошибка пула записи : \n\n\n		", err, "\n\n\n   Нажмите Enter для выхода")
		fmt.Scanln()
		os.Exit(1)
	}

	dublesPool, err = ants.NewPool(len(searchRequests))
	if err != nil {
		PrintErr()
		ColorRed.Print("Невозможно запустить сортер : Ошибка пула удаления дублей : \n\n\n		", err, "\n\n\n   Нажмите Enter для выхода")
		fmt.Scanln()
		os.Exit(1)
	}

	appDirRaw, adrerr := os.Executable()
	if adrerr != nil {
		PrintErr()
		fmt.Print("Не удалось определить папку запуска : Сохранение в ")
		ColorBlue.Print("Папку пользователя\n")
		var hmerr error
		appDir, hmerr = os.UserHomeDir()
		if hmerr != nil {
			fmt.Println("Я не буду обрабатывать ошибки из за того что ты запустил софт на какой то ебанине, жми Enter, установи последнюю винду и не еби себе мозги")
			fmt.Scanln()
			os.Exit(1)
		}

	} else {
		appDir = filepath.Dir(appDirRaw)
	}
}
