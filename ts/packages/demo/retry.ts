
// 重试功能封装示例
export function retryTask<RequestFn extends (...args: any) => any>(requestFn: RequestFn, { retryTimes = 3, sleep = 0 }) {
  return async function (...args): Promise<Awaited<ReturnType<RequestFn>> | undefined> {
    let count = 0
    while (count < retryTimes) {
      try {
        console.log('count:', count+1)
        return await requestFn(...args)
      } catch (e) {
        const fn = () => {
          count++
          if (count === retryTimes) {
            throw e
          }
        }
        if (sleep) {
          setTimeout(fn, sleep)
        } else {
          fn()
        }
      }
    }
  }
}
