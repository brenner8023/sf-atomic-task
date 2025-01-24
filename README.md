# sf-atomic-task

多任务编排能力封装，声明式定义多个任务的依赖关系，简化数据聚合

- 提供ts、golang两种实现，支持浏览器端、服务端
- 支持串并行复杂链路的声明
- 支持业务自定义功能，比如重试功能
- 生产代码零外部依赖，体积小
- 单测覆盖率90%以上

golang用法
```go
import (
	. "github.com/brenner8023/sf-atomic-task"
)

func main() {
  // 声明不同任务之间的依赖关系
  deps := map[string][]string{
    "productSize": {"productId"},
    "productPrice": {"productSize"},
  }
  // 定义任务的具体实现
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
  // 根据实际场景选择你要执行的任务
  result, err := run([]string{"productPrice"})
}
```

TypeScript用法
```ts
import { defineTasks } from 'sf-atomic-task'

const deps = {
  productSize: ['productId'],
  productPrice: ['productSize'],
} as const
const tasks = {
  productId: () => Promise.resolve('001'),
  productSize: ({ pickedDeps }) => {
    console.log('productSize-pickedDeps', pickedDeps.productId)
    return Promise.resolve('XXL')
  },
  productPrice: ({ pickedDeps }) => {
    console.log('productPrice-pickedDeps', pickedDeps)
    return Promise.resolve(100)
  },
}

const { run } = defineTasks(deps, tasks)
const result = await run(['productPrice'])
```
