package work_modules

import (
	"bufio"
	"fmt"
	"golang.org/x/text/transform"
	"io"
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

	reader := bufio.NewReader(transform.NewReader(file, fileDecoder))
	bar := CreateBar()

	var rderr error

Loop:
	for {
		line, rderr = reader.ReadString('\n')
		switch rderr {
		case nil:
			continue Loop
		case io.EOF:
			break Loop
		}
		fmt.Println(line)
		readLines = append(readLines, line)
	}

	linesLen := len(readLines)                    // Кол во строк в слайсе
	sorterWG.Add(linesLen)                        //	 	
	ProcessLines(readLines)                       // ОТПРАВЛЯЕМ СТРОКИ В ПУЛ И ЖДЁМ
	sorterWG.Wait()                               // 
	checkedLines = checkedLines + int64(linesLen) // Прибавляем строки
	clear(readLines)                              // Чистим строки
	barerr := bar.Add(linesLen)                   // Увеличиваем бар
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
			Worker(line)
		})
		if err != nil {
			PrintErr()
			fmt.Print("Ошибка отправки строки в пул : ", err, "\n")
		}
	}
}

func Worker(line string) {
	defer sorterWG.Done()
	fmt.Println(line)
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
	_, _ = ColorBlue.Print(request)
	fmt.Print(" : Записано ")
	_, _ = ColorBlue.Print(len(requestStructMap[request].resultStrings))
	fmt.Print(" уникальных строк\n")

}
