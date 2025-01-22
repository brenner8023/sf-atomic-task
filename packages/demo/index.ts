
import { defineTasks } from 'sf-atomic-task'
import type { AtomicTasks } from 'sf-atomic-task'
import { fetchProductSize } from './api'
import { retryWrapper } from './retry'

// 串行链路：productId -> productSize -> productPrice
async function serialDemo() {
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
  console.log('res:', result.productPrice)
}

// 并行链路：productId -> productSize, productRank, productBrand
async function parallelDemo() {
  const deps = {
    productSize: ['productId'],
    productRank: ['productId'],
    productBrand: ['productId'],
  } as const
  const tasks = {
    productId: () => Promise.resolve('001'),
    productSize: ({ pickedDeps }) => {
      console.log('productSize-pickedDeps', pickedDeps)
      return Promise.resolve('XXL')
    },
    productRank: ({ pickedDeps }) => {
      console.log('productRank-pickedDeps', pickedDeps)
      return Promise.resolve('NO.1')
    },
    productBrand: ({ pickedDeps }) => {
      console.log('productBrand-pickedDeps', pickedDeps)
      return Promise.resolve('Nike')
    },
  }

  const { run } = defineTasks(deps, tasks)

  const result = await run(['productSize', 'productRank', 'productBrand'])
  console.log('res:', result.productSize, result.productRank, result.productBrand)
}

// 业务自定义重试功能
async function retryDemo() {
  const deps = {
    productSize: ['productId'],
  } as const
  const tasks = {
    productId: () => Promise.resolve('001'),
    productSize: retryWrapper(fetchProductSize, 3),
  }
  const { run } = defineTasks(deps, tasks)

  const result = await run(['productSize'])
  console.log(result.productSize)
}

// 类型声明示例
async function typesDemo() {
  interface TasksDeps {
    productId: [],
    productSize: ['productId'],
    productPrice: ['productSize'],
  }

  interface TasksRes {
    productId: string
    productSize: 'S' | 'M' | 'L' | 'XL' | 'XXL'
    productPrice: number
  }

  const deps = {
    productId: [],
    productSize: ['productId'],
    productPrice: ['productSize'],
  } as const
  const tasks: AtomicTasks<TasksDeps, TasksRes> = {
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
  console.log('res:', result.productPrice)
}

// 每个原子任务的耗时统计
async function hrtimeDemo() {
  const deps = {
    productSize: ['productId'],
    productPrice: ['productSize'],
  } as const
  const tasks = {
    productId: () => {
      return new Promise((resolve) => {
        setTimeout(() => {
          resolve('001')
        }, 1000)
      })
    },
    productSize: () => {
      return new Promise((resolve) => {
        setTimeout(() => {
          resolve('L')
        }, 1500)
      })
    },
    productPrice: () => {
      return new Promise((resolve) => {
        setTimeout(() => {
          resolve(100)
        }, 500)
      })
    },
  }

  const { run } = defineTasks(deps, tasks)

  const result = await run(['productPrice'], { hrtime: true })
  console.log('time:', (result as any)._timeRecord)
}

async function main() {
  console.log('\n---serialDemo---')
  await serialDemo()
  console.log('\n---parallelDemo---')
  await parallelDemo()
  console.log('\n---retryDemo---')
  await retryDemo()
  console.log('\n---hrtimeDemo---')
  hrtimeDemo()
}

main()
