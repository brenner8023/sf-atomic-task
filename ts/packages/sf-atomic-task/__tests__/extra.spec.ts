import { beforeEach, describe, expect, test, vi } from 'vitest'
import { defineTasks } from '../index'

describe('extra features', () => {
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
      C: ['A'],
    } as const
    const tasks = {
      A: taskA,
      B: taskB,
      C: taskC,
    }
    run = defineTasks(deps, tasks).run
  })

  test('param hrtime should be work', async () => {
    const result = await run(['C'], { hrtime: true })
    expect(taskA).toBeCalledTimes(1)
    expect(taskC).toBeCalledTimes(1)
    expect(result.C).toBe('DataC')
    expect(typeof result._timeRecord.A).toBe('number')
    expect(result._timeRecord).not.toHaveProperty('B')
    expect(typeof result._timeRecord.C).toBe('number')
  })
})
