package user_modules

import (
	"fmt"
	"github.com/klauspost/cpuid/v2"
	"github.com/pbnjay/memory"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

func PrintLogoStart(appVersion string) {

	ColorBlue.Print(`
     _   _   _                  _     _                        
    | \ | | (_)                | |   (_)                 
    |  \| |  _    ___    ___   | |_   _   _ __     ___    
    | .   | | |  / __|  / _ \  | __| | | | '_ \   / _ \   
    | |\  | | | | (__  | (_) | | |_  | | | | | | |  __/  
    |_| \_| |_|  \___|  \___/   \__| |_| |_| |_|  \___|
														`)
	time.Sleep(300 * time.Millisecond)
	ColorBlue.Print(`
     _____   _           _                       _____                  _ 
    / ____| | |         (_)                     / ____|                | |              
   | (___   | |_   _ __   _   _ __     __ _    | (___     ___    _ __  | |_    ___   _ __ 
    \___ \  | __| | '__| | | | '_ \   / _  |    \___ \   / _ \  | '__| | __|  / _ \ | '__|
    ____) | | |_  | |    | | | | | | | (_| |    ____) | | (_) | | |    | |_  |  __/ | |   
   |_____/   \__| |_|    |_| |_| |_|  \__, |   |_____/   \___/  |_|     \__|  \___| |_|   
                                       __/ |    
                                      |___/  

`)
	time.Sleep(150 * time.Millisecond)
	ColorMagenta.Print("    v", appVersion)
	fmt.Print(" | ")
	ColorMagenta.Print(runtime.Version())
	ColorBlue.Print("     #")
	fmt.Print(" t.me/rx580_work     ")
	ColorGreen.Print("#")
	fmt.Print(" zelenka.guru/rx580    # НикотиновыйКодер\n\n")
	PrintInfo()
	fmt.Print(cpuid.CPU.BrandName, " @ ", cpuid.CPU.PhysicalCores, "/", cpuid.CPU.LogicalCores, " потоков | ")
	fmt.Print(math.Round(float64(memory.FreeMemory()/1073741824)), "/", math.Ceil(float64(memory.TotalMemory()/1073741824)), " Гб доступной памяти\n\n")
	isLogoPrinted = true
}

func PrintLogoFast(appVersion string) {

	ColorBlue.Print(`
     _   _   _                  _     _                        
    | \ | | (_)                | |   (_)                 
    |  \| |  _    ___    ___   | |_   _   _ __     ___    
    | .   | | |  / __|  / _ \  | __| | | | '_ \   / _ \   
    | |\  | | | | (__  | (_) | | |_  | | | | | | |  __/  
    |_| \_| |_|  \___|  \___/   \__| |_| |_| |_|  \___|

     _____   _           _                       _____                  _ 
    / ____| | |         (_)                     / ____|                | |              
   | (___   | |_   _ __   _   _ __     __ _    | (___     ___    _ __  | |_    ___   _ __ 
    \___ \  | __| | '__| | | | '_ \   / _  |    \___ \   / _ \  | '__| | __|  / _ \ | '__|
    ____) | | |_  | |    | | | | | | | (_| |    ____) | | (_) | | |    | |_  |  __/ | |   
   |_____/   \__| |_|    |_| |_| |_|  \__, |   |_____/   \___/  |_|     \__|  \___| |_|   
                                       __/ |    
                                      |___/ 
 
`)
	ColorMagenta.Print("    v", appVersion)
	fmt.Print(" | ")
	ColorMagenta.Print(runtime.Version())
	ColorBlue.Print("     #")
	fmt.Print(" t.me/rx580_work     ")
	ColorGreen.Print("#")
	fmt.Print(" zelenka.guru/rx580    # НикотиновыйКодер\n\n")
}

func PrintInputInfo(appVersion string) {
	ClearTerm()
	PrintLogoFast(appVersion)

	PrintInfo()
	fmt.Print("Всего файлов : ")
	ColorBlue.Print(len(filePathList))
	fmt.Print(" : Объем : ")
	ColorBlue.Print(filesSize / 1048576)
	fmt.Print(" Мб ")
	fmt.Print(": Строк : ")
	ColorBlue.Print("~", filesSize/80, "\n")
}
func PrintSorterData() {
	PrintInfo()
	fmt.Printf("Всего запросов : ")

	reqLen := len(searchRequests)

	switch {
	case reqLen <= 3:
		ColorBlue.Print(reqLen)
		fmt.Print(" : ")
		for i, req := range searchRequests {
			ColorBlue.Print(req)
			if i != reqLen-1 {
				fmt.Print(", ")
			}
		}
		fmt.Print("\n")
	case reqLen > 3 && reqLen <= 10:
		ColorBlue.Print(reqLen, "\n")
		for _, request := range searchRequests {
			fmt.Println("    ", request)
		}
		fmt.Print("\n")
	case reqLen > 10:
		ColorBlue.Print(reqLen, "\n\n")

	}
	PrintInfo()
	fmt.Print("Формат сохранениея : ")
	switch saveType {
	case "1":
		ColorBlue.Print("Log:Pass\n")
	case "2":
		ColorBlue.Print("Url:Log:Pass\n")
	}
}

func PrintInputData(appVersion string) (uSelect string) {
LoopData:
	for true {
		PrintInputInfo(appVersion)

		switch workMode {
		case "sorter":
			PrintSorterData()
		case "cleaner":
		}

		PrintInput()
		fmt.Print("Выберите действие:\n\n")

		ColorBlue.Print("	1")
		fmt.Print(" - Запустить\n")
		if workMode == "sorter" {
			ColorBlue.Print("	2")
			fmt.Print(" - Выбрать разделитель строк - '")
			ColorBlue.Print(":")
			fmt.Print("' по умолчанию\n")
		}
		ColorBlue.Print("	3")
		fmt.Print(" - Ввести данные заново\n\n")

	LoopMode:
		for true {
			fmt.Print("> ")
			userSelect, _ := userInputReader.ReadString('\n')
			userSelect = strings.TrimSpace(userSelect)

			switch userSelect {
			case "1":
				uSelect = "continue"
				break LoopMode
			case "2":
				if workMode == "sorter" {
					delimetr = GetDelimetrInput()
					continue LoopData
				} else {
					continue LoopData
				}
			case "3":
				uSelect = "restart"
				break LoopMode
			default:
				continue LoopMode
			}
		}
		break LoopData
	}
	ClearTerm()
	return uSelect
}

func PrintTimeDuration(duration time.Duration) {
	fmt.Print("\n")
	PrintSuccess()
	fmt.Print("Время сортировки : ")
	ColorBlue.Print(duration, "\n\n\n")

	PrintInfo()
	fmt.Print("Нажмите ")
	ColorBlue.Print("Enter")
	fmt.Print(" для выхода")
	fmt.Scanln()
	os.Exit(0)
}

func PrintInput() {
	fmt.Print("[")
	ColorBlue.Print("#")
	fmt.Print("] ")
}

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

func _() {
	fmt.Print("[")
	ColorYellow.Print("*")
	fmt.Print("] ")
}

func PrintInfo() {
	fmt.Print("[")
	ColorMagenta.Print("*")
	fmt.Print("] ")
}

func PrintWorkModes() {
	PrintInfo()
	fmt.Print("Поддерживаемые типы работы:\n\n")
	ColorBlue.Print("       1")
	fmt.Print(" - Сортер строк\n")
	fmt.Print("       Поиск строк в базе подходящих под запросы и запись в отдельный файл с удалением повторов\n")
	fmt.Print("       Запрос должен быть в формате 'google.com' или 'google'\n\n")
	ColorBlue.Print("       2")
	fmt.Print(" - Клинер базы от невалид строк и дубликатов\n")
	fmt.Print("       Удаление повторов и строк не подходящих под 'A-z / 0-9 / Специмволы | 10-256 символов | без UNKNOWN'\n\n")
	ColorBlue.Print("       4")
	fmt.Print(" - Закрыть программу\n\n")
}
