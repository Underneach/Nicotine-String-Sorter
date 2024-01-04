package user_modules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CheckUpdate(appVersion string) {

	defer updateWG.Done()

	time.Sleep(time.Second)

	type GitResponse struct {
		TagName    string `json:"tag_name"`
		ReleaseUrl string `json:"html_url"`
	}
	var apiResponse GitResponse

	httpClient := &http.Client{Timeout: 5 * time.Second}

	resp, err := httpClient.Get("https://api.github.com/repos/Underneach/Nicotine-String-Sorter/releases/latest")
	if err != nil || resp.StatusCode != 200 {
		WaitLogo()
		PrintErr()
		fmt.Print("Не удалось проверить обновления : Код HTTP ", resp.StatusCode, "\n\n")
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		WaitLogo()
		PrintErr()
		fmt.Print("Не удалось проверить обновления : ", err, "\n\n")
		return
	}

	if apiResponse.TagName != appVersion {
		WaitLogo()
		PrintInfo()
		fmt.Print("Доступна новая версия : ")
		ColorBlue.Print(apiResponse.TagName)
		fmt.Print(" : ")
		ColorBlue.Print(apiResponse.ReleaseUrl, "\n\n")
	} else {
		WaitLogo()
		PrintSuccess()
		fmt.Print("У вас последняя версия сортера\n\n")
	}
	_ = resp.Body.Close()
}

func WaitLogo() {
	for isLogoPrinted == false {
		time.Sleep(time.Millisecond * 100)
	}
}
