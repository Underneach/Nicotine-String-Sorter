package work_modules

import (
	"fmt"
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

// PrintLinesChunk PrintCheckedFiles PrintFileInfo PrintFileSorted Инфа о работе сортера

func PrintLinesChunk() {
	PrintSuccess()
	fmt.Print("Чтение файла по ")
	if GetAviableStringsCount() > int(currentFileLines) {
		ColorBlue.Print(currentFileLines)
	} else {
		ColorBlue.Print(GetAviableStringsCount())
	}
	fmt.Print(" строк\n")
}

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
	fmt.Print("Сортировка файла ")
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

func PrintFileSorted(path string) {
	PrintInfo()
	PrintCheckedFiles()
	ColorBlue.Print(path)
	fmt.Print(" : Файл отсортирован\n\n")
}

func PrintSortInfo() {
	switch {
	case reqLen <= 10:
		for _, request := range searchRequests {
			strLen := len(requestStructMap[request].resultStrings)
			if strLen > 0 {
				PrintSuccess()
				ColorBlue.Print(request)
				fmt.Print(" : ")
				ColorBlue.Print(strLen)
				fmt.Print(" строк\n")
			}
		}
	case reqLen > 10:
		PrintSuccess()
		fmt.Print("Найдено ")
		ColorBlue.Print(matchLines)
		fmt.Print(" подходящих строк\n")
	}
	PrintWarn()
	ColorYellowLight.Print("Невалид")
	fmt.Print(" : ")
	ColorYellowLight.Print(currentFileInvalidLines)
	fmt.Print(" строк\n")
}

func CreatePBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(
		int(currentFileLines),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetItsString("Str"),
		progressbar.OptionSetRenderBlankState(true),
	)
}

func PBarUpdater() {
	pBar = CreatePBar()
	for isFileInProcessing {
		if currFileCheckedLines > int(currentFileLines) {
			_ = pBar.Set64(currentFileLines)
		} else {
			_ = pBar.Set(currFileCheckedLines)
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

func PrintRemoveDublesErr(request string, err error) {
	PrintErr()
	ColorBlue.Print(request)
	fmt.Print(" : Ошибка удаления дублей : ")
	ColorRed.Print(err, "\n")
}
