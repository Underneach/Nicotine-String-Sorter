package work_modules

import (
	"bufio"
	"fmt"
	"github.com/pbnjay/memory"
	"github.com/saintfish/chardet"
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
	var lines []string

	detector := chardet.NewTextDetector()
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		PrintErr()
		fmt.Printf("Ошибка определения кодировки: %s : Используется : ", err)
		ColorBlue.Print(" UTF-8/n")
		return unicode.UTF8.NewDecoder()
	}

	defer file.Close()
	defer clear(lines)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for {

		for i := 0; i < 50 && scanner.Scan(); i++ {
			lines = append(lines, scanner.Text())
		}

		if len(lines) == 0 {
			PrintEndodingLinesEnd()
			decoder = unicode.UTF8.NewDecoder()
			break
		}

		if result, err := detector.DetectBest([]byte(strings.Join(lines, ""))); err != nil {
			PrintEncodingErr(err)
			decoder = unicode.UTF8.NewDecoder()
			break
		} else if result.Confidence >= 90 {
			detectedEncoding, _ = charset.Lookup(result.Charset)
			decoder = detectedEncoding.NewDecoder()
			PrintEncoding(result)
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

	if cfl := int64(math.Round(float64(currentFileSize) / 80)); cfl == 0 {
		currentFileLines = 10 // да да блять, это же ебучий костыль
	} else {
		currentFileLines = cfl
	}

	return nil
}

func GetRunDir() string {

	var path string

	if dir, err := os.Executable(); err != nil {
		path = "."
	} else {
		path = filepath.Dir(dir)
	}

	return path
}
