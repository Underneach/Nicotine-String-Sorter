package work_modules

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RemoveDublesResultFiles() {

	if len(resultFilesList) == 0 {
		fmt.Println("Ебаный рот где файлы")
	}

	if len(resultFilesList) == 1 {
		fullPath, err := filepath.Abs(resultFilesList[0])
		if err != nil {
			return
		}
		PrintSuccess()
		fmt.Print(fullPath)
		return
	}

	for _, path := range Unique(resultFilesList) {

		dublesWG.Add(1)

		err := dublesPool.Submit(func() {
			DublesRemove(path)
		})

		if err != nil {
			PrintErr()
			fmt.Print("Ошибка удаления дублей из полученных файлов : ", err)
		}
	}

	dublesWG.Wait()

}

func DublesRemove(path string) {
	defer dublesWG.Done()
	var lines []string
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	if err != nil {
		PrintErr()
		fmt.Printf("%s : Ошибка удаления дубликатов : %s\n", path, err)
		return
	}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	oldLen := len(lines)
	lines = Unique(lines)

	_, err = bufio.NewWriter(file).WriteString(strings.Join(lines, "\n") + "\n")
	if err != nil {
		PrintErr()
		fmt.Printf("%s : Ошибка удаления дубликатов : %s\n", path, err)
		return
	}
	file.Close()

	PrintSuccess()
	fullPath, err := filepath.Abs(path)
	if err != nil {
		ColorBlue.Print(request)
		fmt.Print(" : ")
	} else {
		fmt.Print(fullPath, "\n")
	}
	fmt.Print("Всего записано уникальных строк : ")
	ColorBlue.Print(len(lines))
	fmt.Print(" : Всего удалено дубликатов : ")
	ColorBlue.Print(oldLen-len(lines), "\n\n")
}
