import { defineConfig } from 'tsup'

export default defineConfig((options) => {
  return {
    entry: ['index.ts'],
    outDir: 'dist',
    format: ['cjs', 'esm'],
    shims: false,
    dts: true,
    watch: options.watch,
    minify: false,
    sourcemap: true,
    clean: true,
  }
})
