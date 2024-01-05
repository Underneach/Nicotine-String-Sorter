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
	ColorBlue        = color.New(color.FgBlue).Add(color.Bold)
	ColorGreen       = color.New(color.FgGreen).Add(color.Bold)
	ColorRed         = color.New(color.FgRed).Add(color.Bold)
	ColorMagenta     = color.New(color.FgMagenta).Add(color.Bold)
	ColorYellow      = color.New(color.FgYellow).Add(color.Bold)
	ColorYellowLight = color.New(color.FgYellow)

	isFileInProcessing   bool                                                                  // Обрабатывается ли файл
	isResultWrited       bool                                                                  // Записан ли файл
	badSymbolsPattern, _                                     = regexp.Compile(`[^a-zA-Z0-9]+`) // Разрешенные для файла символы
	checkedLines         int64                               = 0                               // Колво отработанных строк
	invalidLines         int64                               = 0                               // Колво невалидных строк, 64 байта что бы блять наверняка
	currFileMatchLines   int64                               = 0                               //
	matchLines           int64                               = 0                               // Кол во подошедших строк
	checkedFiles                                             = 0                               // Колво отработанных файлов
	currentFileSize      int64                               = 0                               // Размер текущего файла в сорте
	currentFileLines     int64                               = 0                               // Размер текущего файла в сорте
	currFileInvalidLines int64                               = 0                               // Кол во невалидных строк текущего файла
	currFileCheckedLines                                     = 0                               // Прочеканные строки в текущем файле
	TMPlinesLen                                              = 0                               // Чанк строк в файле
	reqLen                                                   = 0                               // Кол во запросов
	currPath             string                                                                // Текущий файл
	invalidPattern, _    = regexp.Compile(`.{201,}|UNKNOWN`)                                   // Паттерн невалид строк
	requestStructMap     = make(map[string]*Work)                                              // Карта со структурой для каждого запроса
	fileDecoder          *encoding.Decoder                                                     // Декодер файла
	workerPool           *ants.MultiPoolWithFunc                                               // Пул сортера
	sorterWG             sync.WaitGroup                                                        // Синхронизатор очка пула сортера
	writerPool           *ants.PoolWithFunc                                                    // Пул записи
	writerWG             sync.WaitGroup                                                        // Синхронизатор очка пула записи
	dublesPool           *ants.PoolWithFunc                                                    // Пул удаления дублей
	dublesWG             sync.WaitGroup                                                        // Синхронизатор очка пула дублей
	cacheMutex           sync.Mutex                                                            // Мютекс кеша метода получения колва  доступных строк
	cachedStrCount       int                                                                   // Колво доступных строк | кешируется
	lastUpdate           time.Time                                                             // Время с последней обновы cachedStrCount
	pBar                 *progressbar.ProgressBar                                              // Прогресс бар
	runDir               = GetRunDir()                                                         // Папка запуска
	fileChannelMap       = make(map[string]chan [2]string)                                     // Мапа каналов
	dublesMutex          sync.Mutex                                                            // Мютекс вывода результата дублей
	RSMMutex             sync.RWMutex                                                          // Мютекс карты со структурой для каждого запроса
	filePathList         []string
	searchRequests       []string
	saveType             string
)

type Work struct {
	requestPattern *regexp.Regexp // Регулярка запроса
	resultStrings  []string       // Найденые строки
	resultFile     string         // Название файла с найдеными строками
}

func InitVar(_filePathList []string, _searchRequests []string, _saveType string) {
	filePathList = _filePathList
	searchRequests = _searchRequests
	saveType = _saveType
	reqLen = len(searchRequests)

	// Инициализируем каналы
	for _, path := range filePathList {
		fileChannelMap[path] = make(chan [2]string)
	}

	/*

		А сейчас самое ахуительное, sizePerPool в размере что 10000, что 1000000 не даёт сильного различия (1-3 секунды в тестах файла на 4,3кк строк),
		но при преаллоке жрет на 150мб больше. Преклоака тема интересная, сразу кушает память, что бы в дальнейшем не тратить на это время,
		что по сравнению с пулом неорграниченного размера даёт прирост в 20-30%

		БИМ БИМ БАМ БАМ блять

	*/

	var poolerr error
	workerPool, poolerr = ants.NewMultiPoolWithFunc(
		runtime.NumCPU(),
		100000,
		func(line interface{}) { ProcessLine(line.(string)) },
		ants.RoundRobin,
		ants.WithPreAlloc(true),
	)

	if poolerr != nil {
		PrintErr()
		ColorRed.Print("Невозможно запустить сортер : Ошибка пула сортера : \n\n\n		", poolerr, "\n\n\n   Нажмите Enter для выхода")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}

	writerPool, poolerr = ants.NewPoolWithFunc(
		len(searchRequests),
		func(request interface{}) { Writer(request.(string)) },
		ants.WithPreAlloc(true),
	)

	if poolerr != nil {
		PrintErr()
		ColorRed.Print("Невозможно запустить сортер : Ошибка пула записи : \n\n\n		", poolerr, "\n\n\n   Нажмите Enter для выхода")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}

	dublesPool, poolerr = ants.NewPoolWithFunc(
		len(searchRequests),
		func(request interface{}) { DublesRemove(request.(string)) },
		ants.WithPreAlloc(true))

	if poolerr != nil {
		PrintErr()
		ColorRed.Print("Невозможно запустить сортер : Ошибка пула удаления дублей : \n\n\n		", poolerr, "\n\n\n   Нажмите Enter для выхода")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}
}
