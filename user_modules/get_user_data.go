package user_modules

func GetUserInputData(appVersion string) ([]string, []string, string) {

	PrintLogoAnimation(appVersion)
	PrintAutorLinks()
	PrintUserMachineSpecs()

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
