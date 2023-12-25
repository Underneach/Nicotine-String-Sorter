package work_modules

import (
	"fmt"
	"os"
)

// PrintErr PrintSuccess PrintWarn PrintInfo Значки

func PrintErr() {
	fmt.Print("[")
	_, _ = ColorRed.Print("-")
	fmt.Print("] ")
}

func PrintSuccess() {
	fmt.Print("[")
	_, _ = ColorGreen.Print("+")
	fmt.Print("] ")
}

func PrintWarn() {
	fmt.Print("[")
	_, _ = ColorYellow.Print("*")
	fmt.Print("] ")
}

func PrintInfo() {
	fmt.Print("[")
	_, _ = ColorMagenta.Print("*")
	fmt.Print("] ")
}

// PrintLinesChunk PrintCheckedfiles PrintFileInfo PrintFileSorted Инфа о работе сортера

func PrintLinesChunk() {
	PrintInfo()
	fmt.Print("Чтение файла по ")
	if GetAviableStringsCount() > int(currentFileLines) {
		_, _ = ColorBlue.Print(currentFileLines)
	} else {
		_, _ = ColorBlue.Print(GetAviableStringsCount())
	}
	fmt.Print(" строк\n")
}

func PrintCheckedfiles() {
	fmt.Print("[")
	_, _ = ColorBlue.Print(checkedFiles + 1)
	fmt.Print("/")
	_, _ = ColorBlue.Print(len(filePathList))
	fmt.Print("] ")
}

func PrintFileInfo(path string) {
	PrintInfo()
	PrintCheckedfiles()
	fmt.Print("Сортировка файла ")
	_, _ = ColorBlue.Print(path)
	fmt.Print(" : ")
	if currentFileSize < 1610612736 {
		_, _ = ColorBlue.Print(currentFileSize / 1048576)
		fmt.Print(" Мб : ")
	} else {
		_, _ = ColorBlue.Print(currentFileSize / 1073741824)
		fmt.Print(" Гб : ")
	}
	_, _ = ColorBlue.Print("~", currentFileLines)
	fmt.Print(" Строк\n")
}

func PrintFileSorted(path string) {
	PrintSuccess()
	PrintCheckedfiles()
	_, _ = ColorBlue.Print(path)
	fmt.Print(" : Файл отсортирован\n\n")
}

// PrintFileReadErr PrintZeroRequestsErr Ошибки

func PrintFileReadErr(path string, err error) {
	PrintErr()
	fmt.Printf("%s : Ошибка чтения файла : %s\n\n", path, err)
}

func PrintZeroRequestsErr() {
	PrintErr()
	fmt.Print("Нету запросов для сорта : Перезапустите сортер\n")
	PrintErr()
	fmt.Print("Нажмите ")
	_, _ = ColorBlue.Print("Enter")
	fmt.Print(" для выхода")
	_, _ = fmt.Scanln()
	os.Exit(1)
}
