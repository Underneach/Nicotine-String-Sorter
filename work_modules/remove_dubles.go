package work_modules

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RemoveDublesResultFiles() {
	for _, path := range Unique(resultFilesList) {

		var lines []string
		file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		if err != nil {
			PrintErr()
			fmt.Printf("%s : Ошибка удаления дубликатов : %s\n", path, err)
			continue
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
			continue
		}
		file.Close()
		PrintSuccess()
		fmt.Printf("%s : Записано уникальных %d строк : Удалено %d дубликатов\n", request, len(lines), oldLen-len(lines))
	}
}
