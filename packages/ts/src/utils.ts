import { APP_NAME } from './const'

export function throwError({ title, msg }) {
  throw new Error(`${APP_NAME}: ${title} - ${msg}`)
}

export function hrtime(start = 0) {
  return Date.now() - start
}
