import { describe, test, expect, beforeAll } from 'vitest';
import { readFileSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { setLocale } from './i18n/index.js';
import {
  getScoreGradeClass,
  getScoreGradeText,
  escapeHtml,
  parseMarkdown,
  findHistoryItem,
  getVersionData,
} from './utils.js';

// 获取 testdata 目录路径
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const testdataDir = join(__dirname, '..', '..', 'testdata');

// 加载测试数据文件
function loadTestData(filename) {
  const filepath = join(testdataDir, filename);
  const content = readFileSync(filepath, 'utf-8');
  return JSON.parse(content);
}

// ========== 集成测试 - 使用预制 testdata ==========

describe('Integration: Basic Input (v1)', () => {
  let inputData;

  beforeAll(() => {
    inputData = loadTestData('input_v1_basic.json');
  });

  test('loads version correctly', () => {
    expect(inputData.version).toBe(1);
  });

  test('loads original prompt', () => {
    expect(inputData.originalPrompt).toBe('帮我写一个登录页面');
  });

  test('loads current iteration', () => {
    expect(inputData.current.iterationId).toBe('iter-001');
    expect(inputData.current.score).toBe(75);
  });

  test('has suggested directions', () => {
    expect(inputData.current.suggestedDirections).toHaveLength(3);
    expect(inputData.current.suggestedDirections[0].id).toBe('error-handling');
    expect(inputData.current.suggestedDirections[0].label).toBe('错误处理');
  });

  test('has no history', () => {
    expect(inputData.history).toHaveLength(0);
  });

  test('getScoreGradeClass returns correct class', () => {
    expect(getScoreGradeClass(inputData.current.score)).toBe('good');
  });

  test('getScoreGradeText returns correct text', () => {
    setLocale('en-US');
    expect(getScoreGradeText(inputData.current.score)).toBe('Good');
  });
});

describe('Integration: Input with History (v3)', () => {
  let inputData;

  beforeAll(() => {
    inputData = loadTestData('input_v3_with_history.json');
  });

  test('loads version correctly', () => {
    expect(inputData.version).toBe(3);
  });

  test('loads history with correct count', () => {
    expect(inputData.history).toHaveLength(2);
  });

  test('history items have correct structure', () => {
    const history1 = inputData.history[0];
    expect(history1.iterationId).toBe('iter-001');
    expect(history1.score).toBe(55);
    expect(history1.userFeedback.selectedDirections).toContain('structure');
  });

  test('findHistoryItem finds correct item', () => {
    const item = findHistoryItem(inputData.history, 'iter-002');
    expect(item).toBeDefined();
    expect(item.score).toBe(72);
  });

  test('findHistoryItem returns undefined for non-existent', () => {
    const item = findHistoryItem(inputData.history, 'iter-999');
    expect(item).toBeUndefined();
  });

  test('getVersionData returns current version', () => {
    const data = getVersionData(inputData, 'current');
    expect(data.version).toBe(3);
    expect(data.score).toBe(88);
  });

  test('getVersionData returns history version', () => {
    const data = getVersionData(inputData, 'iter-001');
    expect(data.version).toBe(1);
    expect(data.score).toBe(55);
  });
});

describe('Integration: Session Data (v4)', () => {
  let sessionData;

  beforeAll(() => {
    sessionData = loadTestData('session_v4.json');
  });

  test('loads version correctly', () => {
    expect(sessionData.version).toBe(4);
  });

  test('has sessionId', () => {
    expect(sessionData.sessionId).toBe('session_1705632000000');
  });

  test('has timestamps', () => {
    expect(sessionData.createdAt).toBe('2026-01-22T10:00:00Z');
    expect(sessionData.updatedAt).toBe('2026-01-22T10:05:00Z');
  });

  test('has project metadata', () => {
    expect(sessionData.projectPath).toBe('/path/to/project');
    expect(sessionData.lang).toBe('cn');
    expect(sessionData.mode).toBe('professional');
  });

  test('has status and lastAction', () => {
    expect(sessionData.status).toBe('active');
    expect(sessionData.lastAction).toBe('submit');
  });

  test('loads current iteration', () => {
    expect(sessionData.current.iterationId).toBe('iter-002');
    expect(sessionData.current.score).toBe(85);
  });

  test('has suggested directions', () => {
    expect(sessionData.current.suggestedDirections).toHaveLength(3);
    expect(sessionData.current.suggestedDirections[0].id).toBe('refresh-token');
  });

  test('loads history with correct count', () => {
    expect(sessionData.history).toHaveLength(1);
    expect(sessionData.history[0].iterationId).toBe('iter-001');
    expect(sessionData.history[0].score).toBe(65);
  });

  test('history item has user feedback', () => {
    const feedback = sessionData.history[0].userFeedback;
    expect(feedback.selectedDirections).toContain('security');
    expect(feedback.userInput).toBeTruthy();
  });

  test('getScoreGradeClass returns correct class for current score', () => {
    // Score 85 is > 80, so it should be 'excellent'
    expect(getScoreGradeClass(sessionData.current.score)).toBe('excellent');
  });

  test('getVersionData works with session data', () => {
    const data = getVersionData(sessionData, 'current');
    expect(data.score).toBe(85);
  });
});

describe('Integration: Long History (v5)', () => {
  let inputData;

  beforeAll(() => {
    inputData = loadTestData('input_v5_long_history.json');
  });

  test('loads version correctly', () => {
    expect(inputData.version).toBe(5);
  });

  test('loads all history items', () => {
    expect(inputData.history).toHaveLength(4);
  });

  test('history shows score progression', () => {
    const scores = inputData.history.map(h => h.score);
    expect(scores).toEqual([40, 55, 70, 82]);
  });

  test('current score is highest', () => {
    expect(inputData.current.score).toBe(92);
    expect(getScoreGradeClass(inputData.current.score)).toBe('excellent');
  });

  test('all history items have user feedback', () => {
    inputData.history.forEach((item, index) => {
      expect(item.userFeedback).toBeDefined();
      expect(item.userFeedback.selectedDirections.length).toBeGreaterThan(0);
      expect(item.userFeedback.userInput).toBeTruthy();
    });
  });

  test('getVersionData works for all history versions', () => {
    for (let i = 0; i < inputData.history.length; i++) {
      const id = inputData.history[i].iterationId;
      const data = getVersionData(inputData, id);
      expect(data.version).toBe(i + 1);
      expect(data.score).toBe(inputData.history[i].score);
    }
  });
});

describe('Integration: Empty Input', () => {
  let inputData;

  beforeAll(() => {
    inputData = loadTestData('input_empty.json');
  });

  test('handles empty original prompt', () => {
    expect(inputData.originalPrompt).toBe('');
  });

  test('handles empty optimized prompt', () => {
    expect(inputData.current.optimizedPrompt).toBe('');
  });

  test('handles zero score', () => {
    expect(inputData.current.score).toBe(0);
    expect(getScoreGradeClass(0)).toBe('poor');
  });

  test('handles empty directions', () => {
    expect(inputData.current.suggestedDirections).toHaveLength(0);
  });

  test('handles empty history', () => {
    expect(inputData.history).toHaveLength(0);
  });

  test('getVersionData handles empty data gracefully', () => {
    const data = getVersionData(inputData, 'current');
    expect(data.score).toBe(0);
    expect(data.optimizedPrompt).toBe('');
  });
});

describe('Integration: Unicode Input', () => {
  let inputData;

  beforeAll(() => {
    inputData = loadTestData('input_unicode.json');
  });

  test('preserves Chinese characters', () => {
    expect(inputData.originalPrompt).toContain('中文测试');
  });

  test('preserves Japanese characters', () => {
    expect(inputData.originalPrompt).toContain('日本語');
  });

  test('preserves Korean characters', () => {
    expect(inputData.originalPrompt).toContain('한국어');
  });

  test('preserves Arabic characters', () => {
    expect(inputData.originalPrompt).toContain('العربية');
  });

  test('preserves Emoji', () => {
    expect(inputData.originalPrompt).toContain('🎉');
    expect(inputData.current.optimizedPrompt).toContain('🌍');
  });

  test('escapeHtml preserves unicode', () => {
    const escaped = escapeHtml('中文 & English');
    expect(escaped).toBe('中文 &amp; English');
  });

  test('parseMarkdown handles unicode headers', () => {
    const result = parseMarkdown('# 中文标题');
    expect(result).toContain('<h1>中文标题</h1>');
  });
});

describe('Integration: Markdown Content Parsing', () => {
  let inputData;

  beforeAll(() => {
    inputData = loadTestData('input_v3_with_history.json');
  });

  test('optimized prompt contains Markdown headers', () => {
    const prompt = inputData.current.optimizedPrompt;
    expect(prompt).toContain('# Role:');
    expect(prompt).toContain('## Goals');
  });

  test('parseMarkdown converts headers correctly', () => {
    const result = parseMarkdown(inputData.current.optimizedPrompt);
    expect(result).toContain('<h1>');
    expect(result).toContain('<h2>');
  });

  test('parseMarkdown converts list items', () => {
    const markdown = '- Item 1\n- Item 2';
    const result = parseMarkdown(markdown);
    expect(result).toContain('<li>Item 1</li>');
    expect(result).toContain('<ul>');
  });

  test('parseMarkdown converts code blocks', () => {
    const prompt = inputData.current.optimizedPrompt;
    if (prompt.includes('```')) {
      const result = parseMarkdown(prompt);
      expect(result).toContain('<pre>');
      expect(result).toContain('<code>');
    }
  });
});

describe('Integration: Direction Selection Simulation', () => {
  let inputData;

  beforeAll(() => {
    inputData = loadTestData('input_v1_basic.json');
  });

  test('can extract direction IDs for selection', () => {
    const directions = inputData.current.suggestedDirections;
    const ids = directions.map(d => d.id);
    expect(ids).toContain('error-handling');
    expect(ids).toContain('ui-style');
    expect(ids).toContain('security');
  });

  test('direction labels are displayable', () => {
    const directions = inputData.current.suggestedDirections;
    directions.forEach(d => {
      expect(d.label).toBeTruthy();
      expect(typeof d.label).toBe('string');
      // Labels should be escaped for display
      const escaped = escapeHtml(d.label);
      expect(escaped).toBeTruthy();
    });
  });

  test('direction descriptions are displayable', () => {
    const directions = inputData.current.suggestedDirections;
    directions.forEach(d => {
      expect(d.description).toBeTruthy();
      const escaped = escapeHtml(d.description);
      expect(escaped).toBeTruthy();
    });
  });
});

describe('Integration: Result Structure Simulation', () => {
  test('submit result structure', () => {
    const result = {
      action: 'submit',
      selectedDirections: ['error-handling', 'ui-style'],
      userInput: '请添加错误处理',
    };

    expect(result.action).toBe('submit');
    expect(result.selectedDirections).toHaveLength(2);
    expect(result.userInput).toBeTruthy();
  });

  test('cancel result structure', () => {
    const result = {
      action: 'cancel',
      selectedDirections: [],
      userInput: '',
    };

    expect(result.action).toBe('cancel');
    expect(result.selectedDirections).toHaveLength(0);
    expect(result.userInput).toBe('');
  });

  test('rollback result structure', () => {
    const result = {
      action: 'rollback',
      selectedDirections: ['examples'],
      userInput: '基于第一版优化',
      rollbackToIteration: 'iter-001',
    };

    expect(result.action).toBe('rollback');
    expect(result.rollbackToIteration).toBe('iter-001');
  });

  test('timeout result structure', () => {
    const result = {
      action: 'timeout',
      selectedDirections: [],
      userInput: '',
    };

    expect(result.action).toBe('timeout');
  });
});
