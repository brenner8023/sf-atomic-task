package Atom

import (
	"sync"

	core "github.com/brenner8023/sf-atomic-task/core"
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
    var mutex sync.Mutex
    config := core.TaskConfig{
      Params: &paramsMap,
      TaskChanMap: &taskChanMap,
      Mu: &mutex,
      OnceMap: onceMap,
    }
    core.Debug(&config, "fields value", fields)
    core.Debug(&config, "paramsMap value", paramsMap)

    err = core.ParallelRun(fields, &deps, tasks, &config)
    if err != nil {
      return nil, err
    }

    result := make(map[string]any, length)

    for _, field := range fields {
      taskChan := taskChanMap[field]
      result[field] = <-taskChan
    }
    for _, taskChan := range taskChanMap {
      close(taskChan)
    }
    core.Debug(&config, "finished, result value", result)
    return result, nil
  }
  return run, nil
}
