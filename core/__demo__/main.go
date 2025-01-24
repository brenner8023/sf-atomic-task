package main

import (
	"fmt"
	"time"

	. "github.com/brenner8023/sf-atomic-task"
)

// 串行链路：productId -> productSize -> productPrice
func serialDemo() {
  deps := map[string][]string{
    "productSize": {"productId"},
    "productPrice": {"productSize"},
  }
  tasks := map[string]TaskFunc{
    "productId": func(depContext map[string]any, params map[string]any) (any, error) {
      return "001", nil
    },
    "productSize": func(depContext map[string]any, params map[string]any) (any, error) {
      return "XXL", nil
    },
    "productPrice": func(depContext map[string]any, params map[string]any) (any, error) {
      return 100, nil
    },
  }
  run, err := DefineTasks(deps, &tasks)
  if err != nil {
    println(err.Error())
    return
  }
  result, err := run([]string{"productPrice"})
  if err != nil {
    println(err.Error())
    return
  }
  fmt.Printf("res: %v\n", result)
}

// 并行链路：productId -> productSize, productRank, productBrand
func parallelDemo() {
  deps := map[string][]string{
    "productSize": {"productId"},
    "productRank": {"productId"},
    "productBrand": {"productId"},
  }
  tasks := map[string]TaskFunc{
    "productId": func(depContext map[string]any, params map[string]any) (any, error) {
      return "001", nil
    },
    "productSize": func(depContext map[string]any, params map[string]any) (any, error) {
      fmt.Printf("productSize-dep: %v\n", depContext)
      return "XXL", nil
    },
    "productRank": func(depContext map[string]any, params map[string]any) (any, error) {
      fmt.Printf("productRank-dep: %v\n", depContext)
      return "NO.1", nil
    },
    "productBrand": func(depContext map[string]any, params map[string]any) (any, error) {
      fmt.Printf("productBrand-dep: %v\n", depContext)
      return "Nike", nil
    },
  }
  run, err := DefineTasks(deps, &tasks)
  if err != nil {
    println(err.Error())
    return
  }
  result, err := run([]string{"productSize", "productRank", "productBrand"})
  if err != nil {
    println(err.Error())
    return
  }
  fmt.Printf("res: %v\n", result)
}

// 业务自定义重试功能
func retryDemo() {
  retryTask := func(task TaskFunc, retryTime int, sleep int) TaskFunc {
    return func(depContext map[string]any, params map[string]any) (any, error) {
      var result any
      var err error
      for i := 0; i < retryTime; i++ {
        result, err = task(depContext, params)
        if err == nil {
          break
        }
        if sleep > 0 {
          time.Sleep(time.Duration(sleep) * time.Second)
        }
      }
      return result, err
    }
  }
  productId := func(depContext map[string]any, params map[string]any) (any, error) {
    return "211", nil
  }
  tasks := map[string]TaskFunc{
    "productId": retryTask(productId, 3, 1),
  }
  run, err := DefineTasks(nil, &tasks)
  if err != nil {
    println(err.Error())
    return
  }
  result, err := run([]string{"productId"})
  if err != nil {
    println(err.Error())
    return
  }
  fmt.Printf("res: %v\n", result)
}

func main() {
  serialDemo()
  parallelDemo()
  retryDemo()
}