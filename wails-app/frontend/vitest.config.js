import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    // 测试文件匹配模式
    include: ['src/**/*.test.js'],

    // 测试环境 (纯函数测试不需要 DOM)
    environment: 'node',

    // 覆盖率配置
    coverage: {
      provider: 'v8',
      include: ['src/utils.js'],
      exclude: ['src/main.js'], // main.js 包含 DOM 操作，不在此覆盖
      reporter: ['text', 'html'],
      thresholds: {
        statements: 80,
        branches: 80,
        functions: 80,
        lines: 80,
      },
    },
  },
});
