package user_modules

func GetUserInputData(appVersion string) (string, []string, []string, string, string) {

	updateWG.Add(1)
	go CheckUpdate(appVersion)
	SetTermTitle(appVersion)
	PrintLogoStart(appVersion)
	updateWG.Wait()

LoopInput:
	for true {

		workMode = GetWorkMode()
		filePathList = GetFilesInput()

		switch workMode {
		case "sorter":
			searchRequests = GetRequestsInput()
			saveType = GetSaveTypeInput()
			cleanType = ""
		case "cleaner":
			searchRequests = nil
			saveType = ""
			if len(filePathList) > 1 {
				cleanType = GetCleanTypeInput()
			} else {
				cleanType = "1"
			}
		}

		switch PrintInputData(appVersion) {
		case "continue":
			break LoopInput
		case "restart":
			ClearTerm()
			PrintLogoFast(appVersion)
			continue LoopInput
		}
	}

	PrintLogoFast(appVersion)

	return workMode, filePathList, searchRequests, saveType, cleanType
}
