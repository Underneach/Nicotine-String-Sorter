package work_modules

import (
	"bufio"
	"fmt"
	"github.com/pbnjay/memory"
	"github.com/saintfish/chardet"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"math"
	"os"
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
		_, _ = ColorBlue.Print(" utf-8/n")
		detectedEncoding, _ = charset.Lookup("utf-8")
		return detectedEncoding.NewDecoder() // заменить на прямую ссылку с мапы
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
			_, _ = ColorBlue.Print("utf-8\n")
			detectedEncoding, _ = charset.Lookup("utf-8")
			decoder = detectedEncoding.NewDecoder()
			break
		}

		result, err := detector.DetectBest([]byte(strings.Join(lines, "")))
		if err != nil {
			PrintErr()
			fmt.Printf("Ошибка определения кодировки: %s : Используется ", err)
			_, _ = ColorBlue.Print("utf-8\n")
			detectedEncoding, _ = charset.Lookup("utf-8")
			decoder = detectedEncoding.NewDecoder()
			break
		}

		if result.Confidence >= 90 {
			detectedEncoding, _ = charset.Lookup(result.Charset)
			decoder = detectedEncoding.NewDecoder()
			PrintSuccess()
			fmt.Print("Определена кодировка : ")
			_, _ = ColorBlue.Print(result.Charset)
			fmt.Printf(" : Вероятность : ")
			_, _ = ColorBlue.Print(result.Confidence, "\n")
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

func CreateBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(
		int(currentFileLines),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetItsString("Str"),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetRenderBlankState(true),
	)
}
