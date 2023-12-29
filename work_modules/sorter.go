package work_modules

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"os"
	"strings"
	"sync"
	"time"
)

func Sorter(path string) {
	isFileInProcessing = false
	isResultWrited = false
	var tmpLines []string
	TMPlinesLen = 0
	currFileCheckedLines = 0

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

	isFileInProcessing = true
	//go PBarUpdater()
	go ProcessResult()

	scanner := bufio.NewScanner(transform.NewReader(file, fileDecoder))

	for ; scanner.Scan(); TMPlinesLen++ {
		if TMPlinesLen >= GetAviableStringsCount() {
			sorterWG.Add(TMPlinesLen)
			SendLinesToPool(tmpLines)
			currFileCheckedLines += TMPlinesLen
			TMPlinesLen = 0
		} else {
			tmpLines = append(tmpLines, scanner.Text())
		}
	}

	if len(tmpLines) > 0 {
		sorterWG.Add(TMPlinesLen)
		SendLinesToPool(tmpLines)
		currFileCheckedLines += TMPlinesLen
		TMPlinesLen = 0
	}

	checkedLines += int64(currFileCheckedLines)

	isFileInProcessing = false
	for isResultWrited == false {
		time.Sleep(time.Millisecond * 100)
	}

	file.Close()        // Закрываем файл
	_ = PBar.Finish()   // Завершаем бар
	_ = PBar.Exit()     // Закрываем бар
	workerPool.Reboot() // Ребутим пул

	for _, request := range searchRequests {
		clear(requestStructMap[request].resultStrings) // Чистим список
	}

	PrintFileSorted(path) // Пишем файл отсортрован
	checkedFiles++        // Прибавляем пройденные файлы
}

func SendLinesToPool(lines []string) {
	for _, line := range lines {
		if err := workerPool.Submit(func() {
			ProcessLine(line)
		}); err != nil {
			PrintErr()
			fmt.Printf("%s : Ошибка отправки строки в пул : %s \n", line, err)
			continue
		}
	}
	sorterWG.Wait()
}

/*

Обрабатываем строку

*/

func ProcessLine(line string) {
	defer sorterWG.Done()

	if invalidPattern.MatchString(line) {
		invalidLines++
		return
	}

	for _, request := range searchRequests {

		result = requestStructMap[request].requestPattern.FindString(line)
		if result != "" {
			ResultChannel <- [2]string{request, result}
		}
	}
}

/*

ЗАПИСЬ ФАЙЛОВ В НЕСКОЛЬКО ПОТОКОВ

*/

func WriteResult() {

	fmt.Print("\n")
	PrintInfo()
	fmt.Print("Удаление дублей и запись в файл\n")

	for _, request := range searchRequests {

		writerWG.Add(1)

		if wrterr := writerPool.Submit(func() {
			Writer(request)
		}); wrterr != nil {
			PrintResultWriteErr(request, wrterr)
			continue
		}
	}

	writerWG.Wait()
	isResultWrited = true

}

/*

ЗАПИСЬ В ФАЙЛ

*/

func Writer(request string) {
	defer writerWG.Done()

	if len(requestStructMap[request].resultStrings) == 0 {
		PrintErr()
		ColorBlue.Print(request)
		fmt.Print(" : Нет строк для записи\n")
		return
	}

	resultFileName := badSymbolsPattern.ReplaceAllString(request, "_") + ".txt"

	resultFile, err := os.OpenFile(resultFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		PrintErr()
		fmt.Printf("%s : Ошибка записи найденных строк : %s\n", request, err)
		PrintInfo()
		fmt.Print("Запустите сортер с правами Админиcтратора, если ошибка связана с доступом\n")
		return
	}

	if _, err = bufio.NewWriter(transform.NewWriter(resultFile, unicode.UTF8.NewDecoder())).WriteString(strings.Join(requestStructMap[request].resultStrings, "\n")); err != nil {
		PrintErr()
		fmt.Printf("%s : Ошибка записи найденных строк : %s\n", request, err)
		return
	}

	requestStructMap[request].resultFile = resultFileName

	_ = resultFile.Close()
	PrintSuccess()
	ColorBlue.Print(request)
	fmt.Print(" : Записано ")
	ColorBlue.Print(len(requestStructMap[request].resultStrings))
	fmt.Print(" уникальных строк\n")

}

func PBarUpdater() {

	PBar = CreatePBar()

	for isFileInProcessing {
		if currFileCheckedLines > int(currentFileLines) {
			_ = PBar.Set64(currentFileLines)
		} else {
			_ = PBar.Set(currFileCheckedLines)
		}

		time.Sleep(time.Millisecond * 250)

	}
}

func ProcessResult() {

	var ResultListMap sync.Map

	for _, request := range searchRequests {
		ResultListMap.Store(request, requestStructMap[request].resultStrings)
	}

	for isFileInProcessing {

		data, ok := <-ResultChannel

		if !ok {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		resultWG.Add(1)
		if rterr := resultPool.Submit(func() {
			if curStrs, ok := ResultListMap.Load(data[0]); ok {
				ResultListMap.Store(data[0], append(curStrs.([]string), data[1]))
			} else {
				PrintErr()
				fmt.Println("Ошибка загрузки найденной строки : ResultListMap")
				resultWG.Done()
			}
		}); rterr != nil {
			PrintErr()
			fmt.Print("Ошибка распределения найденой строки : ", rterr, "\n")
			resultWG.Done()
		}
	}

	resultWG.Wait()

	for _, request := range searchRequests {
		if curStrs, ok := ResultListMap.Load(request); ok {
			requestStructMap[request].resultStrings = curStrs.([]string)
		} else {
			PrintErr()
			fmt.Println("Ошибка при загрузке данных из ResultListMap для запроса", request)
		}
	}

	WriteResult() // Пишем результат в файл
	isResultWrited = true
}
