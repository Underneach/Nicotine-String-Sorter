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
	ColorBlue.Print("   v")
	ColorMagenta.Print(appVersion)
	fmt.Print(" | ")
	ColorMagenta.Print(runtime.Version())
	ColorBlue.Print("     #")
	fmt.Print(" t.me/rx580_work     ")
	ColorGreen.Print("#")
	fmt.Print(" zelenka.guru/rx580    # НикотиновыйКодер\n\n")
	PrintInfo()
	fmt.Print(cpuid.CPU.BrandName, " @ ", math.Round(float64(cpuid.CPU.BoostFreq/1000000000)), "GHz @ ", cpuid.CPU.PhysicalCores, "/", cpuid.CPU.LogicalCores, " потоков\n")
	PrintInfo()
	fmt.Print(math.Round(float64(memory.FreeMemory()/1073741824)), "/", math.Round(float64(memory.TotalMemory()/1073741824)), " Гб доступной памяти\n\n")
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
	ColorBlue.Print("   v")
	ColorMagenta.Print(appVersion)
	fmt.Print(" | ")
	ColorMagenta.Print(runtime.Version())
	ColorBlue.Print("     #")
	fmt.Print(" t.me/rx580_work     ")
	ColorGreen.Print("#")
	fmt.Print(" zelenka.guru/rx580    # НикотиновыйКодер\n\n")
}

func PrintInputData(appVersion string) string {
	ClearTerm()
	PrintLogoFast(appVersion)
	PrintInfo()
	fmt.Print("Всего файлов : ")
	ColorBlue.Print(len(filePathList))
	fmt.Print(" : Объем : ")
	if filesSize < 1610612736 {
		ColorBlue.Print(filesSize / 1048576)
		fmt.Print(" Мб ")
	} else {
		ColorBlue.Print(filesSize / 1073741824)
		fmt.Print(" Гб ")
	}

	fmt.Print(": Строк : ")
	ColorBlue.Print("~", filesSize/80, "\n")

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

	PrintInput()
	fmt.Print("Выберите действие:\n\n")

	ColorBlue.Print("	1")
	fmt.Print(" : запустить сортер\n")
	ColorBlue.Print("	2")
	fmt.Print(" : ввести данные заново\n\n")
	for true {
		fmt.Print("> ")
		userSelect, _ := userInputReader.ReadString('\n')
		userSelect = strings.TrimSpace(userSelect)
		if userSelect == "1" {
			returnData = "continue"
			break
		} else if userSelect == "2" {
			returnData = "restart"
			break
		}
	}
	ClearTerm()
	return returnData
}

func PrintResult(Duration time.Duration, checkedLines int64, invalidLines int64, resultLinesCount int64, totalFiles int, checkedFiles int) {

	fmt.Print("\n\n")
	PrintSuccess()
	fmt.Print("Файлов отсортировано : ")
	ColorBlue.Print(checkedFiles)
	fmt.Print(" из ")
	ColorBlue.Print(totalFiles, "\n")

	PrintSuccess()
	fmt.Print("Строк отсортировано : ")
	ColorBlue.Print(checkedLines, "\n")

	PrintSuccess()
	fmt.Print("Подходящих строк : ")
	ColorBlue.Print(resultLinesCount, "\n")

	PrintSuccess()
	fmt.Print("Невалидных строк : ")
	ColorBlue.Print(invalidLines, "\n")

	fmt.Print("\n")
	PrintSuccess()
	fmt.Print("Время выполнения : ")
	ColorBlue.Print(Duration, "\n\n\n")

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
