package main

import (
	"String-Sorter/user_modules"
	"String-Sorter/work_modules"
	"time"
)

var (
	appVersion = "1.0.0"
	startTime  time.Time
)

func main() {

	filePathList, searchRequests, saveType := user_modules.GetUserInputData(appVersion)

	startTime = time.Now() // Получаем время начала сортинга

	invalidLines, checkedLines, matchLines, checkedFiles := work_modules.RunSorter(filePathList, searchRequests, saveType)

	user_modules.PrintResult(time.Since(startTime), checkedLines, invalidLines, matchLines, len(filePathList), checkedFiles)

}
