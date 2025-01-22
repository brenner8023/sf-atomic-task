import { beforeEach, describe, expect, test, vi } from 'vitest'
import { defineTasks } from '../index'

describe('parallel process', () => {
  let taskA
  let taskB
  let taskC
  let taskD
  let run

  beforeEach(() => {
    taskA = vi.fn(() => 'DataA')
    taskB = vi.fn(() => Promise.resolve('DataB'))
    taskC = vi.fn(async () => {
      return 'DataC'
    })
    taskD = vi.fn(() => Promise.resolve('DataD'))
    const deps = {
      B: ['A'],
      C: ['A'],
      D: ['A'],
    } as const
    const tasks = {
      A: taskA,
      B: taskB,
      C: taskC,
      D: taskD,
    }
    run = defineTasks(deps, tasks).run
  })

  test('run task A->B|C|D', async () => {
    const result = await run(['B', 'C', 'D'])
    expect(taskA).toBeCalledTimes(1)
    expect(taskB).toBeCalledTimes(1)
    expect(taskC).toBeCalledTimes(1)
    expect(taskD).toBeCalledTimes(1)
    expect(result.B).toBe('DataB')
    expect(result.C).toBe('DataC')
    expect(result.D).toBe('DataD')
    expect(result).not.toHaveProperty('A')
  })

  test('with no deps, run task B|C|D', async () => {
    taskB = vi.fn(() => Promise.resolve('DataB'))
    taskC = vi.fn(async () => {
      return 'DataC'
    })
    taskD = vi.fn(() => Promise.resolve('DataD'))
    const deps = {}
    const tasks = {
      B: taskB,
      C: taskC,
      D: taskD,
    }
    run = defineTasks(deps, tasks).run
    const result = await run(['B', 'C', 'D'])
    expect(taskB).toBeCalledTimes(1)
    expect(taskC).toBeCalledTimes(1)
    expect(taskD).toBeCalledTimes(1)
    expect(result.B).toBe('DataB')
    expect(result.C).toBe('DataC')
    expect(result.D).toBe('DataD')
  })

  // 既有串联又有并联的链路
  test('run task A->B|C->D', async () => {
    taskA = vi.fn(() => 'DataA')
    taskB = vi.fn(({ pickedDeps }) => Promise.resolve(`${pickedDeps.A}->DataB`))
    taskC = vi.fn(async ({ pickedDeps }) => {
      return `${pickedDeps.A}->DataC`
    })
    taskD = vi.fn(({ pickedDeps }) =>
      Promise.resolve(`DataD:${pickedDeps.B}_${pickedDeps.C}`),
    )
    const deps = {
      B: ['A'],
      C: ['A'],
      D: ['B', 'C'],
    } as const
    const tasks = {
      A: taskA,
      B: taskB,
      C: taskC,
      D: taskD,
    }
    run = defineTasks(deps, tasks).run
    const { D } = await run(['D'])
    expect(taskA).toBeCalledTimes(1)
    expect(taskB).toBeCalledTimes(1)
    expect(taskC).toBeCalledTimes(1)
    expect(taskD).toBeCalledTimes(1)
    expect(D).toBe('DataD:DataA->DataB_DataA->DataC')
  })
})
