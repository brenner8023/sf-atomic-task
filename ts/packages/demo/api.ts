
let count = 0

export async function fetchProductSize() {
  if (count < 3) {
    count++
    return Promise.reject('fetchProductSize-error')
  }
  return Promise.resolve('XXL')
}
