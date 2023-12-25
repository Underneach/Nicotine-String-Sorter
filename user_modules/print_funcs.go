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

func PrintLogoAnimation(appVersion string) {

	_, _ = ColorBlue.Print(`
     _   _   _                  _     _                        
    | \ | | (_)                | |   (_)                 
    |  \| |  _    ___    ___   | |_   _   _ __     ___    
    | .   | | |  / __|  / _ \  | __| | | | '_ \   / _ \   
    | |\  | | | | (__  | (_) | | |_  | | | | | | |  __/  
    |_| \_| |_|  \___|  \___/   \__| |_| |_| |_|  \___|
														`)
	time.Sleep(300 * time.Millisecond)
	_, _ = ColorBlue.Print(`
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
	_, _ = ColorMagenta.Printf("                 v%s | %s\n\n\n", appVersion, runtime.Version())

}

func PrintLogoFast(appVersion string) {

	_, _ = ColorBlue.Print(`
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
	_, _ = ColorMagenta.Printf("                 v%s | %s\n\n\n", appVersion, runtime.Version())

}

func PrintAutorLinks() {
	_, _ = ColorBlue.Print("#")
	fmt.Print(" t.me/rx580_work     ")
	_, _ = ColorGreen.Print("#")
	fmt.Print(" zelenka.guru/rx580    # НикотиновыйКодер\n\n")

}

func PrintUserMachineSpecs() {
	PrintInfo()
	fmt.Print(cpuid.CPU.BrandName, " @ ", math.Round(float64(cpuid.CPU.BoostFreq/1000000000)), "GHz @ ", cpuid.CPU.PhysicalCores, "/", cpuid.CPU.LogicalCores, " потоков\n")
	PrintInfo()
	fmt.Print(math.Round(float64(memory.FreeMemory()/1073741824)), "/", math.Round(float64(memory.TotalMemory()/1073741824)), " Гб доступной памяти\n\n")
}

func PrintInputData() string {
	ClearTerm()
	PrintInfo()
	fmt.Print("Всего файлов : ")
	_, _ = ColorBlue.Print(len(filePathList))
	fmt.Print(" : Объем : ")
	if filesSize < 1610612736 {
		_, _ = ColorBlue.Print(filesSize / 1048576)
		fmt.Print(" Мб ")
	} else {
		_, _ = ColorBlue.Print(filesSize / 1073741824)
		fmt.Print(" Гб ")
	}

	fmt.Print(": Строк : ")
	_, _ = ColorBlue.Print("~", filesSize/80, "\n")

	PrintInfo()
	fmt.Printf("Всего запросов : ")

	reqLen := len(searchRequests)

	switch {
	case reqLen <= 3:
		_, _ = ColorBlue.Print(reqLen)
		fmt.Print(" : ")
		for _, req := range searchRequests {
			_, _ = ColorBlue.Print(req)
			fmt.Print(", ")
		}
		fmt.Print("\n")
	case reqLen > 3 && reqLen <= 10:
		_, _ = ColorBlue.Print(reqLen, "\n")
		for _, request := range searchRequests {
			fmt.Println("    ", request)
		}
		fmt.Print("\n")
	case reqLen > 10:
		_, _ = ColorBlue.Print(reqLen, "\n\n")

	}

	PrintInput()
	fmt.Print("Выберите действие:\n\n")

	_, _ = ColorBlue.Print("	1")
	fmt.Print(" : запустить сортер\n")
	_, _ = ColorBlue.Print("	2")
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
	return returnData
}

func PrintResult(Duration time.Duration, checkedLines int64, invalidLines int64, resultLinesCount int64, totalFiles int, checkedFiles int) {

	PrintInfo()
	fmt.Print("Файлов отсортировано : ")
	_, _ = ColorBlue.Print(checkedFiles)
	fmt.Print(" из ")
	_, _ = ColorBlue.Print(totalFiles, "\n")

	PrintInfo()
	fmt.Print("Строк отсортировано : ")
	_, _ = ColorBlue.Print(checkedLines, "\n")

	PrintInfo()
	fmt.Print("Подходящих строк : ")
	_, _ = ColorBlue.Print(resultLinesCount, "\n")

	PrintInfo()
	fmt.Print("Невалидных строк : ")
	_, _ = ColorBlue.Print(invalidLines, "\n")

	fmt.Print("\n\n\n")
	PrintInfo()
	fmt.Print("Время выполнения : ")
	_, _ = ColorBlue.Print(Duration, "\n\n\n")

	PrintInfo()
	fmt.Print("Нажмите ")
	_, _ = ColorBlue.Print("Enter")
	fmt.Print("для выхода")
	_, _ = fmt.Scanln()
	os.Exit(0)
}

func PrintInput() {
	fmt.Print("[")
	_, _ = ColorBlue.Print("#")
	fmt.Print("] ")
}

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

func _() {
	fmt.Print("[")
	_, _ = ColorYellow.Print("*")
	fmt.Print("] ")
}

func PrintInfo() {
	fmt.Print("[")
	_, _ = ColorMagenta.Print("*")
	fmt.Print("] ")
}
