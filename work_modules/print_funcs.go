package work_modules

import (
	"fmt"
	"github.com/saintfish/chardet"
	"github.com/schollz/progressbar/v3"
	"os"
	"time"
)

// PrintErr PrintSuccess PrintWarn PrintInfo Значки

func PrintErr() {
	fmt.Print("[")
	ColorRed.Print("-")
	fmt.Print("] ")
}

func PrintSuccess() {
	fmt.Print("[")
	ColorGreen.Print("+")
	fmt.Print("] ")
}

func PrintWarn() {
	fmt.Print("[")
	ColorYellow.Print("*")
	fmt.Print("] ")
}

func PrintInfo() {
	fmt.Print("[")
	ColorMagenta.Print("*")
	fmt.Print("] ")
}

// PrintLinesChunk PrintCheckedFiles PrintFileInfo PrintFileDone Инфа о работе сортера

func PrintCheckedFiles() {
	fmt.Print("[")
	ColorBlue.Print(checkedFiles + 1)
	fmt.Print("/")
	ColorBlue.Print(len(filePathList))
	fmt.Print("] ")
}

func PrintFileInfo(path string) {
	PrintInfo()
	PrintCheckedFiles()
	fmt.Print("Обработка файла ")
	ColorBlue.Print(path)
	fmt.Print(" : ")
	if currentFileSize < 1610612736 {
		ColorBlue.Print(currentFileSize / 1048576)
		fmt.Print(" Мб : ")
	} else {
		ColorBlue.Print(currentFileSize / 1073741824)
		fmt.Print(" Гб : ")
	}
	ColorBlue.Print("~", currentFileLines)
	fmt.Print(" Строк\n")
}

func PrintFileDone(path string) {
	PrintSuccess()
	PrintCheckedFiles()
	ColorBlue.Print(path)
	fmt.Print(" : Файл обработан\n\n")
}

func PrintSortInfo() {
	fmt.Print("\n")
	switch {
	case reqLen <= 10:
		for _, request := range searchRequests {
			PrintSuccess()
			ColorBlue.Print(request)
			fmt.Print(" : ")
			ColorBlue.Print(sorterRequestStatMapCurrFile[request])
			fmt.Print(" строк\n")

		}
	case reqLen > 10:
		PrintSuccess()
		fmt.Print("Найдено ")
		ColorBlue.Print(currFileMatchLines)
		fmt.Print(" подходящих строк по всем запросам\n")
	}
}

func PrintClearInfo() {
	fmt.Print("\n")
	PrintInfo()
	ColorBlue.Print(TMPlinesLen)
	fmt.Print(" строк : ")
	ColorBlue.Print(currFileWritedString)
	fmt.Print(" Уникальных : ")
	ColorBlue.Print(currFileDubles)
	fmt.Print(" Повторов : ")
	ColorBlue.Print(currFileInvalidLen)
	fmt.Print(" Невалидных\n")
}

func PrintChunk() {
	PrintInfo()
	fmt.Print("Чтение по ")
	if GetAviableStringsCount() > currentFileLines {
		ColorBlue.Print(currentFileLines)
	} else {
		ColorBlue.Print(GetAviableStringsCount())
	}
	fmt.Print(" строк : ")
}

func PrintEncoding(result *chardet.Result) {
	ColorBlue.Print(result.Charset)
	fmt.Print(" - ")
	ColorBlue.Print(result.Confidence)
	fmt.Print(" %\n")
}

func CreatePBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(
		int(currentFileLines),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetItsString("Str"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("[*]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[blue]█[reset]",
			SaucerHead:    "[green]░[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

func PBarUpdater() {
	for isFileInProcessing {
		if TMPlinesLen > int(currentFileLines) {
			_ = pBar.Set64(currentFileLines)
		} else {
			_ = pBar.Set(TMPlinesLen)
		}
		time.Sleep(time.Millisecond * 250)
	}
}

// Ошибки

func PrintFileReadErr(path string, err error) {
	PrintErr()
	fmt.Printf("%s : Ошибка чтения файла : %s\n\n", path, err)
}

func PrintZeroRequestsErr() {
	PrintErr()
	fmt.Print("Нету запросов для сорта : Перезапустите сортер\n")
	PrintErr()
	fmt.Print("Нажмите ")
	ColorBlue.Print("Enter")
	fmt.Print(" для выхода")
	fmt.Scanln()
	os.Exit(1)
}

func PrintResultWriteErr(request string, err error) {
	PrintErr()
	ColorBlue.Print(request)
	fmt.Print(" : Ошибка записи найденных строк : ")
	ColorRed.Print(err, "\n")
	PrintInfo()
	fmt.Print("Запустите сортер с правами Администратора, если ошибка связана с доступом\n")
}

func PrintEncodingErr(err error) {
	PrintErr()
	fmt.Printf(" Ошибка определения кодировки : %s : Используется ", err)
	ColorBlue.Print("UTF-8")
}

func PrintEndodingLinesEnd() {
	PrintWarn()
	fmt.Print(" Недостаточно строк для определения кодировки : Используется ")
	ColorBlue.Print("UTF-8")
}

func PrintSorterResult() {

	fmt.Print("\n")
	for _, request := range searchRequests {
		PrintSuccess()
		ColorBlue.Print(request)
		fmt.Print(" : ")
		ColorBlue.Print(sorterRequestStatMap[request])
		fmt.Print(" строк : ")
		ColorBlue.Print(requestStructMap[request].resultFile, "\n")
	}
	fmt.Print("\n\n")

	PrintSuccess()
	fmt.Print("Файлов отсортировано : ")
	ColorBlue.Print(checkedFiles)
	fmt.Print(" из ")
	ColorBlue.Print(len(filePathList), "\n")

	PrintSuccess()
	fmt.Print("Строк отсортировано : ")
	ColorBlue.Print(checkedLines, "\n")

	PrintSuccess()
	fmt.Print("Подходящих строк : ")
	ColorGreen.Print(matchLines, "\n")

	PrintWarn()
	fmt.Print("Повторяющихся строк : ")
	ColorGreen.Print(sorterDubles, "\n\n")
}

func PrintCleanerResult() {
	fmt.Print("\n\n")
	PrintSuccess()
	fmt.Print("Файлов очищено : ")
	ColorBlue.Print(checkedFiles)
	fmt.Print(" из ")
	ColorBlue.Print(len(filePathList), "\n")
	PrintSuccess()
	fmt.Print("Повторов удалено : ")
	ColorBlue.Print(cleanerDublesLen, "\n")
	PrintSuccess()
	fmt.Print("Невалида удалено : ")
	ColorBlue.Print(cleanerInvalidLen, "\n")
	PrintSuccess()
	fmt.Print("Записано уникальных строк : ")
	ColorBlue.Print(cleanerWritedString, "\n\n")

	switch cleanType {
	case "1":
		for _, path := range filePathList {
			PrintSuccess()
			fmt.Print(cleanerOutputFilesMap[path] + "\n")
		}
	case "2":
		PrintSuccess()
		fmt.Print(cleanerOutputFilesMap[filePathList[0]] + "\n")
	}

}
