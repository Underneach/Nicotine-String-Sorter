package user_modules

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetFilesInput() []string {

	var result []string

	for true {
		PrintInput()
		fmt.Print("Введите путь к файлу или папке для сортировки: ")

		rawPath, _ := userInputReader.ReadString('\n')
		rawPath = filepath.Clean(strings.TrimSpace(rawPath))

		if rawPath == "" {
			continue
		}

		if fileInfo, fierr := os.Stat(rawPath); fierr == nil {

			if fileInfo.IsDir() {
				PrintSuccess()
				fmt.Printf("Папка '")
				ColorBlue.Print(rawPath)
				fmt.Print("' существует:\n\n")

				_ = filepath.Walk(rawPath, func(path string, info os.FileInfo, fwerr error) error {

					if fwerr != nil {
						PrintErr()
						fmt.Print(fwerr, "\n")
						return fwerr
					}

					if !info.IsDir() {
						if filepath.Ext(path) == ".txt" {
							fmt.Printf("    %s\n", path)
							result = append(result, path)
						}
					}
					return nil
				})
				fmt.Print("\n")
				break

			} else {
				PrintSuccess()
				fmt.Print("Файл со строками найден\n\n")
				result = append(result, rawPath)
				break
			}

		} else {
			PrintErr()
			fmt.Printf("Путь '%s' не существует\n", rawPath)
			continue
		}
	}

	result = Unique(result)
	GetFilesSize(result)
	return result
}

func GetRequestsInput() []string {

	var result []string

	PrintInfo()
	fmt.Print("Поддерживаемые типы ввода:\n\n")
	ColorBlue.Print("       1")
	fmt.Print(" - Ввод из терминала\n")
	ColorBlue.Print("       2")
	fmt.Print(" - Ввод из файла\n\n")
	for true {

		PrintInput()
		fmt.Print("Выберите ввод запросов: ")

		inputType, _ := userInputReader.ReadString('\n')

		switch strings.TrimSpace(inputType) {
		case "1":
			for true {
				PrintInput()
				fmt.Print("Введите запросы через пробел: ")
				rawRequests, _ := userInputReader.ReadString('\n')
				for _, request := range strings.Split(rawRequests, " ") {
					request = strings.TrimSpace(strings.ToLower(request))
					_, err := regexp.Compile(".*" + request + ".*:(.+:.+)")
					if err != nil {
						PrintErr()
						fmt.Printf("%s : Ошибка создания регулярного выражения : %s\n", request, err)
						continue
					}
					result = append(result, request)
				}

				if len(result) == 0 {
					PrintErr()
					fmt.Print("Нет запросов для поиска\n")
					continue
				}
				fmt.Print("\n")
				break
			}
		case "2":
			for true {
				PrintInput()
				fmt.Print("Введите путь к файлу без пробелов: ")
				rawRequests, _ := userInputReader.ReadString('\n')
				rawRequests = strings.TrimSpace(rawRequests)
				_, sterr := os.Stat(rawRequests)
				if sterr != nil {
					PrintErr()
					fmt.Print("Файл не существует\n")
					continue
				}
				file, operr := os.Open(rawRequests)
				if operr != nil {
					PrintErr()
					fmt.Printf("Ошибка чтения файла с запросами : %s\n", operr)
					fmt.Println(operr)
					continue
				}

				defer file.Close()

				scanner := bufio.NewScanner(file)
				scanner.Split(bufio.ScanLines)

				for scanner.Scan() {
					request := strings.TrimSpace(strings.ToLower(scanner.Text()))
					_, err := regexp.Compile(regexp.QuoteMeta(request) + ".*:.+:.+")
					if err != nil {
						PrintErr()
						fmt.Printf("%s : Ошибка создания регулярного выражения : %s\n", request, err)
						continue
					}
					result = append(result, request)
				}

				PrintSuccess()
				fmt.Print("Файл с запросами найден : ")
				ColorBlue.Print(len(result))
				fmt.Print(" запросов\n")

				if len(result) == 0 {
					PrintErr()
					fmt.Print("Нет запросов для поиска\n")
					continue
				}
				fmt.Print("\n")
				break
			}
		default:
			continue
		}
		break
	}
	return Unique(result)
}

func GetSaveTypeInput() string {

	var result string

	PrintInfo()
	fmt.Print("Поддерживаемые типы сохранения:\n\n")
	ColorBlue.Print("       1")
	fmt.Print(" - log:pass (")
	ColorBlue.Print("по умолчанию")
	fmt.Print(")\n")
	ColorBlue.Print("       2")
	fmt.Print(" - url:log:pass\n\n")
	for true {
		PrintInput()
		fmt.Print("Выберите тип сохранения: ")
		rawSaveType, _ := userInputReader.ReadString('\n')
		rawSaveType = strings.TrimSpace(rawSaveType)

		if rawSaveType == "1" || rawSaveType == "2" {
			result = rawSaveType
			fmt.Print("\n")
			break
		}
	}

	return result
}
