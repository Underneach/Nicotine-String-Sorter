package work_modules

import (
	"bufio"
	"fmt"
	"golang.org/x/text/transform"
	"os"
	"time"
)

func Sorter(path string) {
	currPath = path
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
	go PBarUpdater()
	go ProcessResult()

	scanner := bufio.NewScanner(transform.NewReader(file, fileDecoder))

	for ; scanner.Scan(); TMPlinesLen++ {
		if TMPlinesLen >= GetAviableStringsCount() {
			sorterWG.Add(TMPlinesLen)
			SendLinesToPool(tmpLines)
			currFileCheckedLines += TMPlinesLen
			TMPlinesLen = 0
			clear(tmpLines)
		} else {
			tmpLines = append(tmpLines, scanner.Text())
		}
	}

	if len(tmpLines) > 0 {
		sorterWG.Add(TMPlinesLen)
		SendLinesToPool(tmpLines)
		currFileCheckedLines += TMPlinesLen
		TMPlinesLen = 0
		clear(tmpLines)
	}
	close(fileChannelMap[currPath])

	checkedLines += int64(currFileCheckedLines) // Прибавляем строки
	_ = pBar.Finish()                           // Завершаем бар
	_ = pBar.Exit()                             // Закрываем бар

	isFileInProcessing = false
	for isResultWrited == false {
		time.Sleep(time.Millisecond * 100)
	}

	file.Close() // Закрываем файл

	for _, request := range searchRequests {
		clear(requestStructMap[request].resultStrings) // Чистим список
	}

	PrintFileSorted(path)                   // Пишем файл отсортрован
	checkedFiles++                          // Прибавляем пройденные файлы
	invalidLines += currentFileInvalidLines // Суммируем невалид строки
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
		currentFileInvalidLines++
		return
	}

	for _, request := range searchRequests {
		result := requestStructMap[request].requestPattern.FindStringSubmatch(line)
		if len(result) == 2 {
			fileChannelMap[currPath] <- [2]string{request, result[1]}
		}
	}
}

func ProcessResult() {
	/*

		Обработка результата отдельная ебатория, в питоне все результаты работы из ThreadPoolExecuror сохранялись структом в список,
		и эта хуйня жрала кучу памяти. Делать пул в данном случае смысла нет, карта не умеет в потокобезопасность и будет сосать бибу.
		Делать sync.Map - нахуй пойдет скорость. Сейчас реализована дефолт FIFO очередь, самый оптимальный подход по моему мнению.


	*/

	ResultListMap := make(map[string][]string)

	for _, request := range searchRequests {
		ResultListMap[request] = requestStructMap[request].resultStrings
	}

	for {
		if data, ok := <-fileChannelMap[currPath]; !ok {
			break
		} else {
			ResultListMap[data[0]] = append(ResultListMap[data[0]], data[1])
			continue
		}
	}

	for _, request := range searchRequests {
		matchLines += int64(len(ResultListMap[request]))
		requestStructMap[request].resultStrings = ResultListMap[request]
	}

	clear(ResultListMap)  // чистим список
	WriteResult()         // Пишем результат в файл
	isResultWrited = true // сообщаем о том, что файл записан
}
