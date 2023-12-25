package main

import (
	"String-Sorter/user_modules"
	"String-Sorter/work_modules"
	"runtime"
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

	runtime.GOMAXPROCS(min(runtime.NumCPU()-1, 1)) // В Go нет лимита используемых ядер | Вычитаем одно ядро что бы не было зависаний системы | Делаем костыль для мегамозгов с дедом на 1 ядро

	go CheckUpdate(appVersion) // Проверяем обнову в отдельном потоке

	filePathList, searchRequests, saveType := user_modules.GetUserInputData(appVersion)

	startTime = time.Now() // Получаем время начала сортинга
	invalidLines, checkedLines, checkedFiles = work_modules.RunSorter(filePathList, searchRequests, saveType)

	user_modules.PrintResult(time.Since(startTime), checkedLines, invalidLines, checkedLines, len(filePathList), checkedFiles)

}
