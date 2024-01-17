package work_modules

import "time"

func Run(_workMode string, _filePathList []string, _searchRequests []string, _saveType string, _numParts int) (startTime time.Time) {
	InitVar(_workMode, _filePathList, _searchRequests, _saveType, _numParts)

	startTime = time.Now()

	switch workMode {
	case "sorter":
		InitSorter()
		RunSorter()
		for _, filePath := range filePathList {
			Sorter(filePath)
		}
		SorterRemoveDublesResultFiles()
		PrintSorterResult()

	case "cleaner":
		InitCleaner()
		RunCleaner()
		for _, filePath := range filePathList {
			Cleaner(filePath)
		}
		PrintCleanerResult()
	case "replacer":
		InitPartSwitcher()
	}

	return startTime
}
