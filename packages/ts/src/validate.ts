import type { Deps, Field, TasksDuck } from './types'

import { throwError } from './utils'

export function validateDeps<Tasks extends TasksDuck>(
  fieldMap: Record<Field<Tasks>, true>,
  deps: Deps<Tasks>,
) {
  Object.keys(deps).forEach((field) => {
    if (!fieldMap[field]) {
      throwError({
        title: 'validateDeps',
        msg: `deps[${field}] is not defined in tasks`,
      })
    }
    deps[field]?.some((depItem) => {
      if (depItem === field) {
        throwError({
          title: 'validateDeps',
          msg: `deps[${field}]${depItem} must not be the same`,
        })
      }
      if (!fieldMap[depItem]) {
        throwError({
          title: 'validateDeps',
          msg: `deps[${field}]${depItem} is not defined in tasks`,
        })
      }
      return true
    })
  })
}

export function validateTasks<Tasks extends TasksDuck>(tasks: Tasks) {
  Object.keys(tasks).forEach((field) => {
    const task = tasks[field]
    if (typeof task !== 'function') {
      throwError({
        title: 'validateTasks',
        msg: `task[${field}] must be a function`,
      })
    }
  })
}

export function validateRunningTasks<Field extends keyof any>(
  fieldMap: Record<Field, true>,
  fields: Field[],
) {
  fields.forEach((field) => {
    if (!fieldMap[field]) {
      throwError({
        title: 'validateRunningTasks',
        msg: `${field as string} is not defined in tasks`,
      })
    }
  })
}
