import { describe, test, expect, beforeEach, vi } from 'vitest';
import {
  containsChinese,
  detectLanguage,
  setLocale,
  getLocale,
  t,
  initI18n,
} from './index.js';

// ========== containsChinese 测试 ==========

describe('containsChinese', () => {
  test('returns true for Chinese text', () => {
    expect(containsChinese('你好')).toBe(true);
    expect(containsChinese('Hello 世界')).toBe(true);
    expect(containsChinese('中文')).toBe(true);
  });

  test('returns false for English only', () => {
    expect(containsChinese('Hello World')).toBe(false);
    expect(containsChinese('test123')).toBe(false);
  });

  test('returns false for empty or null', () => {
    expect(containsChinese('')).toBe(false);
    expect(containsChinese(null)).toBe(false);
    expect(containsChinese(undefined)).toBe(false);
  });

  test('returns false for special characters only', () => {
    expect(containsChinese('!@#$%^&*()')).toBe(false);
  });

  test('returns true for mixed content with Chinese', () => {
    expect(containsChinese('Hello 你好 World')).toBe(true);
    expect(containsChinese('123中文456')).toBe(true);
  });
});

// ========== detectLanguage 测试 ==========

describe('detectLanguage', () => {
  test('returns zh-CN for Chinese text', () => {
    expect(detectLanguage('请优化这个提示词')).toBe('zh-CN');
    expect(detectLanguage('Hello 你好')).toBe('zh-CN');
  });

  test('returns en-US for English text', () => {
    expect(detectLanguage('Please optimize this prompt')).toBe('en-US');
    expect(detectLanguage('Hello World')).toBe('en-US');
  });

  test('returns en-US for empty or null', () => {
    expect(detectLanguage('')).toBe('en-US');
    expect(detectLanguage(null)).toBe('en-US');
    expect(detectLanguage(undefined)).toBe('en-US');
  });
});

// ========== setLocale / getLocale 测试 ==========

describe('setLocale and getLocale', () => {
  beforeEach(() => {
    // Reset to default locale before each test
    setLocale('zh-CN');
  });

  test('sets and gets locale correctly', () => {
    setLocale('en-US');
    expect(getLocale()).toBe('en-US');

    setLocale('zh-CN');
    expect(getLocale()).toBe('zh-CN');
  });

  test('ignores invalid locale', () => {
    setLocale('zh-CN');
    setLocale('invalid-locale');
    expect(getLocale()).toBe('zh-CN');
  });
});

// ========== t (翻译函数) 测试 ==========

describe('t translation function', () => {
  beforeEach(() => {
    setLocale('zh-CN');
  });

  test('returns Chinese translation when locale is zh-CN', () => {
    setLocale('zh-CN');
    expect(t('currentVersion')).toBe('当前版本');
    expect(t('cancel')).toBe('取消');
    expect(t('confirm')).toBe('确定');
  });

  test('returns English translation when locale is en-US', () => {
    setLocale('en-US');
    expect(t('currentVersion')).toBe('Current');
    expect(t('cancel')).toBe('Cancel');
    expect(t('confirm')).toBe('Confirm');
  });

  test('returns key for unknown translation', () => {
    expect(t('unknownKey')).toBe('unknownKey');
  });

  test('handles parameter interpolation', () => {
    setLocale('zh-CN');
    expect(t('continueFromVersionN', { version: 3 })).toBe('基于 v3 继续优化');

    setLocale('en-US');
    expect(t('continueFromVersionN', { version: 3 })).toBe('Continue from v3');
  });

  test('handles multiple parameters', () => {
    // Add a test key if needed, or test with existing keys
    setLocale('zh-CN');
    const result = t('continueFromVersionN', { version: 5 });
    expect(result).toContain('5');
  });
});

// ========== initI18n 测试 ==========

describe('initI18n', () => {
  test('sets locale to zh-CN for Chinese prompt', () => {
    initI18n('请帮我优化这个提示词');
    expect(getLocale()).toBe('zh-CN');
  });

  test('sets locale to en-US for English prompt', () => {
    initI18n('Please optimize this prompt');
    expect(getLocale()).toBe('en-US');
  });

  test('sets locale to en-US for empty prompt', () => {
    initI18n('');
    expect(getLocale()).toBe('en-US');
  });
});

// ========== 评分等级翻译测试 ==========

describe('Score grade translations', () => {
  test('Chinese score grades', () => {
    setLocale('zh-CN');
    expect(t('scorePoor')).toBe('较差');
    expect(t('scoreFair')).toBe('一般');
    expect(t('scoreGood')).toBe('良好');
    expect(t('scoreExcellent')).toBe('优秀');
  });

  test('English score grades', () => {
    setLocale('en-US');
    expect(t('scorePoor')).toBe('Poor');
    expect(t('scoreFair')).toBe('Fair');
    expect(t('scoreGood')).toBe('Good');
    expect(t('scoreExcellent')).toBe('Excellent');
  });
});
