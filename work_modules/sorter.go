package work_modules

import (
	"bufio"
	"fmt"
	"golang.org/x/text/transform"
	"math"
	"os"
	"regexp"
	"time"
)

func RunSorter() {

	PrintInfo()
	fmt.Print("Запуск сортера...")

	var (
		compiledRegEx *regexp.Regexp
		err           error
	)

	for _, request := range searchRequests {

		switch saveType {
		case "1":
			compiledRegEx, err = regexp.Compile(".*" + regexp.QuoteMeta(request) + ".*:(.+:.+)")
		case "2":
			compiledRegEx, err = regexp.Compile("(" + ".*" + regexp.QuoteMeta(request) + ".*:.+:.+)")
		}

		if err != nil {
			PrintErr()
			fmt.Printf("%s : Ошибка компиляции запроса : %s\n", request, err)
			continue
		}

		currentStruct := new(Work)
		currentStruct.requestPattern = compiledRegEx
		currentStruct.resultFile = runDir + `\` + fileBadSymbolsPattern.ReplaceAllString(request, "_") + ".txt"
		requestStructMap[request] = currentStruct
	}

	if len(requestStructMap) == 0 {
		PrintZeroRequestsErr()
	}

	fmt.Print("\r")
	PrintSuccess()
	fmt.Print("Сортер запущен   \n\n")
}

func Sorter(path string) {

	currPath = path
	sorterStringChannelMap[currPath] = make(chan string)
	isFileInProcessing = false
	isResultWrited = false
	TMPlinesLen = 0

	if err := GetCurrentFileSize(path); err != nil {
		PrintFileReadErr(path, err)
		return
	}

	PrintFileInfo(path)
	PrintLinesChunk()
	fileDecoder = GetEncodingDecoder(path)

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)

	if err != nil {
		PrintFileReadErr(path, err)
		return
	}

	if GetAviableStringsCount() > currentFileLines {
		sorterPool.Tune(int(math.Round(float64(currentFileLines) / 3)))
	} else {
		sorterPool.Tune(int(math.Round(float64(GetAviableStringsCount()) / 3)))
	}

	isFileInProcessing = true
	go PBarUpdater()
	go SorterProcessResult()
	go SorterProcessInputLines()

	scanner := bufio.NewScanner(transform.NewReader(file, fileDecoder))

	for ; scanner.Scan(); TMPlinesLen++ {
		workWG.Add(1)
		sorterStringChannelMap[currPath] <- scanner.Text()
	}

	workWG.Wait()
	close(sorterStringChannelMap[currPath])

	checkedLines += int64(TMPlinesLen) // Прибавляем строки
	_ = pBar.Finish()                  // Завершаем бар
	_ = pBar.Exit()                    // Закрываем бар
	close(sorterResultChannelMap[currPath])

	isFileInProcessing = false
	for !isResultWrited {
		time.Sleep(time.Millisecond * 100)
	}

	file.Close() // Закрываем файл

	for _, request := range searchRequests {
		requestStructMap[request].resultStrings = nil // Чистим список
	}

	sorterResultChannelMap[currPath] = nil // Чистим канал
	PrintFileSorted(path)                  // Пишем файл отсортрован
	checkedFiles++                         // Прибавляем пройденные файлы
	matchLines += currFileMatchLines       // Суммируем найденые строки
}

func SorterProcessInputLines() {
	for {
		if data, ok := <-sorterStringChannelMap[currPath]; !ok {
			break
		} else {
			_ = sorterPool.Invoke(data)
			continue
		}
	}
}

/*

Обрабатываем строку

*/

func SorterProcessLine(line string) {
	defer workWG.Done()
	for _, request := range searchRequests {
		if result := requestStructMap[request].requestPattern.FindStringSubmatch(line); len(result) == 2 {
			sorterResultChannelMap[currPath] <- [2]string{request, result[1]}
			return
		}
	}
}

func SorterProcessResult() {

	ResultListMap := make(map[string][]string)

	for _, request := range searchRequests {
		ResultListMap[request] = requestStructMap[request].resultStrings
	}

	for {
		if data, ok := <-sorterResultChannelMap[currPath]; ok {
			ResultListMap[data[0]] = append(ResultListMap[data[0]], data[1])
			continue
		} else {
			break
		}
	}

	for _, request := range searchRequests {
		currFileMatchLines = int64(len(ResultListMap[request]))
		requestStructMap[request].resultStrings = ResultListMap[request]
	}

	ResultListMap = nil   // чистим
	SorterWriteResult()   // Пишем результат в файл
	isResultWrited = true // сообщаем о том, что файл записан
}
