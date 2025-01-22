import type { Deps, ExtractedTasks, Field, TasksDuck } from './types'

import { hrtime } from './utils'
import { validateDeps, validateRunningTasks, validateTasks } from './validate'

export function defineTasks<Tasks extends TasksDuck>(
  deps: Deps<Tasks>,
  tasks: Tasks,
) {
  validateTasks<Tasks>(tasks)
  const fieldMap = Object.fromEntries(
    Object.keys(tasks).map((field) => [field, true]),
  ) as Record<Field<Tasks>, true>
  validateDeps<Tasks>(fieldMap, deps)

  type Run = <TaskName extends Field<Tasks>>(
    fields: TaskName[],
    params?: Record<string, any> & { hrtime?: boolean },
  ) => Promise<Pick<ExtractedTasks<Tasks>, TaskName>>
  const run: Run = async (fields, params = {}) => {
    validateRunningTasks<Field<Tasks>>(fieldMap, fields)
    const totalRes = {} as ExtractedTasks<Tasks>
    const promiseMap = new Map<Field<Tasks>, Promise<any>>()
    const timeRecord = {} as Record<Field<Tasks>, number>

    const dealDeps = async (field: Field<Tasks>) => {
      const dependencies = deps[field] || []
      await parallelRun(dependencies)
      return (dependencies as Field<Tasks>[]).reduce(
        (total, curr) => {
          total[curr] = totalRes[curr]
          return total
        },
        {} as typeof totalRes,
      )
    }

    const serialRun = async (field: Field<Tasks>) => {
      const pickedDeps = await dealDeps(field)
      return await tasks[field]({ pickedDeps, ...params })
    }

    const doTask = async (field: Field<Tasks>) => {
      if (totalRes[field]) {
        return
      }
      if (promiseMap.has(field)) {
        totalRes[field] = await promiseMap.get(field)
        return
      }
      const start = params.hrtime ? hrtime() : undefined
      const currPromise = serialRun(field)
      promiseMap.set(field, currPromise)
      totalRes[field] = await currPromise
      if (params.hrtime && !(field in timeRecord)) {
        timeRecord[field] = hrtime(start)
      }
    }

    const parallelRun = async (fields: readonly Field<Tasks>[] = []) => {
      if (fields.length === 0) {
        return
      }
      const promises = fields.map((field) => doTask(field))
      await Promise.allSettled(promises)
    }

    await parallelRun(fields)

    const result = fields.reduce(
      (total, curr) => {
        total[curr] = totalRes[curr]
        return total
      },
      {} as typeof totalRes,
    )

    params.hrtime && Object.assign(result, { _timeRecord: timeRecord })

    return result
  }

  return {
    run,
  }
}
