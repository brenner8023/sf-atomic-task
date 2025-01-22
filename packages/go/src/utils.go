package core

import (
	"fmt"
	"strings"
)

func FormatError(title, msg string) error {
  message := fmt.Sprintf("%s: %s - %s", APP_NAME, title, msg)
  return fmt.Errorf(message)
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

func assignError(errorChan *chan error) error {
  var strBuilder strings.Builder

  for err := range *errorChan {
    strBuilder.WriteString(fmt.Sprintf("%v\n", err))
  }
  // close(*errorChan)
  if strBuilder.Len() > 0 {
    return FormatError("assignError", strBuilder.String())
  }
  return nil
}
