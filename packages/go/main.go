package main

import (
	"fmt"
	core "sf-atomic-task/src"
	"time"
)

type TaskFunc = core.TaskFunc
type runFunc = func(fields []string, params ...map[string]any) (map[string]any, error)

func DefineTasks(deps map[string][]string, tasks *map[string]TaskFunc) (runFunc, error) {
  fieldMap, onceMap, err := core.InitConfig(tasks)
  if err != nil {
    return nil, err
  }
  validateError := core.ValidateDeps(fieldMap, deps)
  if validateError != nil {
    return nil, validateError
  }

  run := func(fields []string, params ...map[string]any) (map[string]any, error) {
    var paramsMap map[string]any
    if len(params) > 1 {
      err := core.FormatError("run", "too many parameters, expected at most 1")
      return nil, err
    } else if len(params) == 0 {
      paramsMap = map[string]any{"debug": false}
    } else {
      paramsMap = params[0]
    }

    validateError := core.ValidateRunningTasks(fieldMap, fields)
    if validateError != nil {
      return nil, validateError
    }

    length := len(fields)
    taskChanMap := make(map[string]chan any, length)

    config := core.TaskConfig{
      Params: &paramsMap,
      TaskChanMap: &taskChanMap,
      OnceMap: onceMap,
    }
    core.Debug(&config, "fields value", fields)
    core.Debug(&config, "paramsMap value", paramsMap)

    core.ParallelRun(fields, &deps, tasks, &config)

    result := make(map[string]any, length)

    for field, taskChan := range taskChanMap {
      result[field] = <-taskChan
      close(taskChan)
    }
    core.Debug(&config, "finished, result value", result)
    return result, nil
  }
  return run, nil
}

func main() {

  deps := map[string][]string{
    "A": {},
    "B": {"A"},
    "C": {"A"},
  }
  tasks := map[string]TaskFunc{
    "A": func(depContext map[string]any, params map[string]any) (any, error) {
      time.Sleep(2 * time.Second)
      return 8023, nil
    },
    "B": func(context map[string]any, params map[string]any) (any, error) {
      return "B", nil
    },
    "C": func(context map[string]any, params map[string]any) (any, error) {
      return true, nil
    },
  }
  run, err := DefineTasks(deps, &tasks)
  if err != nil {
    println(err.Error())
    return
  }
  result, err := run([]string{"A", "B", "C"}, map[string]any{"debug": true})
  if err != nil {
    println(err.Error())
    return
  }
  if false {
    fmt.Printf("result: %v\n", result)
  }
}
