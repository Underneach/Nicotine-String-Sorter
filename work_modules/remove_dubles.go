package work_modules

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RemoveDublesResultFiles() {

	if len(filePathList) == 1 {
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
	_, _ = ColorBlue.Print(request)
	fmt.Print(" : Всего записано уникальных строк : ")
	_, _ = ColorBlue.Print(len(lines))
	fmt.Print(" : Всего удалено дубликатов : ")
	_, _ = ColorBlue.Print(oldLen-len(lines), "\n")
}
