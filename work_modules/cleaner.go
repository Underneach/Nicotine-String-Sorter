package work_modules

import (
	"bufio"
	"fmt"
	"github.com/zeebo/xxh3"
	"golang.org/x/text/transform"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func RunCleaner() {

	PrintInfo()
	fmt.Print("Запуск Клинера...")

	for _, path := range filePathList {
		cleanerOutputFilesMap[path] = GetRunDir() + strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + "_cleaned.txt"
	}
	fmt.Print("\r")
	PrintSuccess()
	fmt.Print("Клинер запущен   \n\n")

}

func Cleaner(path string) {
	currPath = path
	cleanerStringChannelMap[currPath] = make(chan string)
	cleanerResultChannelMap[currPath] = make(chan string)
	cleanerStringHashMap = make(map[uint64]bool)
	TMPlinesLen = 0
	currFileDubles = 0
	currFileWritedString = 0
	currFileInvalidLen = 0

	if err := GetCurrentFileSize(path); err != nil {
		PrintFileReadErr(path, err)
		return
	}

	PrintFileInfo(path)
	PrintLinesChunk()
	fileDecoder = GetEncodingDecoder(path)

	cleanerReadFile, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)

	if err != nil {
		PrintFileReadErr(path, err)
		return
	}

	cleanerWriteFile, err = os.OpenFile(cleanerOutputFilesMap[currPath], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		PrintFileReadErr(path, err)
		return
	}

	if i := GetAviableStringsCount(); i > currentFileLines {
		cleanerPool.Tune(int(math.Round(float64(currentFileLines) / 3)))
	} else {
		cleanerPool.Tune(int(math.Round(float64(i) / 3)))
	}

	scanner := bufio.NewScanner(transform.NewReader(cleanerReadFile, fileDecoder))
	isFileInProcessing = true

	go PBarUpdater()
	go CleanerProcessInputLines()
	go CleanerWriteLine()

	for ; scanner.Scan(); TMPlinesLen++ {
		workWG.Add(1)
		cleanerStringChannelMap[currPath] <- scanner.Text()
	}

	workWG.Wait()                               // Ждем горутины
	isFileInProcessing = false                  // Останавливаем пбар
	close(cleanerStringChannelMap[currPath])    // Закрываем каналы
	close(cleanerResultChannelMap[currPath])    //
	checkedLines += int64(TMPlinesLen)          // Прибавляем строки
	cleanerDublesLen += currFileDubles          //
	cleanerWritedString += currFileWritedString //
	cleanerInvalidLen += currFileInvalidLen     //
	_ = pBar.Finish()                           // Завершаем бар
	_ = pBar.Exit()                             // Закрываем бар
	cleanerReadFile.Close()                     // Закрываем файл
	cleanerWriteFile.Close()                    // Закрываем файл
	cleanerStringChannelMap[currPath] = nil     // Чистим канал
	cleanerResultChannelMap[currPath] = nil     // 
	cleanerStringHashMap = nil                  //
	fmt.Print("\n")                             //
	PrintClearInfo()                            //
	PrintFileSorted(path)                       // Пишем файл отсортрован
	checkedFiles++                              // Прибавляем пройденные файлы
}

func CleanerProcessInputLines() {
	for {
		if data, ok := <-cleanerStringChannelMap[currPath]; !ok {
			break
		} else {
			_ = cleanerPool.Invoke(data)
			continue
		}
	}
}

func CleanerProcessString(line string) {
	defer workWG.Done()
	if validPattern.MatchString(line) && !uncknownPattern.MatchString(line) {
		hash := xxh3.HashString(line)
		CHMMutex.Lock()
		if _, ok := cleanerStringHashMap[hash]; !ok {
			cleanerStringHashMap[hash] = true
			cleanerResultChannelMap[currPath] <- line
		} else {
			currFileDubles++
		}
		CHMMutex.Unlock()
	} else {
		currFileInvalidLen++
	}
}

func CleanerWriteLine() {
	for {
		if data, ok := <-cleanerResultChannelMap[currPath]; !ok {
			break
		} else {
			_, _ = cleanerWriteFile.WriteString(data + "\n")
			currFileWritedString++
			continue
		}
	}
}
