package work_modules

import (
	"fmt"
	"regexp"
)

func BeforeRun() {

	PrintInfo()
	fmt.Print("Запуск сортера...")

	var compiledRegEx *regexp.Regexp
	var err error

	for _, request := range searchRequests {

		switch saveType {
		case "1":
			compiledRegEx, err = regexp.Compile(".*" + regexp.QuoteMeta(request) + ".*:(.+:.+)")
		case "2":
			compiledRegEx, err = regexp.Compile("(" + ".*" + regexp.QuoteMeta(request) + ".*:.+:.+)")
		}

		if err != nil {
			PrintErr()
			fmt.Printf("%s : Ошибка компиляции запроса : %s\n", request, err)
			continue
		}

		currentStruct := new(Work)
		currentStruct.requestPattern = compiledRegEx
		currentStruct.resultFile = runDir + `\` + badSymbolsPattern.ReplaceAllString(request, "_") + ".txt"
		requestStructMap[request] = currentStruct
	}

	if len(requestStructMap) == 0 {
		PrintZeroRequestsErr()
	}

	fmt.Print("\r")
	PrintSuccess()
	fmt.Print("Сортер запущен   \n\n")
}
