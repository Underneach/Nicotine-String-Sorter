package user_modules

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
		case "4":
			os.Exit(0)
		default:
			continue LoopWork
		}
	}
	return work
}

func GetFilesInput() (result []string) {

Loop:
	for {
		PrintInput()
		fmt.Print("Введите путь к файлу или папке для обработки: ")

		rawPath, _ := userInputReader.ReadString('\n')
		rawPath = strings.TrimSpace(rawPath)
		if rawPath == "" {
			continue Loop
		}

		rawPath = filepath.Clean(rawPath)

		if fileInfo, fierr := os.Stat(rawPath); fierr == nil {

			if fileInfo.IsDir() {
				PrintSuccess()
				fmt.Printf("Папка '")
				ColorBlue.Print(rawPath)
				fmt.Print("' найдена:\n\n")

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

				if len(result) >= 1 {
					fmt.Print("\n")
					break Loop
				} else {
					PrintErr()
					fmt.Print("Нет файлов для обработки\n")
					continue Loop
				}

			} else {
				PrintSuccess()
				fmt.Print("Файл со строками найден\n\n")
				result = append(result, rawPath)
				break Loop
			}

		} else {
			PrintErr()
			fmt.Printf("Путь '%s' не существует\n", rawPath)
			continue Loop
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

LoopA:
	for {

		PrintInput()
		fmt.Print("Выберите ввод запросов: ")

		inputType, _ := userInputReader.ReadString('\n')

		switch strings.TrimSpace(inputType) {
		case "1":
		LoopB:
			for true {
				PrintInput()
				fmt.Print("Введите запросы через пробел: ")
				rawRequests, _ := userInputReader.ReadString('\n')
				rawRequests = strings.TrimSpace(rawRequests)
				if rawRequests == "" {
					continue LoopB
				}
				for _, request := range strings.Split(rawRequests, " ") {
					request = strings.TrimSpace(strings.ToLower(request))
					_, err := regexp.Compile(".*" + request + ".*:(.+:.+)")
					if err != nil {
						PrintErr()
						fmt.Printf("%s : Ошибка создания регулярного выражения : %s\n", request, err)
						continue LoopB
					}
					requests = append(requests, request)
				}

				if len(requests) == 0 {
					PrintErr()
					fmt.Print("Нет запросов для поиска\n")
					continue LoopB
				}
				fmt.Print("\n")
				break LoopA
			}
		case "2":
		LoopC:
			for true {
				PrintInput()
				fmt.Print("Введите путь к файлу: ")
				rawRequests, _ := userInputReader.ReadString('\n')
				rawRequests = strings.TrimSpace(rawRequests)
				_, sterr := os.Stat(rawRequests)
				if sterr != nil {
					PrintErr()
					fmt.Print("Файл не существует\n")
					continue LoopC
				}
				file, operr := os.Open(rawRequests)
				if operr != nil {
					PrintErr()
					fmt.Printf("Ошибка чтения файла с запросами : %s\n", operr)
					fmt.Println(operr)
					continue LoopC
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
						continue LoopC
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
					continue LoopA
				}
				fmt.Print("\n")
				break LoopA
			}
		default:
			continue LoopA
		}
	}
	return Unique(requests)
}

func GetSaveTypeInput() (saveType string) {

	PrintInfo()
	fmt.Print("Поддерживаемые типы сохранения:\n\n")
	ColorBlue.Print("       1")
	fmt.Print(" - Log:Pass\n")
	ColorBlue.Print("       2")
	fmt.Print(" - Url:Log:Pass\n\n")

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

func GetCleanTypeInput() (cleanType string) {

	PrintInfo()
	fmt.Print("Поддерживаемые режимы клинера:\n\n")
	ColorBlue.Print("       1")
	fmt.Print(" - Чистка и сохранение каждой базы отдельно\n")
	ColorBlue.Print("       2")
	fmt.Print(" - Чистка всех баз вместе и сохранение в один файл\n\n")

Loop:
	for true {
		PrintInput()
		fmt.Print("Выберите режим чистки: ")
		rawcleanType, _ := userInputReader.ReadString('\n')
		rawcleanType = strings.TrimSpace(rawcleanType)

		switch rawcleanType {
		case "1", "2":
			cleanType = rawcleanType
			fmt.Print("\n")
			break Loop
		default:
			continue Loop
		}
	}
	return cleanType
}

func GetDelimetrInput() (delimetr string) {

LoopDel:
	for true {
		PrintInput()
		fmt.Print("Введите разделитель строк: ")
		var rawDelTrim string

		rawDel, _ := userInputReader.ReadString('\n')

		switch rawDel {
		case "":
			continue LoopDel
		case " ":
			rawDelTrim = rawDel
		default:
			rawDelTrim = strings.TrimSpace(rawDel)
		}

		PrintInfo()
		fmt.Print("Разделитель строк - '")
		ColorBlue.Print(rawDelTrim)
		fmt.Print("'\n\n")
		ColorBlue.Print("       1")
		fmt.Print(" - Продолжить\n")
		ColorBlue.Print("       2")
		fmt.Print(" - Ввести заново\n\n")
	LoopAction:
		for true {
			PrintInput()
			fmt.Print("Выберите действие: ")
			action, _ := userInputReader.ReadString('\n')
			action = strings.TrimSpace(action)
			switch action {
			case "1":
				delimetr = rawDelTrim
				break LoopDel
			case "2":
				continue LoopDel
			default:
				continue LoopAction
			}
		}
	}

	return
}
