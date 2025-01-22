package core

import "fmt"

func ValidateDeps(fieldMap map[string]bool, deps map[string][]string) error {
  for field, depFields := range deps {
    if !fieldMap[field] {
      message := fmt.Sprintf("deps[%s] is not defined in tasks", field)
      return FormatError("validateDeps", message)
    }
    for _, depItem := range depFields {
      if depItem == field {
        message := fmt.Sprintf("deps[%s] has a circular dependency", field)
        return FormatError("validateDeps", message)
      }
      if !fieldMap[depItem] {
        message := fmt.Sprintf("deps[%s]%s is not defined in tasks", field, depItem)
        return FormatError("validateDeps", message)
      }
    }
  }
  return nil
}

func ValidateRunningTasks(fieldMap map[string]bool, fields []string) error {
  for _, field := range fields {
    if !fieldMap[field] {
      message := fmt.Sprintf("%s is not defined in tasks", field)
      return FormatError("validateRunningTasks", message)
    }
  }
  return nil
}
