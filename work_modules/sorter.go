package work_modules

import (
	"bufio"
	"fmt"
	"github.com/zeebo/xxh3"
	"golang.org/x/text/encoding/unicode"
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
			compiledRegEx, err = regexp.Compile("(.*" + regexp.QuoteMeta(request) + ".*:.+:.+)")
		}

		if err != nil {
			PrintErr()
			fmt.Printf("%s : Ошибка компиляции запроса : %s\n", request, err)
			continue
		}

		currentStruct := new(Work)
		currentStruct.requestPattern = compiledRegEx
		currentStruct.resultFile = runDir + fileBadSymbolsPattern.ReplaceAllString(request, "_") + ".txt"
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
	sorterStringHashMap = make(map[uint64]bool)
	isFileInProcessing = false
	isResultWrited = false
	TMPlinesLen = 0
	currFileDubles = 0
	for _, req := range searchRequests {
		sorterRequestStatMapCurrFile[req] = 0
	}

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
	pBar = CreatePBar()
	go PBarUpdater()
	go SorterWriteResult()

	scanner := bufio.NewScanner(transform.NewReader(file, fileDecoder))

	for ; scanner.Scan(); TMPlinesLen++ {
		workWG.Add(1)
		_ = sorterPool.Invoke(scanner.Text())
	}

	workWG.Wait()

	checkedLines += int64(TMPlinesLen) // Прибавляем строки
	_ = pBar.Finish()                  // Завершаем бар
	_ = pBar.Exit()                    // Закрываем бар
	close(sorterWriteChannelMap[currPath])

	isFileInProcessing = false
	for !isResultWrited {
		time.Sleep(time.Millisecond * 100)
	}

	file.Close() // Закрываем файл

	sorterWriteChannelMap[currPath] = nil // Чистим канал
	PrintSortInfo()
	PrintFileSorted(path)            // Пишем файл отсортрован
	checkedFiles++                   // Прибавляем пройденные файлы
	matchLines += currFileMatchLines // Суммируем найденые строки
	sorterDubles += currFileDubles
}

func SorterProcessLine(line string) {
	defer workWG.Done()
	for _, request := range searchRequests {
		if result := requestStructMap[request].requestPattern.FindStringSubmatch(line); len(result) == 2 {
			sorterWriteChannelMap[currPath] <- [2]string{request, result[1]}
			return
		}
	}
}

func SorterWriteResult() {

	for _, request := range searchRequests {
		if resultFile, err := os.OpenFile(requestStructMap[request].resultFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			sorterResultFileMap[request] = resultFile
			sorterResultWriterMap[request] = bufio.NewWriter(transform.NewWriter(resultFile, unicode.UTF8.NewDecoder()))
		} else {
			PrintResultWriteErr(request, err)
		}
	}

	for {
		if data, ok := <-sorterWriteChannelMap[currPath]; ok {
			hash := xxh3.HashString(data[1])
			if _, ok := sorterStringHashMap[hash]; !ok {
				sorterStringHashMap[hash] = true
				_, _ = sorterResultWriterMap[data[0]].WriteString(data[1] + "\n")
				sorterRequestStatMapCurrFile[data[0]]++
			} else {
				currFileDubles++
			}
			continue
		} else {
			break
		}
	}

	for _, request := range searchRequests {
		sorterResultFileMap[request].Close()
		currFileMatchLines += sorterRequestStatMapCurrFile[request]
		sorterRequestStatMap[request] += sorterRequestStatMapCurrFile[request]
	}

	isResultWrited = true // сообщаем о том, что файл записан
}
