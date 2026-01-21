// ========== i18n 国际化核心模块 ==========
import zhCN from './zh-CN.js';
import enUS from './en-US.js';

const locales = {
  'zh-CN': zhCN,
  'en-US': enUS
};

let currentLocale = 'zh-CN';

/**
 * 检测文本是否包含中文字符
 * @param {string} text - 待检测文本
 * @returns {boolean} 是否包含中文
 */
export function containsChinese(text) {
  if (!text) return false;
  return /[\u4e00-\u9fff]/.test(text);
}

/**
 * 根据文本内容检测语言
 * @param {string} text - 待检测文本（通常是 originalPrompt）
 * @returns {string} 语言代码 ('zh-CN' 或 'en-US')
 */
export function detectLanguage(text) {
  return containsChinese(text) ? 'zh-CN' : 'en-US';
}

/**
 * 设置当前语言
 * @param {string} locale - 语言代码
 */
export function setLocale(locale) {
  if (locales[locale]) {
    currentLocale = locale;
    // 仅在浏览器环境中更新 DOM
    if (typeof document !== 'undefined') {
      document.documentElement.lang = locale;
      updateDOMTranslations();
    }
  }
}

/**
 * 获取当前语言
 * @returns {string} 当前语言代码
 */
export function getLocale() {
  return currentLocale;
}

/**
 * 翻译函数
 * @param {string} key - 翻译键
 * @param {Object} params - 插值参数 (可选)
 * @returns {string} 翻译后的文本
 */
export function t(key, params = {}) {
  let message = locales[currentLocale]?.[key] || locales['en-US']?.[key] || key;

  // 处理 {variable} 格式的插值
  Object.keys(params).forEach(k => {
    message = message.replace(new RegExp(`\\{${k}\\}`, 'g'), params[k]);
  });

  return message;
}

/**
 * 更新 DOM 中所有带 data-i18n 属性的元素
 */
export function updateDOMTranslations() {
  // 仅在浏览器环境中执行
  if (typeof document === 'undefined') {
    return;
  }

  // 更新文本内容
  document.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.getAttribute('data-i18n');
    el.textContent = t(key);
  });

  // 更新 placeholder
  document.querySelectorAll('[data-i18n-placeholder]').forEach(el => {
    const key = el.getAttribute('data-i18n-placeholder');
    el.placeholder = t(key);
  });

  // 更新 title
  document.querySelectorAll('[data-i18n-title]').forEach(el => {
    const key = el.getAttribute('data-i18n-title');
    el.title = t(key);
  });
}

/**
 * 初始化 i18n
 * @param {string} originalPrompt - 原始提示词，用于检测语言
 */
export function initI18n(originalPrompt) {
  const locale = detectLanguage(originalPrompt);
  setLocale(locale);
}

// 默认导出
export default {
  t,
  setLocale,
  getLocale,
  detectLanguage,
  containsChinese,
  initI18n,
  updateDOMTranslations
};
