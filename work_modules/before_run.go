package work_modules

import (
	"fmt"
	"regexp"
)

func BeforeRun() {

	PrintInfo()
	fmt.Print("Запуск сортера...\n\n")

	var compiledRegEx *regexp.Regexp
	var err error

	for _, request = range searchRequests {

		switch saveType {
		case "1":
			compiledRegEx, err = regexp.Compile(".*" + request + ".*/:(.+/:.+)")
		case "2":
			compiledRegEx, err = regexp.Compile("(.*" + request + ".*/:.+/:.+)")
		}

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
