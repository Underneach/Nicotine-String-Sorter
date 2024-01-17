package main

import (
	"String-Sorter/user_modules"
	"String-Sorter/work_modules"
	"time"
)

var appVersion = "2.0.0"

func main() {

	user_modules.PrintTimeDuration(time.Since(work_modules.Run(user_modules.GetUserInputData(appVersion))))

}
