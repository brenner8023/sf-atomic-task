import { describe, expect, test, vi } from 'vitest'
import { defineTasks } from '../index'

describe('validate process', () => {
  test('function validateDeps should work', async () => {
    const taskA = vi.fn(() => 'DataA')
    const taskB = vi.fn(() => Promise.resolve('DataB'))
    const deps = {
      A: [],
      C: ['A'],
    } as const
    const tasks = {
      A: taskA,
      B: taskB,
    }
    expect(() => defineTasks(deps, tasks)).toThrowError(
      'sf-atomic-task: validateDeps - deps[C] is not defined in tasks',
    )

    const deps2 = {
      A: ['A'],
    } as const
    const tasks2 = {
      A: taskA,
    }
    expect(() => defineTasks(deps2, tasks2)).toThrowError(
      'sf-atomic-task: validateDeps - deps[A]A must not be the same',
    )

    const deps3 = {
      A: ['C'],
    } as any
    const tasks3 = {
      A: taskA,
    }
    expect(() => defineTasks(deps3, tasks3)).toThrowError(
      'sf-atomic-task: validateDeps - deps[A]C is not defined in tasks',
    )
  })

  test('function validateTasks should work', async () => {
    const tasks = {
      A: 'taskA',
    }
    expect(() => defineTasks({}, tasks as any)).toThrowError(
      'sf-atomic-task: validateTasks - task[A] must be a function',
    )
  })

  test('function validateRunningTasks should work', async () => {
    const taskA = vi.fn(() => 'DataA')
    const taskB = vi.fn(() => Promise.resolve('DataB'))
    const tasks = {
      A: taskA,
      B: taskB,
    }
    const { run } = defineTasks({}, tasks)
    await expect(async () => await run(['D' as any])).rejects.toThrowError(
      'sf-atomic-task: validateRunningTasks - D is not defined in tasks',
    )
  })

  test('when task A is rejected, other tasks should work', async () => {
    const taskA = vi.fn(() => Promise.reject('errorA'))
    const taskB = vi.fn(({ pickedDeps }) =>
      Promise.resolve(`B:${pickedDeps.A}`),
    )
    const taskC = vi.fn(() => Promise.resolve('DataC'))
    const deps = {
      A: [],
      B: ['A'],
      C: [],
    } as const
    const tasks = {
      A: taskA,
      B: taskB,
      C: taskC,
    }
    const { run } = defineTasks(deps, tasks)
    const result = await run(['B', 'C'])
    expect(taskA).toBeCalledTimes(1)
    expect(taskB).toBeCalledTimes(1)
    expect(taskC).toBeCalledTimes(1)
    expect(result.B).toBe('B:undefined')
    expect(result.C).toBe('DataC')
  })
})
