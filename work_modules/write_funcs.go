package work_modules

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"os"
	"strings"
)

/*

	ЗАПИСЬ РЕЗУЛЬТАТА СОРТА

*/

func WriteResult() {

	fmt.Print("\n")
	PrintInfo()
	fmt.Print("Запись строк в файл\n")

	for _, req := range searchRequests {
		writerWG.Add(1)
		if wrterr := writerPool.Submit(func() {
			Writer(req, requestStructMap[req].resultStrings, requestStructMap[req].resultFile) // Юзать мютексы? Да пошли они нахуй!
		}); wrterr != nil {
			PrintResultWriteErr(req, wrterr)
			writerWG.Done()
			continue
		}
	}

	writerWG.Wait()
	PrintSortInfo()
}

func Writer(request string, lines []string, path string) {
	defer writerWG.Done()
	fmt.Println(path)

	if len(lines) == 0 {
		PrintErr()
		ColorBlue.Print(request)
		fmt.Print(" : Нет строк для записи\n")
		return
	}

	resultFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		PrintResultWriteErr(request, err)
		return
	}

	if _, err := bufio.NewWriter(transform.NewWriter(resultFile, unicode.UTF8.NewDecoder())).WriteString(strings.Join(lines, "\n")); err != nil {
		PrintResultWriteErr(request, err)
		return
	}

	_ = resultFile.Close()
	clear(lines)
}

/*

УДАЛЕНИЕ ДУБЛЕЙ ИЗ РЕЗУЛЬТАТА

*/

func RemoveDublesResultFiles() {

	dublesWG.Add(reqLen)
	for _, request := range searchRequests {
		_ = dublesPool.Submit(func() {
			DublesRemove(request, requestStructMap[request].resultFile)
			dublesWG.Done()
		})
	}
	dublesWG.Wait()
}

func DublesRemove(request string, path string) {
	fmt.Println(request)

	var lines []string
	rdfile, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	scanner := bufio.NewScanner(rdfile)
	scanner.Split(bufio.ScanLines)

	if os.IsNotExist(err) {
		PrintErr()
		ColorBlue.Print(request)
		fmt.Print(" : Нет файла для удаления дублей\n\n")
		return
	} else if err != nil {
		PrintErr()
		ColorBlue.Print(request)
		fmt.Print(" : ")
		ColorRed.Print(err, "\n\n")
		return
	}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	rdfile.Close()

	oldLen := len(lines)
	lines = Unique(lines)

	// Открываем файл и чистим
	wrfile, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		PrintRemoveDublesErr(request, err)
		return
	}

	_ = wrfile.Truncate(0)
	_, _ = wrfile.Seek(0, 0)

	if _, err = bufio.NewWriter(wrfile).WriteString(strings.Join(lines, "\n") + "\n"); err != nil {
		PrintRemoveDublesErr(request, err)
		wrfile.Close()
		return
	}
	wrfile.Close()

	dublesMutex.Lock()
	PrintSuccess()
	ColorBlue.Print(request)
	fmt.Print(" : ")
	ColorBlue.Print(path, "\n")

	PrintInfo()
	fmt.Print("Уникальных строк : ")
	ColorBlue.Print(len(lines))
	fmt.Print(" : Дубликатов : ")
	ColorBlue.Print(oldLen-len(lines), "\n\n")
	dublesMutex.Unlock()
	clear(lines)
}
