package user_modules

func GetUserInputData(appVersion string) ([]string, []string, string) {

	updateWG.Add(1)
	go CheckUpdate(appVersion)
	PrintLogoStart(appVersion)
	updateWG.Wait()

	for true {

		filePathList = GetFilesInput()
		searchRequests = GetRequestsInput()
		saveType = GetSaveTypeInput()

		PrintLogoFast(appVersion)
		if PrintInputData(appVersion) == "continue" {
			break
		} else {
			continue
		}
	}

	PrintLogoFast(appVersion)

	return filePathList, searchRequests, saveType

}
