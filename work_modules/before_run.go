package work_modules

import (
	"fmt"
	"regexp"
)

func BeforeRun() {

	PrintInfo()
	fmt.Print("Запуск сортера...\n\n")

	for _, request = range searchRequests {

		compiledRegEx, err := regexp.Compile("(.*" + request + ".*):(.+):(.+)")

		if err != nil {
			PrintErr()
			fmt.Printf("%s : Ошибка компиляции запроса : %s\n", request, err)
			continue
		}

		currentStruct := new(Work)
		currentStruct.requestPattern = compiledRegEx
		requestStructMap[request] = currentStruct
	}

	if len(requestStructMap) == 0 {
		PrintZeroRequestsErr()
	}
}
