package user_modules

func GetUserInputData(appVersion string) ([]string, []string, string) {
	
	SetTermTitle(appVersion)
	updateWG.Add(1)
	go CheckUpdate(appVersion)
	PrintLogoStart(appVersion)
	updateWG.Wait()

LoopInput:
	for true {

		filePathList = GetFilesInput()
		searchRequests = GetRequestsInput()
		saveType = GetSaveTypeInput()

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

	return filePathList, searchRequests, saveType

}
