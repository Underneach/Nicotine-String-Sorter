package work_modules

import "time"

func Run(_workMode string, _filePathList []string, _searchRequests []string, _saveType string, _cleanType string, _delimetr string) time.Duration {
	InitVar(_workMode, _filePathList, _searchRequests, _saveType, _cleanType, _delimetr)

	startTime := time.Now()

	switch workMode {
	case "sorter":
		InitSorter()
		RunSorter()
		for _, filePath := range filePathList {
			Sorter(filePath)
		}
		PrintSorterResult()

	case "cleaner":
		RunCleaner()
		for _, filePath := range filePathList {
			Cleaner(filePath)
		}
		PrintCleanerResult()
	}

	return time.Since(startTime)
}
