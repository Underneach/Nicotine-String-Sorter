package main

import (
	"String-Sorter/user_modules"
	"String-Sorter/work_modules"
	"time"
)

var (
	checkedLines int64
	invalidLines int64
	matchLines   int64
	checkedFiles int
	appVersion   = "1.0.0"
	startTime    time.Time
)

func main() {

	filePathList, searchRequests, saveType := []string{`C:\Users\truew\GolandProjects\String-Sorter\test\@urlcloudFREE — копия.txt` /*, `C:\Users\truew\GolandProjects\String-Sorter\test\@urlcloudFREE.txt`*/}, []string{"google", "yandex", "netflix"}, "1" //user_modules.GetUserInputData(appVersion)

	startTime = time.Now() // Получаем время начала сортинга

	invalidLines, checkedLines, matchLines, checkedFiles = work_modules.RunSorter(filePathList, searchRequests, saveType)

	user_modules.PrintResult(time.Since(startTime), checkedLines, invalidLines, matchLines, len(filePathList), checkedFiles)

}
