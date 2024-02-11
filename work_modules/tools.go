package work_modules

import (
	"bufio"
	"fmt"
	"github.com/pbnjay/memory"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetAviableStringsCount() int64 {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if time.Since(lastUpdate) > time.Second*15 { // Если прошло более 1\4 минуты с момента последнего обновления, обновляем кеш
		cachedStrCount = getAviableStringsCountCached()
		lastUpdate = time.Now()
	}
	return cachedStrCount
}

func getAviableStringsCountCached() int64 {
	freeMemory := memory.FreeMemory()
	if freeMemory != 0 {
		return int64(math.Round(float64(freeMemory / (80 * 4 * 0.90)))) // 80 - Предпологаемый размер строки, 4 - размер символа в байтах, 0.90 - оставляем часть памяти для других элементов сортера
	} else {
		PrintWarn()
		fmt.Print(" Не удалось получить количество доступной памяти : Чтение по чанкам в 2Гб")
		return 6700000 // Возвращаем ~2 гига, если не получили доступный размер
	}
}

func GetFileProcessInfo(path string) *encoding.Decoder {

	result := make(chan *encoding.Decoder, 1)

	go func() {
		result <- GetFileDecoder(path)
	}()

	PrintChunk()

	select {
	case <-time.After(5 * time.Second):
		PrintErr()
		fmt.Print(" Таймаут определения кодировки : Используется ")
		ColorBlue.Print(" UTF-8\n")
		return unicode.UTF8.NewDecoder()
	case result := <-result:
		return result
	}
}

func GetFileDecoder(path string) *encoding.Decoder {

	var detectedEncoding encoding.Encoding
	var decoder *encoding.Decoder
	var lines []string

	detector := chardet.NewTextDetector()
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		PrintErr()
		fmt.Printf(" Ошибка определения кодировки : %s : Используется : ", err)
		ColorBlue.Print(" UTF-8\n")
		return unicode.UTF8.NewDecoder()
	}

	defer file.Close()

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
	lines = nil
	return decoder
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

func GetRunDir() (rundir string) {
	var path string

	if dir, cerr := os.Executable(); cerr != nil {
		path = "."
	} else {
		path = filepath.Dir(dir)
	}

	if _, aerr := os.Stat(path + `\result\`); os.IsNotExist(aerr) {
		if verr := os.Mkdir(path+`\result\`, os.ModePerm); verr == nil {
			rundir = path + `\result\`
		} else {
			rundir = path + `\`
		}
	} else {
		rundir = path + `\result\`
	}

	return rundir
}

func RemoveFromSliceByValue(list []string, item string) (result []string) {
	var success = false
	for i, other := range list {
		if other == item {
			result = append(list[:i], list[i+1:]...)
			success = true
		}
	}
	if success {
		return result
	} else {
		return list
	}
}

func IsDirEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}
