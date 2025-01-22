package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 模拟一个任务函数，可能会返回错误
func doTask(ctx context.Context, field string) error {
    select {
    case <-time.After(2 * time.Second):
        if field == "B" {
            return fmt.Errorf("error in task %s", field)
        }
        fmt.Printf("Task %s completed\n", field)
        return nil
    case <-ctx.Done():
        fmt.Printf("Task %s cancelled\n", field)
        return ctx.Err()
    }
}

func main2() {
    fmt.Println("start")
    fields := []string{"A", "B", "C"}

    // 创建一个带取消功能的上下文
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var wg sync.WaitGroup
    errChan := make(chan error, len(fields))

    for _, field := range fields {
        wg.Add(1)
        go func(field string) {
            defer wg.Done()
            if err := doTask(ctx, field); err != nil {
                errChan <- err
                cancel() // 取消其他任务
            }
        }(field)
    }

    // 等待所有任务完成
    go func() {
      println("gogo")
        wg.Wait()
        close(errChan)
    }()

    println("test")

    // 处理错误
    for err := range errChan {
      println("gg")
        if err != nil {
            fmt.Printf("Received error: %v\n", err)
            break
        }
    }

    fmt.Println("All tasks completed or cancelled")
}
