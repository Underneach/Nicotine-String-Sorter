package user_modules

import (
	"fmt"
	"os"
	"os/exec"
)

func GetFilesSize(flist []string) {
	for _, path := range flist {
		info, err := os.Stat(path)
		if err != nil {
			PrintErr()
			fmt.Printf("%s : Ошибка получения размера файла : %s\n", path, err)
			continue
		}

		filesSize = +info.Size()
	}
}

func ClearTerm() {
	switch userOs {
	case "windows":
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	case "linux":
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested 
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
