package work_modules

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/panjf2000/ants/v2"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/text/encoding"
	"os"
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

	// Общие
	isFileInProcessing       bool                                                       // Обрабатывается ли файл
	isResultWrited           bool                                                       // Записан ли файл
	fileBadSymbolsPattern, _                          = regexp.Compile(`[^a-zA-Z0-9]+`) // Разрешенные для файла символы
	checkedLines             int64                    = 0                               // Колво отработанных строк
	checkedFiles                                      = 0                               // Колво отработанных файлов
	currentFileSize          int64                    = 0                               // Размер текущего файла в сорте
	currentFileLines         int64                    = 0                               // Размер текущего файла в сорте
	fileDecoder              *encoding.Decoder                                          // Декодер файла
	cacheMutex               sync.Mutex                                                 // Мютекс кеша метода получения колва  доступных строк
	cachedStrCount           int64                                                      // Колво доступных строк | кешируется
	lastUpdate               time.Time                                                  // Время с последней обновы cachedStrCount
	pBar                     *progressbar.ProgressBar                                   // Прогресс бар
	runDir                   = GetRunDir()                                              // Папка запуска
	currPath                 string                                                     // Текущий файл
	poolerr                  error                                                      // Ошибка создания пула
	workWG                   sync.WaitGroup                                             // Синхронизатор очка
	TMPlinesLen              = 0                                                        // Чанк строк в файле
	currPathCut              string                                                     // Текущий файл без полного пути

	// Сортер
	currFileMatchLines           int64                             = 0                      //
	matchLines                   int64                             = 0                      // Кол во подошедших строк
	reqLen                                                         = 0                      // Кол во запросов
	sorterDubles                 int64                             = 0                      // Кол во повторяющихся строк
	requestStructMap                                               = make(map[string]*Work) // Карта со структурой для каждого запроса
	sorterPool                   *ants.MultiPoolWithFunc                                    // Пул сортера
	sorterWriteChannelMap        = make(map[string]chan [2]string)                          // Мапа каналов
	sorterRequestStatMap         = make(map[string]int64)                                   // Колво найденных для каждого запроса
	sorterRequestStatMapCurrFile = make(map[string]int64)                                   // Колво найденных для каждого запроса
	sorterResultWriterMap        = make(map[string]*bufio.Writer)                           // Мапа врайтера для каждого запроса
	sorterResultFileMap          = make(map[string]*os.File)                                // Мапа файла для каждого запроса
	sorterStringHashMap          = make(map[uint64]bool)                                    // Мапа хешей строк

	// Клинер
	validPattern, _         = regexp.Compile(`^[a-zA-Z0-9\.\,\!\?\:\;\-\'\"\@\/\#\$\%\^\&\*\(\)\_\+\=\~\x60]{10,256}$`)                         // Паттерн валида
	unknownPattern, _       = regexp.Compile(`UNKNOWN`)                                                                                         // Содержание UNKNOWN
	cleanerOutputFilesMap   = make(map[string]string)                                                                                           // Мапа выходных файлов
	cleanerResultChannelMap = make(map[string]chan string)                                                                                      // Мапа валид строк
	cleanerWriteFile        *os.File                                                                                                            // Файл записи
	cleanerInvalidLen       int64                                                                                       = 0                     // Кол во невалид строк
	currFileInvalidLen      int64                                                                                       = 0                     // Кол во повторяющихся строк
	cleanerDublesLen        int64                                                                                       = 0                     // Колво повторяющихся строк
	currFileDubles          int64                                                                                       = 0                     // 
	cleanerWritedString     int64                                                                                       = 0                     // Кол во записанных строк
	currFileWritedString    int64                                                                                       = 0                     //
	cleanerStringHashMap                                                                                                = make(map[uint64]bool) // Мапа хешей строк

	// Арги
	filePathList   []string
	searchRequests []string
	saveType       string
	workMode       string
	cleanType      string
)

type Work struct {
	requestPattern *regexp.Regexp // Регулярка запроса
	resultFile     string         // Название файла с найдеными строками
}

func InitVar(_workMode string, _filePathList []string, _searchRequests []string, _saveType string, _cleanType string) {
	workMode = _workMode
	filePathList = _filePathList
	searchRequests = _searchRequests
	saveType = _saveType
	cleanType = _cleanType
}

func InitSorter() {

	reqLen = len(searchRequests)

	// Инициализируем каналы
	for _, path := range filePathList {
		sorterWriteChannelMap[path] = make(chan [2]string)
	}

	sorterPool, poolerr = ants.NewMultiPoolWithFunc(
		runtime.NumCPU(),
		100000,
		func(line interface{}) { SorterProcessLine(line.(string)) },
		ants.RoundRobin,
		ants.WithPreAlloc(true),
	)

	if poolerr != nil {
		PrintErr()
		ColorRed.Print("Невозможно запустить сортер : Ошибка пула сортера : \n\n\n		", poolerr, "\n\n\n   Нажмите Enter для выхода")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}
}
