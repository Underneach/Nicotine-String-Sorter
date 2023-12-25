package main

import (
	"String-Sorter/user_modules"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CheckUpdate(appVersion string) {

	type GitResponse struct {
		TagName    string `json:"tag_name"`
		ReleaseUrl string `json:"html_url"`
	}
	var apiResponse GitResponse

	httpClient := &http.Client{Timeout: 10 * time.Second}

	resp, err := httpClient.Get("https://api.github.com/repos/Underneach/String-Sorter/releases/latest")
	if err != nil || resp.StatusCode != 200 {
		return
	}

	_ = resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return
	}

	if apiResponse.TagName != appVersion {
		user_modules.PrintInfo()
		fmt.Print("Доступна новая версия : ")
		_, _ = user_modules.ColorBlue.Print(apiResponse.TagName)
		fmt.Print(" : ")
		_, _ = user_modules.ColorBlue.Print(apiResponse.ReleaseUrl, "\n")
	}

}
