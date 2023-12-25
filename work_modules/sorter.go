package work_modules

import (
	"bufio"
	"fmt"
	"golang.org/x/text/transform"
	"os"
	"syscall"
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

	scanner := bufio.NewScanner(transform.NewReader(file, fileDecoder))
	bar := CreateBar()

	for scanner.Scan() {
		if len(readLines) < GetAviableStringsCount() {
			readLines = append(readLines, scanner.Text())
			println(line)
		} else {
			break
		}
	}

	linesLen := len(readLines)
	fmt.Println(linesLen)

	sorterWG.Add(linesLen)          //	 	
	ProcessLines(readLines)         // ОТПРАВЛЯЕМ СТРОКИ В ПУЛ И ЖДЁМ
	sorterWG.Wait()                 // 
	checkedLines = +int64(linesLen) // Прибавляем строки и ждём
	clear(readLines)                // Чистим строки
	barerr := bar.Add(linesLen)     // Увеличиваем бар
	if barerr != nil {
		PrintErr()
		fmt.Print("Ошибка Прогресс-Бара : ", barerr)
		syscall.Exit(1)
	}

	_ = bar.Finish()                         // Завершаем бар
	_ = bar.Exit()                           // Закрываем бар
	file.Close()                             // Закрываем файл
	workerPool.Reboot()                      // Ребутим пул
	WriteResult()                            // Пишем результат в файл
	clear(readLines)                         // Чистим переменные
	for _, request := range searchRequests { // Чистим переменные
		clear(requestStructMap[request].resultStrings)
	}
	PrintFileSorted(path)
	checkedFiles++ // Прибавляем пройденные файлы
}

func ProcessLines(readLines []string) {
	for _, line = range readLines {
		err := workerPool.Submit(func() {
			defer sorterWG.Done()
			Worker(line)
		})
		if err != nil {
			PrintErr()
			fmt.Print("Ошибка отправки строки в пул : ", err, "\n")
		}
	}
}

func Worker(line string) {
	for _, request := range searchRequests {
		if invalidPattern.MatchString(line) {
			invalidLines++
			return
		} else {
			result := requestStructMap[request].requestPattern.FindStringSubmatch(line)
			if len(result) == 3 {
				requestStructMap[request].resultStrings[line] = result
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

	for _, request := range searchRequests {

		writerWG.Add(1)
		err := writerPool.Submit(func() {
			defer writerWG.Done()
			Writer(request)
		})

		if err != nil {
			PrintErr()
			fmt.Printf("%s : Ошибка записи найденных строк : %s\n", request, err)
			continue
		}
	}
}

/*

ЗАПИСЬ В ФАЙЛ

*/

func Writer(request string) {
	defer writerWG.Done()

	PrintInfo()
	fmt.Print("Удаление дублей и запись в файл\n")

	resultFileName := appDir + "/" + badSymbolsPattern.ReplaceAllString(request, "_") + ".txt"

	resultFile, err := os.OpenFile(resultFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		PrintErr()
		fmt.Printf("%s : Ошибка записи найденных строк : %s\n", request, err)
		return
	}

	resultFilesList = append(resultFilesList, resultFileName)

	writer := bufio.NewWriter(resultFile)

	switch saveType {
	case "1":
		for _, linePart := range requestStructMap[request].resultStrings {
			_, err = writer.WriteString(linePart[1] + ":" + linePart[2] + "\n")
		}
		_, err = writer.WriteString("\n")
	case "2":
		for _, linePart := range requestStructMap[request].resultStrings {
			_, err = writer.WriteString(linePart[0] + ":" + linePart[1] + ":" + linePart[2] + "\n")
		}
		_, err = writer.WriteString("\n")
	case "3":
	}

	if err != nil {
		PrintErr()
		fmt.Printf("%s : Ошибка записи найденных строк : %s\n", request, err)
		return
	}

	_ = resultFile.Close()
	PrintSuccess()
	fmt.Printf("Записано %d уникальных строк", 1)

}
