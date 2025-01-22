package core

import (
	"fmt"
	"sync"
)

type TaskFunc func(context map[string]any, params map[string]any) (any, error)

type TaskConfig struct {
  Params *map[string]any
  TaskChanMap *map[string]chan any
  OnceMap map[string]*(sync.Once)
}

func dealDeps(field string, deps *map[string][]string, tasks *map[string]TaskFunc, config *TaskConfig) (map[string]any, error) {
  currDeps, ok := (*deps)[field]
  if !ok {
    return nil, nil
  }
  if len(currDeps) == 0 {
    return nil, nil
  }

  Debug(config, "parallel run deps", field, currDeps)
  err := ParallelRun(currDeps, deps, tasks, config)
  if err != nil {
    return nil, err
  }

  depContext := make(map[string]any, len(currDeps))
  taskChanMap := (*config).TaskChanMap

  for _, depField := range currDeps {
    Debug(config, "read data start", depField)
    data := <-(*taskChanMap)[depField]
    depContext[depField] = data
    Debug(config, "read data end", depField, data)

    Debug(config, "dealDeps: write data start", depField, data)
    (*taskChanMap)[depField] <- data
    Debug(config, "dealDeps: write data end", depField, data)
  }
  return depContext, nil
}

func doTask(field string, deps *map[string][]string, tasks *map[string]TaskFunc, config *TaskConfig) error {
  params := (*config).Params
  taskChanMap := (*config).TaskChanMap

  depContext, err := dealDeps(field, deps, tasks, config)
  if err != nil {
    return err
  }

  currTask, ok := (*tasks)[field]
  if !ok {
    msg := fmt.Sprintf("task[%s] is not defined", field)
    return FormatError("doTask", msg)
  }

  taskData, err := currTask(depContext, *params)
  if err != nil {
    return err
  }

  _, isChanInit := (*taskChanMap)[field]
  if !isChanInit {
    Debug(config, "doTask: write data start", field, taskData)
    (*taskChanMap)[field] = make(chan any, 1)
    (*taskChanMap)[field] <- taskData
    Debug(config, "doTask: write data end", field, taskData)
  }
  return nil
}

func ParallelRun(fields []string, deps *map[string][]string, tasks *map[string]TaskFunc, config *TaskConfig) error {
  total := len(fields)
  if total == 0 {
    return nil
  }

  var waitGroup sync.WaitGroup
  waitGroup.Add(total)

  for _, field := range fields {
    go func(field string) {
      defer waitGroup.Done()
      currOnce := config.OnceMap[field]
      currOnce.Do(func() {
        Debug(config, fmt.Sprintf("---doTask %s start---", field))
        doTask(field, deps, tasks, config)
        Debug(config, fmt.Sprintf("---doTask %s end---", field))
      })
    }(field)
  }
  waitGroup.Wait()

  return nil
}

func InitConfig(tasks *map[string]TaskFunc) (map[string]bool, map[string]*sync.Once, error) {
  total := len(*tasks)
  fieldMap := make(map[string]bool, total)
  onceMap := make(map[string]*sync.Once, total)
  for field := range *tasks {
    fieldMap[field] = true
    onceMap[field] = &sync.Once{}
  }
  return fieldMap, onceMap, nil
}
