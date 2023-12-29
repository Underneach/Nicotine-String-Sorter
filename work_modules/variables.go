package work_modules

import (
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

	isFileInProcessing   bool                                                                  // Обрабатывается ли файл
	isResultWrited       bool                                                                  // Записан ли файл
	badSymbolsPattern, _                                     = regexp.Compile(`[^a-zA-Z0-9]+`) // Разрешенные для файла символы
	checkedLines         int64                               = 0                               // Колво отработанных строк
	invalidLines         int64                               = 0                               // Колво невалидных строк, 64 байта что бы блять наверняка
	checkedFiles                                             = 0                               // Колво отработанных файлов
	currentFileSize      int64                               = 0                               // Размер текущего файла в сорте
	currentFileLines     int64                               = 0                               // Размер текущего файла в сорте
	TMPlinesLen                                              = 0                               // Чанк строк в файле
	currFileCheckedLines                                     = 0                               // Прочеканные строки в текущем файле
	request              string                                                                // Текущий запрос для обработки
	result               string                                                                // Результат регулярки
	invalidPattern, _    = regexp.Compile(`.{201,}|UNKNOWN`)                                   // Паттерн невалид строк
	requestStructMap     = make(map[string]*Work)                                              // Карта со структурой для каждого запроса
	fileDecoder          *encoding.Decoder                                                     // Декодер файла
	workerPool           *ants.MultiPool                                                       // Пул сортера
	sorterWG             sync.WaitGroup                                                        // Синхронизатор очка пула сортера
	writerPool           *ants.Pool                                                            // Пул записи
	writerWG             sync.WaitGroup                                                        // Синхронизатор очка пула записи
	dublesPool           *ants.Pool                                                            // Пул удаления дублей
	resultPool           *ants.Pool                                                            // Пул распределения найденых строк
	resultWG             sync.WaitGroup                                                        // Синхронизатор очка Пула распределения
	dublesWG             sync.WaitGroup                                                        // Синхронизатор очка пула дублей
	cacheMutex           sync.Mutex                                                            // Мютекс кеша метода получения колва  доступных строк
	cachedStrCount       int                                                                   // Колво доступных строк | кешируется
	lastUpdate           time.Time                                                             // Время с последней обновы cachedStrCount
	PBar                 *progressbar.ProgressBar                                              // Прогресс бар
	ResultChannel        = make(chan [2]string)                                                // Канал с найдеными строками
	filePathList         []string
	searchRequests       []string
	saveType             string
)

type Work struct {
	requestPattern  *regexp.Regexp // Регулярка запроса
	resultStrings   []string       // Найденые строки
	resultFile      string         // Название файла с найдеными строками
	resultFileExist bool           // Есть ли файл с результатом
}

// InitVar Инициализируем аргументы сортера для глобального доступа
func InitVar(_filePathList []string, _searchRequests []string, _saveType string) {
	filePathList = _filePathList
	searchRequests = _searchRequests
	saveType = _saveType

	var poolerr error
	workerPool, poolerr = ants.NewMultiPool(runtime.NumCPU(), 100000, ants.RoundRobin, ants.WithPreAlloc(true)) // Мультипул горутин для сорта 

	if poolerr != nil {
		PrintErr()
		ColorRed.Print("Невозможно запустить сортер : Ошибка пула сортера : \n\n\n		", poolerr, "\n\n\n   Нажмите Enter для выхода")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}

	resultPool, poolerr = ants.NewPool(1000)
	if poolerr != nil {
		PrintErr()
		ColorRed.Print("Невозможно запустить сортер : Ошибка пула записи : \n\n\n		", poolerr, "\n\n\n   Нажмите Enter для выхода")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}

	writerPool, poolerr = ants.NewPool(len(searchRequests))
	if poolerr != nil {
		PrintErr()
		ColorRed.Print("Невозможно запустить сортер : Ошибка пула записи : \n\n\n		", poolerr, "\n\n\n   Нажмите Enter для выхода")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}

	dublesPool, poolerr = ants.NewPool(len(searchRequests))
	if poolerr != nil {
		PrintErr()
		ColorRed.Print("Невозможно запустить сортер : Ошибка пула удаления дублей : \n\n\n		", poolerr, "\n\n\n   Нажмите Enter для выхода")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}
}
