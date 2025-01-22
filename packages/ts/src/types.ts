export type TasksDuck = Record<string, (...args: any[]) => any>

export type Field<Tasks> = string & keyof Tasks

export type Deps<Tasks extends TasksDuck> = {
  readonly [k in Field<Tasks>]?: readonly Field<Tasks>[]
}

type UnwrapPromise<T> = T extends Promise<infer U> ? U : T

export type ExtractedTasks<T extends Record<string, () => any>> = {
  [K in keyof T]: UnwrapPromise<ReturnType<T[K]>>
}

type DepsDuck<TasksRes> = Readonly<
  Record<keyof TasksRes, readonly (keyof TasksRes)[]>
>
type TaskArg<
  TaskDeps extends DepsDuck<TasksRes>,
  TasksRes,
  K extends keyof TaskDeps,
  Params,
> = {
  pickedDeps: Pick<TasksRes, TaskDeps[K][number]>
} & { [Key in keyof Params]: Params[Key] }

export type AtomicTasks<
  TaskDeps extends DepsDuck<TasksRes>,
  TasksRes extends Record<string, any>,
  Params extends Record<string, any> = Record<string, any>,
> = {
  [K in keyof TasksRes]: (
    taskArg: TaskArg<TaskDeps, TasksRes, K, Params>,
  ) => Promise<TasksRes[K]>
}
