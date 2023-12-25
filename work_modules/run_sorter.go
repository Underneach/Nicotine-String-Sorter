package work_modules

func RunSorter(_filePathList []string, _searchRequests []string, _saveType string) (int64, int64, int) {
	InitVar(_filePathList, _searchRequests, _saveType)
	BeforeRun()

	for _, filePath := range filePathList {
		Sorter(filePath)
	}

	RemoveDublesResultFiles()

	return invalidLines, checkedLines, checkedFiles
}
