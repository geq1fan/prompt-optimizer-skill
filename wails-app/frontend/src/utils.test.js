import { describe, test, expect } from 'vitest';
import {
  getScoreGradeClass,
  getScoreGradeText,
  escapeHtml,
  parseMarkdown,
  findHistoryItem,
  getVersionData,
  formatCountdown,
  getCountdownStatusClass,
} from './utils.js';

// ========== getScoreGradeClass 测试 ==========

describe('getScoreGradeClass', () => {
  test.each([
    [0, 'poor'],
    [20, 'poor'],
    [40, 'poor'],
    [41, 'fair'],
    [50, 'fair'],
    [60, 'fair'],
    [61, 'good'],
    [70, 'good'],
    [80, 'good'],
    [81, 'excellent'],
    [90, 'excellent'],
    [100, 'excellent'],
  ])('score %i returns "%s"', (score, expected) => {
    expect(getScoreGradeClass(score)).toBe(expected);
  });

  test('handles negative scores as poor', () => {
    expect(getScoreGradeClass(-10)).toBe('poor');
  });

  test('handles scores above 100 as excellent', () => {
    expect(getScoreGradeClass(150)).toBe('excellent');
  });
});

// ========== getScoreGradeText 测试 ==========

describe('getScoreGradeText', () => {
  test.each([
    [0, 'Poor'],
    [40, 'Poor'],
    [41, 'Fair'],
    [60, 'Fair'],
    [61, 'Good'],
    [80, 'Good'],
    [81, 'Excellent'],
    [100, 'Excellent'],
  ])('score %i returns "%s"', (score, expected) => {
    expect(getScoreGradeText(score)).toBe(expected);
  });
});

// ========== escapeHtml 测试 ==========

describe('escapeHtml', () => {
  test('escapes ampersand', () => {
    expect(escapeHtml('a & b')).toBe('a &amp; b');
  });

  test('escapes less than', () => {
    expect(escapeHtml('a < b')).toBe('a &lt; b');
  });

  test('escapes greater than', () => {
    expect(escapeHtml('a > b')).toBe('a &gt; b');
  });

  test('escapes double quotes', () => {
    expect(escapeHtml('a "b" c')).toBe('a &quot;b&quot; c');
  });

  test('escapes single quotes', () => {
    expect(escapeHtml("a 'b' c")).toBe('a &#39;b&#39; c');
  });

  test('escapes multiple special characters', () => {
    expect(escapeHtml('<script>alert("xss")</script>')).toBe(
      '&lt;script&gt;alert(&quot;xss&quot;)&lt;/script&gt;'
    );
  });

  test('returns empty string for null', () => {
    expect(escapeHtml(null)).toBe('');
  });

  test('returns empty string for undefined', () => {
    expect(escapeHtml(undefined)).toBe('');
  });

  test('returns empty string for empty input', () => {
    expect(escapeHtml('')).toBe('');
  });

  test('converts numbers to string', () => {
    expect(escapeHtml(123)).toBe('123');
  });

  test('handles Chinese characters', () => {
    expect(escapeHtml('中文 & 日本語')).toBe('中文 &amp; 日本語');
  });
});

// ========== parseMarkdown 测试 ==========

describe('parseMarkdown', () => {
  test('returns empty string for null', () => {
    expect(parseMarkdown(null)).toBe('');
  });

  test('returns empty string for undefined', () => {
    expect(parseMarkdown(undefined)).toBe('');
  });

  test('returns empty string for empty input', () => {
    expect(parseMarkdown('')).toBe('');
  });

  describe('headers', () => {
    test('converts h1', () => {
      expect(parseMarkdown('# Title')).toContain('<h1>Title</h1>');
    });

    test('converts h2', () => {
      expect(parseMarkdown('## Subtitle')).toContain('<h2>Subtitle</h2>');
    });

    test('converts h3', () => {
      expect(parseMarkdown('### Section')).toContain('<h3>Section</h3>');
    });
  });

  describe('text formatting', () => {
    test('converts bold', () => {
      expect(parseMarkdown('**bold text**')).toContain('<strong>bold text</strong>');
    });

    test('converts italic', () => {
      expect(parseMarkdown('*italic text*')).toContain('<em>italic text</em>');
    });

    test('converts inline code', () => {
      expect(parseMarkdown('use `code` here')).toContain('<code>code</code>');
    });
  });

  describe('lists', () => {
    test('converts unordered list', () => {
      const input = '- item 1\n- item 2';
      const result = parseMarkdown(input);
      expect(result).toContain('<li>item 1</li>');
      expect(result).toContain('<li>item 2</li>');
      expect(result).toContain('<ul>');
    });
  });

  describe('XSS prevention', () => {
    test('escapes HTML in content', () => {
      const result = parseMarkdown('<script>alert("xss")</script>');
      expect(result).not.toContain('<script>');
      expect(result).toContain('&lt;script&gt;');
    });
  });
});

// ========== findHistoryItem 测试 ==========

describe('findHistoryItem', () => {
  const mockHistory = [
    { iterationId: 'iter-001', score: 65 },
    { iterationId: 'iter-002', score: 75 },
    { iterationId: 'iter-003', score: 85 },
  ];

  test('finds existing item', () => {
    const result = findHistoryItem(mockHistory, 'iter-002');
    expect(result).toBeDefined();
    expect(result.score).toBe(75);
  });

  test('returns undefined for non-existent id', () => {
    const result = findHistoryItem(mockHistory, 'iter-999');
    expect(result).toBeUndefined();
  });

  test('returns undefined for null history', () => {
    const result = findHistoryItem(null, 'iter-001');
    expect(result).toBeUndefined();
  });

  test('returns undefined for undefined history', () => {
    const result = findHistoryItem(undefined, 'iter-001');
    expect(result).toBeUndefined();
  });

  test('returns undefined for empty history', () => {
    const result = findHistoryItem([], 'iter-001');
    expect(result).toBeUndefined();
  });
});

// ========== getVersionData 测试 ==========

describe('getVersionData', () => {
  const mockInputData = {
    version: 3,
    current: {
      score: 85,
      optimizedPrompt: 'Current prompt',
    },
    history: [
      { iterationId: 'iter-001', score: 65, optimizedPrompt: 'First' },
      { iterationId: 'iter-002', score: 75, optimizedPrompt: 'Second' },
    ],
  };

  test('returns current version data', () => {
    const result = getVersionData(mockInputData, 'current');
    expect(result.version).toBe(3);
    expect(result.score).toBe(85);
    expect(result.optimizedPrompt).toBe('Current prompt');
  });

  test('returns history version data', () => {
    const result = getVersionData(mockInputData, 'iter-001');
    expect(result.version).toBe(1);
    expect(result.score).toBe(65);
    expect(result.optimizedPrompt).toBe('First');
  });

  test('returns second history version', () => {
    const result = getVersionData(mockInputData, 'iter-002');
    expect(result.version).toBe(2);
    expect(result.score).toBe(75);
  });

  test('returns default for non-existent id', () => {
    const result = getVersionData(mockInputData, 'iter-999');
    expect(result.version).toBe(0);
    expect(result.score).toBe(0);
    expect(result.optimizedPrompt).toBe('');
  });

  test('returns default for null inputData', () => {
    const result = getVersionData(null, 'current');
    expect(result.version).toBe(0);
    expect(result.score).toBe(0);
  });

  test('handles missing current data', () => {
    const result = getVersionData({ version: 1 }, 'current');
    expect(result.version).toBe(1);
    expect(result.score).toBe(0);
    expect(result.optimizedPrompt).toBe('');
  });

  test('handles empty history', () => {
    const data = { version: 1, current: { score: 50 }, history: [] };
    const result = getVersionData(data, 'iter-001');
    expect(result.version).toBe(0);
  });
});

// ========== formatCountdown 测试 ==========

describe('formatCountdown', () => {
  test.each([
    [0, '0:00'],
    [30, '0:30'],
    [59, '0:59'],
    [60, '1:00'],
    [90, '1:30'],
    [600, '10:00'],
    [3599, '59:59'],
    [3600, '60:00'],
  ])('%i seconds returns "%s"', (seconds, expected) => {
    expect(formatCountdown(seconds)).toBe(expected);
  });

  test('handles negative seconds as 0:00', () => {
    expect(formatCountdown(-10)).toBe('0:00');
  });
});

// ========== getCountdownStatusClass 测试 ==========

describe('getCountdownStatusClass', () => {
  test('returns empty for normal time', () => {
    expect(getCountdownStatusClass(120)).toBe('');
    expect(getCountdownStatusClass(60)).toBe('');
  });

  test('returns warning for < 60 seconds', () => {
    expect(getCountdownStatusClass(59)).toBe('warning');
    expect(getCountdownStatusClass(30)).toBe('warning');
  });

  test('returns critical for < 30 seconds', () => {
    expect(getCountdownStatusClass(29)).toBe('critical');
    expect(getCountdownStatusClass(0)).toBe('critical');
  });

  test('returns critical for negative seconds', () => {
    expect(getCountdownStatusClass(-5)).toBe('critical');
  });
});
