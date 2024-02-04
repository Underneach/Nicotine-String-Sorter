package main

import (
	"String-Sorter/user_modules"
	"String-Sorter/work_modules"
)

var appVersion = "2.6.0"

func main() {

	user_modules.PrintTimeDuration(work_modules.Run(user_modules.GetUserInputData(appVersion)))

}
