package work_modules

import (
	"bufio"
	"fmt"
	"golang.org/x/text/transform"
	"os"
	"strings"
)

func Sorter(path string) {

	err := GetCurrentFileSize(path)
	if err != nil {
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

	bar := CreateBar()

	scanner := bufio.NewScanner(transform.NewReader(file, fileDecoder))
	for scanner.Scan() {
		readLines = append(readLines, strings.TrimSpace(scanner.Text()))
	}

	linesLen := len(readLines) // Кол во строк в слайсе
	if linesLen == 0 {
		PrintErr()
		fmt.Print(path, " : Не удалось прочитать строки в файле\n")
		return
	}
	sorterWG.Add(linesLen)          //	 	
	ProcessLines(readLines)         // ОТПРАВЛЯЕМ СТРОКИ В ПУЛ И ЖДЁМ
	sorterWG.Wait()                 // 
	checkedLines += int64(linesLen) // Прибавляем строки
	clear(readLines)                // Чистим строки
	_ = bar.Finish()                // Завершаем бар
	_ = bar.Exit()                  // Закрываем бар
	file.Close()                    // Закрываем файл
	workerPool.Reboot()             // Ребутим пул
	WriteResult()                   // Пишем результат в файл
	clear(readLines)                // Чистим переменные
	for _, request := range searchRequests {
		clear(requestStructMap[request].resultStrings) // Чистим переменные
	}
	PrintFileSorted(path) // Пишем файл отсортрован
	checkedFiles++        // Прибавляем пройденные файлы
}

func ProcessLines(readLines []string) {
	for _, line = range readLines {
		err := workerPool.Submit(func() {
			Worker(line)
		})
		if err != nil {
			PrintErr()
			fmt.Print("Ошибка отправки строки в пул : ", err, "\n")
		}
	}

	i := requestStructMap[request] // ПОХУЙ
	for _, a := range tempResultLines {
		i.resultStrings = append(i.resultStrings, a)
	}
	requestStructMap[request] = i
}

func Worker(line string) {
	defer sorterWG.Done()
	for _, request := range searchRequests {
		if invalidPattern.MatchString(line) {
			invalidLines++
			return
		} else {
			result = requestStructMap[request].requestPattern.FindString(line)
			if result != "" {
				tempResultLines = append(tempResultLines, result)
				return
			} else {
				return
			}
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
		err := writerPool.Submit(func() {
			Writer(request)
		})

		if err != nil {
			PrintErr()
			fmt.Printf("%s : Ошибка записи найденных строк : %s\n", request, err)
			PrintInfo()
			fmt.Print("Запустите сортер с правами Админимтратора, если ошибка связана с доступом\n")
			continue
		}
		writerWG.Wait()
	}
}

/*

ЗАПИСЬ В ФАЙЛ

*/

func Writer(request string) {
	defer writerWG.Done()

	if len(requestStructMap[request].resultStrings) == 0 {
		return
	}

	resultFileName := appDir + "/" + badSymbolsPattern.ReplaceAllString(request, "_") + ".txt"

	resultFile, err := os.OpenFile(resultFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		PrintErr()
		fmt.Printf("%s : Ошибка записи найденных строк : %s\n", request, err)
		PrintInfo()
		fmt.Print("Запустите сортер с правами Админимтратора, если ошибка связана с доступом\n")
		return
	}

	_, err = bufio.NewWriter(resultFile).WriteString(strings.Join(requestStructMap[request].resultStrings, "\n"))

	if err != nil {
		PrintErr()
		fmt.Printf("%s : Ошибка записи найденных строк : %s\n", request, err)
		return
	}

	resultFilesList = append(resultFilesList, resultFileName)

	_ = resultFile.Close()
	PrintSuccess()
	ColorBlue.Print(request)
	fmt.Print(" : Записано ")
	ColorBlue.Print(len(requestStructMap[request].resultStrings))
	fmt.Print(" уникальных строк\n")

}
