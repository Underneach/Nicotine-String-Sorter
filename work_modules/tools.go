package work_modules

import (
	"bufio"
	"fmt"
	"github.com/pbnjay/memory"
	"github.com/saintfish/chardet"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetAviableStringsCount() int {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if time.Since(lastUpdate) > time.Minute { // Если прошло более минуты с момента последнего обновления, обновляем кеш
		cachedStrCount = getAviableStringsCountCached()
		lastUpdate = time.Now()
	}
	return cachedStrCount
}

func getAviableStringsCountCached() int {
	freeMemory := memory.FreeMemory()
	if freeMemory != 0 {
		return int(math.Round(float64(freeMemory / (80 * 4 * 0.90)))) // 80 - Предпологаемый размер строки, 4 - размер символа в байтах, 0.90 - оставляем часть памяти для других элементов сортера
	} else {
		PrintWarn()
		fmt.Print("Не удалось получить количество доступной памяти : Чтение по чанкам в 2Гб")
		return 6700000 // Возвращаем ~2 гига, если не получили доступный размер
	}
}

func GetEncodingDecoder(path string) *encoding.Decoder {

	var detectedEncoding encoding.Encoding
	var decoder *encoding.Decoder

	detector := chardet.NewTextDetector()
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		PrintErr()
		fmt.Printf("Ошибка определения кодировки: %s : Используется : ", err)
		ColorBlue.Print(" utf-8/n")
		return unicode.UTF8.NewDecoder()
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for {
		var lines []string

		for i := 0; i < 50 && scanner.Scan(); i++ {
			lines = append(lines, scanner.Text())
		}

		if len(lines) == 0 {
			PrintWarn()
			fmt.Print("Недостаточно строк для определения кодировки : Используется : ")
			ColorBlue.Print("utf-8\n")
			decoder = unicode.UTF8.NewDecoder()
			break
		}

		result, err := detector.DetectBest([]byte(strings.Join(lines, "")))
		if err != nil {
			PrintErr()
			fmt.Printf("Ошибка определения кодировки: %s : Используется ", err)
			ColorBlue.Print("utf-8\n")
			decoder = unicode.UTF8.NewDecoder()
			break
		}

		if result.Confidence >= 90 {
			detectedEncoding, _ = charset.Lookup(result.Charset)
			decoder = detectedEncoding.NewDecoder()
			PrintSuccess()
			fmt.Print("Определена кодировка : ")
			ColorBlue.Print(result.Charset)
			fmt.Printf(" : Вероятность : ")
			ColorBlue.Print(result.Confidence)
			fmt.Print(" %\n")
			break
		}
	}
	return decoder
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

func GetCurrentFileSize(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	currentFileSize = info.Size()
	cfl := int64(math.Round(float64(currentFileSize) / 80))
	if cfl == 0 {
		currentFileLines = 2
	} else {
		currentFileLines = cfl
	}

	return nil
}

func CreatePBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(
		int(currentFileLines),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetItsString("Str"),
		progressbar.OptionSetRenderBlankState(true),
	)
}

func RemoveDublesResultFiles() {

	ExistFileCount := CheckFileExists()

	if ExistFileCount == 0 {
		PrintErr()
		fmt.Print("Нет файлов для удаления дублей\n")
		return
	}

	for _, request := range searchRequests {
		dublesWG.Add(1)

		if requestStructMap[request].resultFileExist {

			err := dublesPool.Submit(func() {
				DublesRemove(request, requestStructMap[request].resultFile)
			})

			if err != nil {
				PrintErr()
				fmt.Print("Ошибка удаления дублей из полученных файлов : ", err)
			}
		}
	}

	dublesWG.Wait()

}

func DublesRemove(request string, path string) {
	defer dublesWG.Done()
	var lines []string
	rdfile, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	scanner := bufio.NewScanner(rdfile)
	scanner.Split(bufio.ScanLines)

	if err != nil {
		PrintErr()
		fmt.Printf("%s : Ошибка удаления дубликатов : %s\n", path, err)
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
		PrintErr()
		fmt.Printf("%s : Ошибка удаления дубликатов : %s\n", path, err)
		return
	}

	_ = wrfile.Truncate(0)
	_, _ = wrfile.Seek(0, 0)

	_, err = bufio.NewWriter(wrfile).WriteString(strings.Join(lines, "\n") + "\n")
	if err != nil {
		PrintErr()
		fmt.Printf("%s : Ошибка удаления дубликатов : %s\n", path, err)
		return
	}
	wrfile.Close()

	PrintSuccess()
	_, _ = ColorBlue.Print(request)
	fmt.Print(" : ")

	fullPath, err := filepath.Abs(path)
	if err == nil {
		fmt.Print(fullPath, "\n")
	}

	fmt.Print("Всего записано уникальных строк : ")
	ColorBlue.Print(len(lines))
	fmt.Print(" : Всего удалено дубликатов : ")
	ColorBlue.Print(oldLen-len(lines), "\n\n")
}

func CheckFileExists() int {

	var ExistFileCount int

	for _, request := range searchRequests {

		path := requestStructMap[request].resultFile

		_, fileerr := os.Stat(path)

		if fileerr != nil {
			requestStructMap[request].resultFileExist = true
			ExistFileCount++
		}
	}
	return ExistFileCount
}
