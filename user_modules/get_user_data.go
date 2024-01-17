package user_modules

func GetUserInputData(appVersion string) (string, []string, []string, string, int) {

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
		case "cleaner":
			searchRequests = nil
			saveType = ""
		case "replacer":
			numParts = GetPartsInput()
			searchRequests = nil
			saveType = ""
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

	return workMode, filePathList, searchRequests, saveType, numParts
}
