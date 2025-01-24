package core

import (
	"errors"
	"fmt"
)

func FormatError(title, msg string) error {
	message := fmt.Sprintf("%s: %s - %s", APP_NAME, title, msg)
	return errors.New(message)
}

func Debug(config *TaskConfig, title string, record ...any) {
	if (*config.Params)["debug"] != true {
		return
	}
	if len(record) == 0 {
		fmt.Printf("%s: %s\n", APP_NAME, title)
		return
	}
	message := fmt.Sprintf("%+v", record)
	fmt.Printf("%s: %s - %s\n", APP_NAME, title, message)
}
