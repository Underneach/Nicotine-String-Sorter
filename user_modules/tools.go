package user_modules

import (
	"fmt"
	"os"
	"os/exec"
)

func GetFilesSize(flist []string) {
	for _, path := range flist {

		if info, err := os.Stat(path); err != nil {
			PrintErr()
			fmt.Printf("%s : Ошибка получения размера файла : %s\n", path, err)
			continue
		} else {
			filesSize += info.Size()
		}
	}
}

func ClearTerm() {
	switch userOs {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested 
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	case "linux":
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	default:
		fmt.Println("\033[2J")
	}
}

func Unique(slice []string) []string {

	inResult := make(map[string]bool)
	var result []string
	for _, str := range slice {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func SetTermTitle(appVersion string) {
	var cmd *exec.Cmd

	switch userOs {
	case "windows":
		cmd = exec.Command("powershell", "-Command", "& { $Host.UI.RawUI.WindowTitle = '"+"Nicotine String Sorter | НикотиновыйКодер | "+appVersion+"' }")
	case "linux":
		cmd = exec.Command("bash", "-c", "echo -ne '\\033]0;"+"Nicotine String Sorter | НикотиновыйКодер | "+appVersion+"\\007'")
	}

	if err := cmd.Run(); err != nil {
		PrintErr()
		fmt.Print(err, "\n")
	}
}
