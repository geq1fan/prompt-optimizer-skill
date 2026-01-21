// ========== i18n 导入 ==========
import { t, initI18n } from './i18n/index.js';
import {
  getScoreGradeClass,
  getScoreGradeText,
  escapeHtml,
  findHistoryItem as findHistoryItemInArray,
  getVersionData as getVersionDataFromInput
} from './utils.js';

// ========== 全局状态 ==========
let inputData = null;
let selectedDirections = new Set();
let countdownInterval = null;
let currentPromptMode = 'render';  // 'render' | 'source'
let currentDiffMode = 'split';     // 'unified' | 'split'
let viewingHistoryVersion = null;  // 当前正在查看的历史版本 iterationId，null 表示查看当前版本

// ========== 初始化 ==========
document.addEventListener('DOMContentLoaded', async () => {
  try {
    await loadData();
  } catch (error) {
    console.error('Failed to load data:', error);
    showError(t('loadError') + ': ' + error.message);
  }
});

// ========== 数据加载 ==========
async function loadData() {
  try {
    inputData = await window.go.main.App.GetInputData();
    // 根据原始 Prompt 初始化 i18n
    initI18n(inputData.originalPrompt);
    hideLoading();
    renderUI();
    startCountdown();
  } catch (error) {
    throw error;
  }
}

function hideLoading() {
  document.getElementById('loading-screen').classList.add('hidden');
  document.getElementById('main-content').classList.remove('hidden');
}

function showError(message) {
  alert(message);
}

// ========== UI 渲染 ==========
function renderUI() {
  renderVersion();
  renderScore(inputData.current.score);
  renderReports();
  renderOptimizedPrompt();
  renderDirections();
  renderHistory();
  initCompareSelectors();
  initUserInput();
  initScorePopover();
}

// 获取当前迭代版本（基于历史记录数量）
function getCurrentIterationVersion() {
  return (inputData.history || []).length + 1;
}

function renderVersion() {
  const versionBadge = document.getElementById('version-badge');
  versionBadge.textContent = 'v' + getCurrentIterationVersion();
}

function renderScore(score) {
  const scoreValue = document.getElementById('score-value');
  scoreValue.textContent = score;

  const scoreRingProgress = document.getElementById('score-ring-progress');
  const circumference = 2 * Math.PI * 85; // r = 85 (Bento hero ring)
  const offset = circumference - (score / 100) * circumference;
  scoreRingProgress.style.strokeDasharray = circumference;
  scoreRingProgress.style.strokeDashoffset = offset;

  const gradeClass = getScoreGradeClass(score);
  scoreRingProgress.setAttribute('class', 'score-ring-progress ' + gradeClass);

  const scoreGrade = document.getElementById('score-grade');
  scoreGrade.textContent = getScoreGradeText(score);
  scoreGrade.className = 'score-grade ' + gradeClass;

  const scoreRingWrapper = document.querySelector('.score-ring-wrapper');
  if (scoreRingWrapper) {
    scoreRingWrapper.className = 'score-ring-wrapper glow ' + gradeClass;
  }
}

function renderReports() {
  // 评审报告 → 右侧卡片
  const reviewReport = document.getElementById('review-report');
  reviewReport.innerHTML = parseMarkdown(inputData.current.reviewReport || t('noReviewReport'));

  // 评估报告 → 弹出层
  renderEvaluationPopover();
}

// ========== 评估报告弹出层解析与渲染 ==========

/**
 * 解析维度评分
 * 输入格式: - **任务表达**: 95 - 清晰表达了优化目标
 * 返回: [{ name: '任务表达', score: 95, comment: '清晰表达了优化目标' }, ...]
 */
function parseEvalDimensions(markdown) {
  const dimensions = [];
  const dimensionRegex = /[-*•]\s*\*{0,2}(.+?)\*{0,2}[:：]\s*(\d+)\s*[-–]\s*(.+)/g;
  let match;

  // 只解析 "### 维度评分" 区块
  const sectionMatch = markdown.match(/###\s*维度评分[\s\S]*?(?=###|$)/i);
  if (!sectionMatch) return dimensions;

  const section = sectionMatch[0];
  while ((match = dimensionRegex.exec(section)) !== null) {
    dimensions.push({
      name: match[1].trim(),
      score: parseInt(match[2], 10),
      comment: match[3].trim()
    });
  }
  return dimensions;
}

/**
 * 解析核心优势
 * 返回: ['结构清晰', '上下文丰富', ...]
 */
function parseEvalAdvantages(markdown) {
  const sectionMatch = markdown.match(/###\s*核心优势[\s\S]*?(?=###|$)/i);
  if (!sectionMatch) return [];

  const items = [];
  const itemRegex = /[-*•]\s*(.+)/g;
  let match;
  while ((match = itemRegex.exec(sectionMatch[0])) !== null) {
    const text = match[1].trim();
    if (text && !text.startsWith('#')) {
      items.push(text);
    }
  }
  return items;
}

/**
 * 解析改进建议
 * 返回: ['添加边界条件', '补充错误处理', ...]
 */
function parseEvalSuggestions(markdown) {
  const sectionMatch = markdown.match(/###\s*改进建议[\s\S]*?(?=###|$)/i);
  if (!sectionMatch) return [];

  const items = [];
  const itemRegex = /[-*•]\s*(.+)/g;
  let match;
  while ((match = itemRegex.exec(sectionMatch[0])) !== null) {
    const text = match[1].trim();
    if (text && !text.startsWith('#')) {
      items.push(text);
    }
  }
  return items;
}

/**
 * 解析迭代方向
 * 返回: { conclusion: '符合预期', nextStep: '可以尝试添加示例' }
 */
function parseEvalDirection(markdown) {
  const sectionMatch = markdown.match(/###\s*迭代方向评估[\s\S]*?(?=###|$)/i);
  if (!sectionMatch) return null;

  const section = sectionMatch[0];
  const conclusionMatch = section.match(/\*\*结论\*\*[:：]\s*(.+)/);
  const nextStepMatch = section.match(/\*\*下一步建议\*\*[:：]\s*(.+)/);

  return {
    conclusion: conclusionMatch ? conclusionMatch[1].trim() : '',
    nextStep: nextStepMatch ? nextStepMatch[1].trim() : ''
  };
}

/**
 * 渲染维度评分（进度条 + 分数 + 简评）
 */
function renderPopoverDimensions(dimensions) {
  if (!dimensions.length) return '';

  const getBarColor = (score) => {
    if (score >= 90) return 'var(--color-green, #22c55e)';
    if (score >= 80) return 'var(--color-accent, #3b82f6)';
    if (score >= 70) return 'var(--color-orange, #eab308)';
    return 'var(--color-red, #ef4444)';
  };

  const items = dimensions.map(d => `
    <div class="eval-dimension-item">
      <span class="eval-dimension-name">${escapeHtml(d.name)}</span>
      <div class="eval-dimension-bar-wrapper">
        <div class="eval-dimension-bar-fill" style="width: ${d.score}%; background: ${getBarColor(d.score)}"></div>
      </div>
      <span class="eval-dimension-value">${d.score}</span>
    </div>
    <div class="eval-dimension-comment">${escapeHtml(d.comment)}</div>
  `).join('');

  return `<div class="eval-dimensions">${items}</div>`;
}

/**
 * 渲染核心优势（标签样式）
 */
function renderPopoverAdvantages(advantages) {
  if (!advantages.length) return '';

  const tags = advantages.map(a => `<span class="eval-tag">${escapeHtml(a)}</span>`).join('');
  return `
    <div class="eval-section-title">${t('coreAdvantages') || '核心优势'}</div>
    <div class="eval-tags">${tags}</div>
  `;
}

/**
 * 渲染改进建议（紧凑列表）
 */
function renderPopoverSuggestions(suggestions) {
  if (!suggestions.length) return '';

  const items = suggestions.map(s => `<li>${escapeHtml(s)}</li>`).join('');
  return `
    <div class="eval-section-title">${t('improvementSuggestions') || '改进建议'}</div>
    <ul class="eval-suggestions">${items}</ul>
  `;
}

/**
 * 渲染迭代方向（状态标签 + 提示）
 */
function renderPopoverDirection(direction) {
  if (!direction || !direction.conclusion) return '';

  const getLabelClass = (conclusion) => {
    if (conclusion.includes('符合预期')) return 'positive';
    if (conclusion.includes('偏离目标')) return 'negative';
    return 'neutral';
  };

  return `
    <div class="eval-direction">
      <span class="eval-direction-label ${getLabelClass(direction.conclusion)}">${escapeHtml(direction.conclusion)}</span>
      <span class="eval-direction-hint">${escapeHtml(direction.nextStep)}</span>
    </div>
  `;
}

/**
 * 渲染评估报告到弹出层
 */
function renderEvaluationPopover() {
  const evaluationEl = document.getElementById('evaluation-report');
  if (!evaluationEl) return;

  const evalMarkdown = inputData.current.evaluationReport || '';
  if (!evalMarkdown) {
    evaluationEl.innerHTML = '<p style="color: var(--text-tertiary); font-size: 12px;">' + (t('noEvaluationReport') || '暂无评估数据') + '</p>';
    return;
  }

  const dimensions = parseEvalDimensions(evalMarkdown);
  const advantages = parseEvalAdvantages(evalMarkdown);
  const suggestions = parseEvalSuggestions(evalMarkdown);
  const direction = parseEvalDirection(evalMarkdown);

  evaluationEl.innerHTML =
    renderPopoverDimensions(dimensions) +
    renderPopoverAdvantages(advantages) +
    renderPopoverSuggestions(suggestions) +
    renderPopoverDirection(direction);
}

/**
 * 调整弹出层位置，防止溢出视口
 */
function adjustPopoverPosition() {
  const popover = document.getElementById('score-popover');
  const scoreCard = document.getElementById('score-card');
  if (!popover || !scoreCard) return;

  const cardRect = scoreCard.getBoundingClientRect();
  const popoverHeight = popover.offsetHeight;
  const cardCenterY = cardRect.top + cardRect.height / 2;

  let offsetY = 0;
  // 检查顶部溢出
  if (cardCenterY - popoverHeight / 2 < 10) {
    offsetY = 10 - (cardCenterY - popoverHeight / 2);
  }
  // 检查底部溢出
  else if (cardCenterY + popoverHeight / 2 > window.innerHeight - 10) {
    offsetY = (window.innerHeight - 10) - (cardCenterY + popoverHeight / 2);
  }

  popover.style.setProperty('--popover-offset-y', `${offsetY}px`);
}

/**
 * 初始化分数弹出层交互（触摸设备兼容）
 */
function initScorePopover() {
  const scoreCard = document.getElementById('score-card');
  const popover = document.getElementById('score-popover');

  if (!scoreCard || !popover) return;

  // hover 进入时调整位置
  scoreCard.addEventListener('mouseenter', () => {
    requestAnimationFrame(adjustPopoverPosition);
  });

  // 触摸设备：点击切换
  scoreCard.addEventListener('click', (e) => {
    if ('ontouchstart' in window) {
      popover.classList.toggle('visible');
      e.stopPropagation();
      requestAnimationFrame(adjustPopoverPosition);
    }
  });

  // 点击外部关闭
  document.addEventListener('click', (e) => {
    if (!scoreCard.contains(e.target)) {
      popover.classList.remove('visible');
    }
  });

  // 窗口resize时重新计算
  window.addEventListener('resize', () => {
    if (window.getComputedStyle(popover).visibility === 'visible') {
      requestAnimationFrame(adjustPopoverPosition);
    }
  }, { passive: true });
}

function renderOptimizedPrompt() {
  const promptContent = document.getElementById('optimized-prompt');
  const promptSource = document.getElementById('optimized-prompt-source');
  const rawPrompt = inputData.current.optimizedPrompt || '';

  // 渲染模式显示
  promptContent.innerHTML = parseMarkdown(rawPrompt);

  // 源码模式显示
  promptSource.value = rawPrompt;

  // 监听源码编辑
  promptSource.addEventListener('input', () => {
    inputData.current.optimizedPrompt = promptSource.value;
  });
}

// ========== Prompt 模式切换 ==========
function switchPromptMode(mode) {
  currentPromptMode = mode;

  const renderBtn = document.getElementById('mode-render-btn');
  const sourceBtn = document.getElementById('mode-source-btn');
  const promptContent = document.getElementById('optimized-prompt');
  const promptSource = document.getElementById('optimized-prompt-source');

  if (mode === 'render') {
    renderBtn.classList.add('active');
    sourceBtn.classList.remove('active');
    promptContent.classList.remove('hidden');
    promptSource.classList.add('hidden');
    // 从源码更新渲染
    promptContent.innerHTML = parseMarkdown(promptSource.value);
  } else {
    renderBtn.classList.remove('active');
    sourceBtn.classList.add('active');
    promptContent.classList.add('hidden');
    promptSource.classList.remove('hidden');
  }
}

function renderDirections() {
  const container = document.getElementById('direction-cards');
  const directions = inputData.current.suggestedDirections || [];

  // 定义颜色循环
  const colors = ['blue', 'pink', 'green', 'orange', 'purple'];

  // 使用胶囊标签样式 (Bento Playful)
  container.innerHTML = directions.map((dir, index) => `
    <button class="tag" data-id="${dir.id}" data-color="${colors[index % colors.length]}" onclick="toggleDirection('${dir.id}')" title="${escapeHtml(dir.description)}">
      <svg class="tag-check" width="14" height="14" viewBox="0 0 14 14" fill="none">
        <path d="M3 7L6 10L11 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <span class="tag-label">${escapeHtml(dir.label)}</span>
    </button>
  `).join('');
}

function toggleDirection(id) {
  const tag = document.querySelector(`.tag[data-id="${id}"]`);
  if (selectedDirections.has(id)) {
    selectedDirections.delete(id);
    tag.classList.remove('selected');
  } else {
    selectedDirections.add(id);
    tag.classList.add('selected');
  }
}

// ========== 历史记录 ==========
function renderHistory() {
  const container = document.getElementById('history-dropdown');
  const history = inputData.history || [];
  const currentVersion = getCurrentIterationVersion();

  let html = '';

  // 当前版本 - 点击直接切换
  const isViewingCurrent = viewingHistoryVersion === null;
  html += `
    <div class="history-item current${isViewingCurrent ? ' viewing' : ''}" onclick="switchToVersion(null)">
      <div class="history-item-header">
        <div class="history-indicator"></div>
        <span class="history-version">v${currentVersion} (${t('current')})</span>
        <span class="history-score">- ${inputData.current.score}${t('points')}</span>
      </div>
      <div class="history-feedback">${t('currentPending')}</div>
    </div>
  `;

  // 历史版本 (倒序显示) - 点击直接切换
  for (let i = history.length - 1; i >= 0; i--) {
    const item = history[i];
    const version = i + 1;
    const feedback = item.userFeedback;
    const isViewing = viewingHistoryVersion === item.iterationId;

    html += `
      <div class="history-item${isViewing ? ' viewing' : ''}" data-iteration-id="${item.iterationId}" onclick="switchToVersion('${item.iterationId}')">
        <div class="history-item-header">
          <div class="history-indicator"></div>
          <span class="history-version">v${version}</span>
          <span class="history-score">- ${item.score}${t('points')}</span>
        </div>
        ${feedback ? `
          <div class="history-feedback">
            ${feedback.selectedDirections && feedback.selectedDirections.length > 0
              ? `${t('selectDirection')}: ${feedback.selectedDirections.join(', ')}<br>`
              : ''}
            ${feedback.userInput ? `${t('userInput')}: "${escapeHtml(feedback.userInput)}"` : ''}
          </div>
        ` : ''}
      </div>
    `;
  }

  container.innerHTML = html;
}

function toggleHistoryDropdown() {
  const dropdown = document.getElementById('history-dropdown');
  dropdown.classList.toggle('hidden');
}

// 统一版本切换函数
function switchToVersion(iterationId) {
  viewingHistoryVersion = iterationId;

  if (iterationId === null) {
    // 切换回当前版本
    renderScore(inputData.current.score);
    document.getElementById('review-report').innerHTML = parseMarkdown(inputData.current.reviewReport || t('noReviewReport'));

    // 渲染评估报告到弹出层
    renderEvaluationPopover();

    document.getElementById('optimized-prompt').innerHTML = parseMarkdown(inputData.current.optimizedPrompt || '');
    document.getElementById('optimized-prompt-source').value = inputData.current.optimizedPrompt || '';
  } else {
    // 切换到历史版本
    const item = findHistoryItem(iterationId);
    if (item) {
      renderScore(item.score);
      document.getElementById('review-report').innerHTML = parseMarkdown(item.reviewReport || '');

      // 渲染历史版本的评估报告到弹出层
      renderEvaluationPopoverForHistory(item);

      document.getElementById('optimized-prompt').innerHTML = parseMarkdown(item.optimizedPrompt || '');
      document.getElementById('optimized-prompt-source').value = item.optimizedPrompt || '';
    }
  }

  // 更新版本指示器
  updateVersionIndicator();
  // 更新历史记录中的选中状态
  renderHistory();
  // 关闭下拉框
  toggleHistoryDropdown();
}

/**
 * 渲染历史版本的评估报告到弹出层
 */
function renderEvaluationPopoverForHistory(historyItem) {
  const evaluationEl = document.getElementById('evaluation-report');
  if (!evaluationEl) return;

  const evalMarkdown = historyItem.evaluationReport || '';
  if (!evalMarkdown) {
    evaluationEl.innerHTML = '<p style="color: var(--text-tertiary); font-size: 12px;">' + (t('noEvaluationReport') || '暂无评估数据') + '</p>';
    return;
  }

  const dimensions = parseEvalDimensions(evalMarkdown);
  const advantages = parseEvalAdvantages(evalMarkdown);
  const suggestions = parseEvalSuggestions(evalMarkdown);
  const direction = parseEvalDirection(evalMarkdown);

  evaluationEl.innerHTML =
    renderPopoverDimensions(dimensions) +
    renderPopoverAdvantages(advantages) +
    renderPopoverSuggestions(suggestions) +
    renderPopoverDirection(direction);
}

// 更新版本状态指示器
function updateVersionIndicator() {
  const indicator = document.getElementById('history-viewing-indicator');
  const versionNumber = document.getElementById('viewing-version-number');
  const continueBtn = document.getElementById('continue-from-history-btn');

  if (viewingHistoryVersion === null) {
    // 查看当前版本，隐藏指示器和按钮
    indicator.classList.add('hidden');
    continueBtn.classList.add('hidden');
  } else {
    // 查看历史版本，显示指示器和按钮
    const history = inputData.history || [];
    const versionIndex = history.findIndex(item => item.iterationId === viewingHistoryVersion);
    const version = versionIndex + 1;

    versionNumber.textContent = `v${version}`;
    indicator.classList.remove('hidden');
    continueBtn.classList.remove('hidden');
  }
}

// 基于当前查看的历史版本继续优化
function continueFromViewingVersion() {
  if (viewingHistoryVersion) {
    rollbackToVersion(viewingHistoryVersion);
  }
}

function rollbackToVersion(iterationId) {
  toggleHistoryDropdown();
  showSuccessAnimation();

  setTimeout(async () => {
    const directions = Array.from(selectedDirections);
    const userInput = document.getElementById('user-input').value;

    try {
      await window.go.main.App.Rollback(iterationId, directions, userInput);
    } catch (error) {
      console.error('Rollback failed:', error);
    }
  }, 500);
}

// 包装 utils.js 中的纯函数
function findHistoryItem(iterationId) {
  return findHistoryItemInArray(inputData.history || [], iterationId);
}

// ========== 版本对比 ==========
function initCompareSelectors() {
  const leftSelect = document.getElementById('compare-left');
  const rightSelect = document.getElementById('compare-right');
  const history = inputData.history || [];
  const currentVersion = getCurrentIterationVersion();

  let options = '';

  // 添加当前版本
  options += `<option value="current">v${currentVersion} (${t('current')}) - ${inputData.current.score}${t('points')}</option>`;

  // 添加历史版本
  for (let i = history.length - 1; i >= 0; i--) {
    const item = history[i];
    const version = i + 1;
    options += `<option value="${item.iterationId}">v${version} - ${item.score}${t('points')}</option>`;
  }

  leftSelect.innerHTML = options;
  rightSelect.innerHTML = options;

  // 默认选择第一个历史版本和当前版本
  if (history.length > 0) {
    leftSelect.value = history[history.length - 1].iterationId;
    rightSelect.value = 'current';
  }
}

function showCompareView() {
  document.getElementById('compare-view').classList.remove('hidden');

  // 如果正在查看历史版本，自动设置对比：左侧=查看的版本，右侧=当前版本
  if (viewingHistoryVersion !== null) {
    const leftSelect = document.getElementById('compare-left');
    const rightSelect = document.getElementById('compare-right');
    leftSelect.value = viewingHistoryVersion;
    rightSelect.value = 'current';
  }

  updateComparison();
}

function hideCompareView() {
  document.getElementById('compare-view').classList.add('hidden');
}

function updateComparison() {
  const leftId = document.getElementById('compare-left').value;
  const rightId = document.getElementById('compare-right').value;

  const leftData = getVersionData(leftId);
  const rightData = getVersionData(rightId);

  // 更新统一 Diff 视图标签
  document.getElementById('compare-unified-left-label').textContent =
    `v${leftData.version} (${leftData.score}${t('points')})`;
  document.getElementById('compare-unified-right-label').textContent =
    `v${rightData.version} (${rightData.score}${t('points')})`;

  // 生成 Diff 高亮内容
  const diffHtml = generateDiffHtml(leftData.optimizedPrompt, rightData.optimizedPrompt);
  document.getElementById('compare-unified-content').innerHTML = diffHtml;

  // 更新并排视图
  document.getElementById('compare-left-header').textContent =
    `v${leftData.version} (${leftData.score}${t('points')})`;
  document.getElementById('compare-left-content').innerHTML =
    generateSideDiffHtml(leftData.optimizedPrompt, rightData.optimizedPrompt, 'left');

  document.getElementById('compare-right-header').textContent =
    `v${rightData.version} (${rightData.score}${t('points')})`;
  document.getElementById('compare-right-content').innerHTML =
    generateSideDiffHtml(leftData.optimizedPrompt, rightData.optimizedPrompt, 'right');

  // 更新回滚按钮
  const rollbackBtn = document.getElementById('rollback-from-compare-btn');
  if (leftId === 'current') {
    rollbackBtn.style.display = 'none';
  } else {
    rollbackBtn.style.display = 'block';
    rollbackBtn.textContent = t('continueFromVersionN', { version: leftData.version });
    rollbackBtn.onclick = () => {
      hideCompareView();
      rollbackToVersion(leftId);
    };
  }
}

// ========== Diff 高亮生成 ==========
function generateDiffHtml(oldText, newText) {
  if (typeof Diff === 'undefined') {
    // Diff 库未加载，显示简单对比
    return `<div class="diff-unchanged">${escapeHtml(newText)}</div>`;
  }

  // 使用 diffWords 进行词级别对比
  const diff = Diff.diffWords(oldText || '', newText || '');
  let html = '';

  diff.forEach(part => {
    const escaped = escapeHtml(part.value);
    if (part.added) {
      html += `<span class="diff-word-added">${escaped}</span>`;
    } else if (part.removed) {
      html += `<span class="diff-word-removed">${escaped}</span>`;
    } else {
      html += escaped;
    }
  });

  return `<div class="diff-content">${html}</div>`;
}

function generateSideDiffHtml(oldText, newText, side) {
  if (typeof Diff === 'undefined') {
    return `<div>${escapeHtml(side === 'left' ? oldText : newText)}</div>`;
  }

  // 使用 diffLines 进行行级别对比
  const diff = Diff.diffLines(oldText || '', newText || '');
  let html = '';

  diff.forEach(part => {
    const lines = part.value.split('\n').filter(line => line !== '');
    lines.forEach(line => {
      const escaped = escapeHtml(line);
      if (part.added) {
        if (side === 'right') {
          html += `<div class="diff-line diff-added">${escaped}</div>`;
        }
      } else if (part.removed) {
        if (side === 'left') {
          html += `<div class="diff-line diff-removed">${escaped}</div>`;
        }
      } else {
        html += `<div class="diff-line diff-unchanged">${escaped}</div>`;
      }
    });
  });

  return html;
}

// ========== Diff 模式切换 ==========
function switchDiffMode(mode) {
  currentDiffMode = mode;

  const unifiedBtn = document.getElementById('diff-unified-btn');
  const splitBtn = document.getElementById('diff-split-btn');
  const unifiedView = document.getElementById('compare-unified');
  const splitView = document.getElementById('compare-split');

  if (mode === 'unified') {
    unifiedBtn.classList.add('active');
    splitBtn.classList.remove('active');
    unifiedView.classList.remove('hidden');
    splitView.classList.add('hidden');
  } else {
    unifiedBtn.classList.remove('active');
    splitBtn.classList.add('active');
    unifiedView.classList.add('hidden');
    splitView.classList.remove('hidden');
  }
}

// 包装 utils.js 中的纯函数，适配 getCurrentIterationVersion
function getVersionData(id) {
  const adaptedData = {
    version: getCurrentIterationVersion(),
    current: inputData.current,
    history: inputData.history
  };
  return getVersionDataFromInput(adaptedData, id);
}

function rollbackFromCompare() {
  const leftId = document.getElementById('compare-left').value;
  if (leftId !== 'current') {
    hideCompareView();
    rollbackToVersion(leftId);
  }
}

// ========== 原始 Prompt 弹窗 ==========
function showOriginalPrompt() {
  const modal = document.getElementById('original-prompt-modal');
  const content = document.getElementById('original-prompt-content');
  content.textContent = inputData.originalPrompt || '';
  modal.classList.remove('hidden');
}

function hideOriginalPrompt() {
  document.getElementById('original-prompt-modal').classList.add('hidden');
}

// ========== 倒计时 ==========
function startCountdown() {
  updateCountdown();
  countdownInterval = setInterval(updateCountdown, 1000);
}

async function updateCountdown() {
  try {
    const remaining = await window.go.main.App.GetRemainingSeconds();
    const timeout = await window.go.main.App.GetTimeoutSeconds();
    const minutes = Math.floor(remaining / 60);
    const seconds = remaining % 60;

    const display = document.getElementById('countdown');
    display.textContent = `${minutes}:${seconds.toString().padStart(2, '0')}`;

    // 更新进度条
    const progressFill = document.getElementById('progress-fill');
    if (progressFill && timeout > 0) {
      const percentage = (remaining / timeout) * 100;
      progressFill.style.width = `${percentage}%`;

      // 颜色警告
      if (remaining < 30) {
        progressFill.style.background = 'linear-gradient(90deg, #EF4444, #F87171)';
      } else if (remaining < 60) {
        progressFill.style.background = 'linear-gradient(90deg, #F97316, #FB923C)';
      } else {
        progressFill.style.background = 'linear-gradient(90deg, var(--color-accent), #60A5FA)';
      }
    }

    // 更新警告状态
    display.classList.remove('warning', 'critical');
    if (remaining < 30) {
      display.classList.add('critical');
    } else if (remaining < 60) {
      display.classList.add('warning');
    }

    if (remaining <= 0) {
      clearInterval(countdownInterval);
    }
  } catch (error) {
    console.error('Failed to get remaining seconds:', error);
  }
}

// ========== 用户输入 ==========
function initUserInput() {
  const textarea = document.getElementById('user-input');
  const charCount = document.getElementById('char-count');

  textarea.addEventListener('input', () => {
    const length = textarea.value.length;
    charCount.textContent = `${length} / 2000`;
  });
}

// ========== 复制功能 ==========
async function copyPrompt() {
  // 从源码编辑器获取最新内容
  const promptSource = document.getElementById('optimized-prompt-source');
  const prompt = promptSource ? promptSource.value : (inputData.current.optimizedPrompt || '');

  try {
    await navigator.clipboard.writeText(prompt);

    const btn = document.querySelector('.copy-btn');
    const textSpan = btn.querySelector('.copy-text');
    const originalText = textSpan.textContent;

    btn.classList.add('copied');
    textSpan.textContent = t('copied');

    setTimeout(() => {
      btn.classList.remove('copied');
      textSpan.textContent = originalText;
    }, 2000);
  } catch (error) {
    console.error('Failed to copy:', error);
  }
}

// ========== 提交操作 ==========
async function handleSubmit() {
  const submitBtn = document.getElementById('submit-btn');
  submitBtn.disabled = true;

  showSuccessAnimation();

  setTimeout(async () => {
    const directions = Array.from(selectedDirections);
    const userInput = document.getElementById('user-input').value;

    try {
      await window.go.main.App.Submit(directions, userInput);
    } catch (error) {
      console.error('Submit failed:', error);
      submitBtn.disabled = false;
      hideSuccessAnimation();
    }
  }, 500);
}

async function handleCancel() {
  try {
    await window.go.main.App.Cancel();
  } catch (error) {
    console.error('Cancel failed:', error);
  }
}

// ========== 成功动画 ==========
function showSuccessAnimation() {
  document.getElementById('success-overlay').classList.remove('hidden');
}

function hideSuccessAnimation() {
  document.getElementById('success-overlay').classList.add('hidden');
}

// ========== Markdown 解析（使用 markdown-it 或回退解析器）==========
function parseMarkdown(text) {
  if (!text) return '';

  // 使用 markdown-it 库进行渲染（如果可用）
  if (typeof markdownit !== 'undefined') {
    try {
      const md = markdownit({
        html: false,        // 禁用 HTML 标签
        breaks: false,      // 不将单个换行符转换为 <br>
        linkify: true,      // 自动识别链接
        typographer: false  // 禁用排版优化
      });
      return md.render(text);
    } catch (e) {
      console.warn('markdown-it parse error, falling back to simple parser:', e);
    }
  } else {
    console.warn('markdown-it not loaded, using fallback parser');
  }

  // 回退到简单解析器
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

  return html;
}

// ========== 点击外部关闭下拉 ==========
document.addEventListener('click', (e) => {
  const historyWrapper = document.querySelector('.history-dropdown-wrapper');
  const historyDropdown = document.getElementById('history-dropdown');

  if (historyWrapper && !historyWrapper.contains(e.target)) {
    historyDropdown.classList.add('hidden');
  }
});

// ========== 暴露全局函数 (供 HTML onclick 调用) ==========
window.toggleHistoryDropdown = toggleHistoryDropdown;
window.showCompareView = showCompareView;
window.hideCompareView = hideCompareView;
window.copyPrompt = copyPrompt;
window.showOriginalPrompt = showOriginalPrompt;
window.hideOriginalPrompt = hideOriginalPrompt;
window.handleSubmit = handleSubmit;
window.handleCancel = handleCancel;
window.toggleDirection = toggleDirection;
window.switchToVersion = switchToVersion;
window.continueFromViewingVersion = continueFromViewingVersion;
window.rollbackToVersion = rollbackToVersion;
window.rollbackFromCompare = rollbackFromCompare;
window.updateComparison = updateComparison;
window.switchPromptMode = switchPromptMode;
window.switchDiffMode = switchDiffMode;
