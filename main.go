package main

import (
	"String-Sorter/user_modules"
	"String-Sorter/work_modules"
	"time"
)

var (
	checkedLines int64
	invalidLines int64
	checkedFiles int
	appVersion   = "1.0.0"
	startTime    time.Time
)

func main() {

	filePathList, searchRequests, saveType := []string{`C:\Users\truew\GolandProjects\String-Sorter\@urlcloudFREE.txt`}, []string{"1"}, "1"  //user_modules.GetUserInputData(appVersion)

	startTime = time.Now() // Получаем время начала сортинга
	invalidLines, checkedLines, checkedFiles = work_modules.RunSorter(filePathList, searchRequests, saveType)

	user_modules.PrintResult(time.Since(startTime), checkedLines, invalidLines, checkedLines, len(filePathList), checkedFiles)

}
