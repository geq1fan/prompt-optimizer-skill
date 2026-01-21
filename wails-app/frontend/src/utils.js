// ========== 纯函数 - 可独立测试 ==========

import { t } from './i18n/index.js';

/**
 * 根据分数返回等级 CSS 类名
 * @param {number} score - 分数 (0-100)
 * @returns {string} CSS 类名
 */
export function getScoreGradeClass(score) {
  if (score <= 40) return 'poor';
  if (score <= 60) return 'fair';
  if (score <= 80) return 'good';
  return 'excellent';
}

/**
 * 根据分数返回等级文字
 * @param {number} score - 分数 (0-100)
 * @returns {string} 等级文字
 */
export function getScoreGradeText(score) {
  if (score <= 40) return t('scorePoor');
  if (score <= 60) return t('scoreFair');
  if (score <= 80) return t('scoreGood');
  return t('scoreExcellent');
}

/**
 * HTML 转义 - XSS 防护
 * @param {string} text - 原始文本
 * @returns {string} 转义后的文本
 */
export function escapeHtml(text) {
  if (!text) return '';
  const escapeMap = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#39;',
  };
  return String(text).replace(/[&<>"']/g, char => escapeMap[char]);
}

/**
 * 简单 Markdown 解析
 * @param {string} text - Markdown 文本
 * @returns {string} HTML 字符串
 */
export function parseMarkdown(text) {
  if (!text) return '';

  let html = escapeHtml(text);

  // 标题
  html = html.replace(/^### (.+)$/gm, '<h3>$1</h3>');
  html = html.replace(/^## (.+)$/gm, '<h2>$1</h2>');
  html = html.replace(/^# (.+)$/gm, '<h1>$1</h1>');

  // 粗体和斜体
  html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
  html = html.replace(/\*(.+?)\*/g, '<em>$1</em>');

  // 代码块
  html = html.replace(/```[\s\S]*?```/g, match => {
    const code = match.slice(3, -3).trim();
    return `<pre><code>${code}</code></pre>`;
  });

  // 行内代码
  html = html.replace(/`([^`]+)`/g, '<code>$1</code>');

  // 列表
  html = html.replace(/^- (.+)$/gm, '<li>$1</li>');
  html = html.replace(/(<li>.*<\/li>\n?)+/g, '<ul>$&</ul>');

  // 表格 (简单处理)
  const tableRegex = /\|(.+)\|\n\|[-|]+\|\n((?:\|.+\|\n?)+)/g;
  html = html.replace(tableRegex, (match, header, body) => {
    const headers = header.split('|').filter(s => s.trim());
    const rows = body.trim().split('\n').map(row =>
      row.split('|').filter(s => s.trim())
    );

    let table = '<table><thead><tr>';
    headers.forEach(h => table += `<th>${h.trim()}</th>`);
    table += '</tr></thead><tbody>';
    rows.forEach(row => {
      table += '<tr>';
      row.forEach(cell => table += `<td>${cell.trim()}</td>`);
      table += '</tr>';
    });
    table += '</tbody></table>';
    return table;
  });

  // 段落
  html = html.replace(/\n\n/g, '</p><p>');
  html = '<p>' + html + '</p>';
  html = html.replace(/<p><\/p>/g, '');
  html = html.replace(/<p>(<h[123]>)/g, '$1');
  html = html.replace(/(<\/h[123]>)<\/p>/g, '$1');
  html = html.replace(/<p>(<ul>)/g, '$1');
  html = html.replace(/(<\/ul>)<\/p>/g, '$1');
  html = html.replace(/<p>(<pre>)/g, '$1');
  html = html.replace(/(<\/pre>)<\/p>/g, '$1');
  html = html.replace(/<p>(<table>)/g, '$1');
  html = html.replace(/(<\/table>)<\/p>/g, '$1');

  return html;
}

/**
 * 在历史记录中查找指定迭代
 * @param {Array} history - 历史记录数组
 * @param {string} iterationId - 迭代 ID
 * @returns {Object|undefined} 找到的历史项
 */
export function findHistoryItem(history, iterationId) {
  if (!Array.isArray(history)) return undefined;
  return history.find(item => item.iterationId === iterationId);
}

/**
 * 获取版本数据
 * @param {Object} inputData - 完整输入数据
 * @param {string} id - 版本 ID ('current' 或 iterationId)
 * @returns {Object} 版本数据
 */
export function getVersionData(inputData, id) {
  if (!inputData) {
    return { version: 0, score: 0, optimizedPrompt: '' };
  }

  if (id === 'current') {
    return {
      version: inputData.version,
      score: inputData.current?.score || 0,
      optimizedPrompt: inputData.current?.optimizedPrompt || ''
    };
  }

  const history = inputData.history || [];
  for (let i = 0; i < history.length; i++) {
    if (history[i].iterationId === id) {
      return {
        version: i + 1,
        score: history[i].score,
        optimizedPrompt: history[i].optimizedPrompt
      };
    }
  }

  return { version: 0, score: 0, optimizedPrompt: '' };
}
