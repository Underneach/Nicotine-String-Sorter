package user_modules

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func GetWorkMode() (work string) {
	PrintWorkModes()

LoopWork:
	for {
		PrintInput()
		fmt.Print("Выберите тип работы: ")
		wmraw, _ := userInputReader.ReadString('\n')
		wmraw = strings.TrimSpace(wmraw)

		switch wmraw {
		case "1":
			work = "sorter"
			break LoopWork
		case "2":
			work = "cleaner"
			break LoopWork
		/*case "3":
			work = "replacer"
			break LoopWork*/
		case "4":
			os.Exit(0)
		default:
			continue LoopWork
		}
	}
	return work
}

func GetFilesInput() (result []string) {

	for {
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

func GetRequestsInput() (requests []string) {

	PrintInfo()
	fmt.Print("Поддерживаемые типы ввода:\n\n")
	ColorBlue.Print("       1")
	fmt.Print(" - Ввод из терминала\n")
	ColorBlue.Print("       2")
	fmt.Print(" - Ввод из файла\n\n")

	for {

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
					requests = append(requests, request)
				}

				if len(requests) == 0 {
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
				fmt.Print("Введите путь к файлу: ")
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
					requests = append(requests, request)
				}

				PrintSuccess()
				fmt.Print("Файл с запросами найден : ")
				ColorBlue.Print(len(requests))
				fmt.Print(" запросов\n")

				if len(requests) == 0 {
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
	return Unique(requests)
}

func GetSaveTypeInput() (saveType string) {

	PrintInfo()
	fmt.Print("Поддерживаемые типы сохранения:\n\n")
	ColorBlue.Print("       1")
	fmt.Print(" - log:pass\n")
	ColorBlue.Print("       2")
	fmt.Print(" - url:log:pass\n\n")

Loop:
	for true {
		PrintInput()
		fmt.Print("Выберите тип сохранения: ")
		rawSaveType, _ := userInputReader.ReadString('\n')
		rawSaveType = strings.TrimSpace(rawSaveType)

		switch rawSaveType {
		case "1", "2":
			saveType = rawSaveType
			fmt.Print("\n")
			break Loop
		default:
			continue Loop
		}
	}
	return saveType
}

func GetPartsInput() (partNum int) {
	PrintInfo()
	fmt.Print("\n       Расположение частей строки 0:1:2\n\n")
LoopPart:
	for {
		PrintInput()
		fmt.Print("Введите новый порядок расположения частей: ")

		inputPartRaw, _ := userInputReader.ReadString('\n')
		inputPartRaw = strings.TrimSpace(inputPartRaw)
		inputPart, err := strconv.Atoi(inputPartRaw)

		if (partRegex.MatchString(inputPartRaw)) && (err != nil) {
			partNum = inputPart
			break LoopPart
		} else {
			continue LoopPart
		}
	}
	return
}
