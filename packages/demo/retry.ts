
// 重试功能封装示例
export function retryWrapper<RequestFn extends (...args: any) => any>(requestFn: RequestFn, retryTimes = 3) {
  return async function (...args): Promise<Awaited<ReturnType<RequestFn>> | undefined> {
    let count = 0
    while (count < retryTimes) {
      try {
        console.log('count:', count+1)
        return await requestFn(...args)
      } catch (e) {
        count++
        if (count === retryTimes) {
          throw e
        }
      }
    }
  }
}
