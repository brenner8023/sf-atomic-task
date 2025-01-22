import { beforeEach, describe, expect, test, vi } from 'vitest'
import { defineTasks } from '../index'

describe('serial process A->B->C', () => {
  let taskA
  let taskB
  let taskC
  let run

  beforeEach(() => {
    taskA = vi.fn(() => 'DataA')
    taskB = vi.fn(() => Promise.resolve('DataB'))
    taskC = vi.fn(async () => {
      return 'DataC'
    })
    const deps = {
      A: [],
      B: ['A'],
      C: ['B'],
    } as const
    const tasks = {
      A: taskA,
      B: taskB,
      C: taskC,
    }
    run = defineTasks(deps, tasks).run
  })

  test('when run C task, task ABC should be called', async () => {
    const result = await run(['C'])
    expect(taskA).toBeCalledTimes(1)
    expect(taskB).toBeCalledTimes(1)
    expect(taskC).toBeCalledTimes(1)
    expect(result.C).toBe('DataC')
    expect(result).not.toHaveProperty('B')
    expect(result).not.toHaveProperty('A')
  })

  test('when run B and C task, task ABC should be called', async () => {
    const result = await run(['B', 'C'])
    expect(taskA).toBeCalledTimes(1)
    expect(taskB).toBeCalledTimes(1)
    expect(taskC).toBeCalledTimes(1)
    expect(result.B).toBe('DataB')
    expect(result.C).toBe('DataC')
    expect(result).not.toHaveProperty('A')
  })

  test('when C depends on A and B, C can get result of AB', async () => {
    taskA = vi.fn(() => Promise.resolve('DataA'))
    taskB = vi.fn(() => 'DataB')
    taskC = vi.fn(async ({ pickedDeps }) => {
      return `C:${pickedDeps.A}_${pickedDeps.B}`
    })
    const deps = {
      C: ['A', 'B'],
    } as const
    const tasks = {
      A: taskA,
      B: taskB,
      C: taskC,
    }
    const { run } = defineTasks(deps, tasks)
    const { C } = await run(['C'])
    expect(C).toBe('C:DataA_DataB')
  })
})
