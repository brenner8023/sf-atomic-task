import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    coverage: {
      enabled: !process.env.CI,
      include: ['index.ts', 'src/**/*'],
    },
  },
})
