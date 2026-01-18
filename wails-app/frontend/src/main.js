// ========== 全局状态 ==========
let inputData = null;
let selectedDirections = new Set();
let countdownInterval = null;
let currentTab = 'review';

// ========== 初始化 ==========
document.addEventListener('DOMContentLoaded', async () => {
  try {
    await loadData();
  } catch (error) {
    console.error('Failed to load data:', error);
    showError('加载数据失败: ' + error.message);
  }
});

// ========== 数据加载 ==========
async function loadData() {
  try {
    inputData = await window.go.main.App.GetInputData();
    hideLoading();
    renderUI();
    startCountdown();
    initTabIndicator();
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
}

function renderVersion() {
  const versionBadge = document.getElementById('version-badge');
  versionBadge.textContent = 'v' + inputData.version;
}

function renderScore(score) {
  // 更新评分值
  const scoreValue = document.getElementById('score-value');
  scoreValue.textContent = score;

  // 更新评分环
  const scoreRingProgress = document.getElementById('score-ring-progress');
  const circumference = 2 * Math.PI * 70; // r = 70
  const offset = circumference - (score / 100) * circumference;
  scoreRingProgress.style.strokeDasharray = circumference;
  scoreRingProgress.style.strokeDashoffset = offset;

  // 设置颜色等级
  const gradeClass = getScoreGradeClass(score);
  scoreRingProgress.className = 'score-ring-progress ' + gradeClass;

  // 更新等级文字
  const scoreGrade = document.getElementById('score-grade');
  scoreGrade.textContent = getScoreGradeText(score);

  // 更新进度条
  const scoreBarFill = document.getElementById('score-bar-fill');
  scoreBarFill.style.width = score + '%';
  scoreBarFill.className = 'score-bar-fill ' + gradeClass;

  const scoreBarValue = document.getElementById('score-bar-value');
  scoreBarValue.textContent = score + '/100';
}

function getScoreGradeClass(score) {
  if (score <= 40) return 'poor';
  if (score <= 60) return 'fair';
  if (score <= 80) return 'good';
  return 'excellent';
}

function getScoreGradeText(score) {
  if (score <= 40) return 'Poor';
  if (score <= 60) return 'Fair';
  if (score <= 80) return 'Good';
  return 'Excellent';
}

function renderReports() {
  const reviewReport = document.getElementById('review-report');
  const evaluationReport = document.getElementById('evaluation-report');

  reviewReport.innerHTML = parseMarkdown(inputData.current.reviewReport || '暂无评审报告');
  evaluationReport.innerHTML = parseMarkdown(inputData.current.evaluationReport || '暂无评估报告');
}

function renderOptimizedPrompt() {
  const promptContent = document.getElementById('optimized-prompt');
  promptContent.textContent = inputData.current.optimizedPrompt || '';
}

function renderDirections() {
  const container = document.getElementById('direction-cards');
  const directions = inputData.current.suggestedDirections || [];

  container.innerHTML = directions.map(dir => `
    <div class="direction-card" data-id="${dir.id}" onclick="toggleDirection('${dir.id}')">
      <div class="direction-card-header">
        <div class="direction-checkbox">
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
            <path d="M2.5 6L5 8.5L9.5 4" stroke="white" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <span class="direction-label">${escapeHtml(dir.label)}</span>
      </div>
      <p class="direction-description">${escapeHtml(dir.description)}</p>
    </div>
  `).join('');
}

function toggleDirection(id) {
  const card = document.querySelector(`.direction-card[data-id="${id}"]`);
  if (selectedDirections.has(id)) {
    selectedDirections.delete(id);
    card.classList.remove('selected');
  } else {
    selectedDirections.add(id);
    card.classList.add('selected');
  }
}

// ========== 历史记录 ==========
function renderHistory() {
  const container = document.getElementById('history-dropdown');
  const history = inputData.history || [];
  const currentVersion = inputData.version;

  let html = '';

  // 当前版本
  html += `
    <div class="history-item current">
      <div class="history-item-header">
        <div class="history-indicator"></div>
        <span class="history-version">v${currentVersion} (当前)</span>
        <span class="history-score">- ${inputData.current.score}分</span>
      </div>
      <div class="history-feedback">当前版本，待提交</div>
    </div>
  `;

  // 历史版本 (倒序显示)
  for (let i = history.length - 1; i >= 0; i--) {
    const item = history[i];
    const version = i + 1;
    const feedback = item.userFeedback;

    html += `
      <div class="history-item" data-iteration-id="${item.iterationId}">
        <div class="history-item-header">
          <div class="history-indicator"></div>
          <span class="history-version">v${version}</span>
          <span class="history-score">- ${item.score}分</span>
        </div>
        ${feedback ? `
          <div class="history-feedback">
            ${feedback.selectedDirections && feedback.selectedDirections.length > 0
              ? `选择方向: ${feedback.selectedDirections.join(', ')}<br>`
              : ''}
            ${feedback.userInput ? `用户输入: "${escapeHtml(feedback.userInput)}"` : ''}
          </div>
        ` : ''}
        <div class="history-actions">
          <button class="history-action-btn" onclick="viewHistoryVersion('${item.iterationId}')">查看详情</button>
          <button class="history-action-btn" onclick="rollbackToVersion('${item.iterationId}')">基于此版本继续优化</button>
        </div>
      </div>
    `;
  }

  container.innerHTML = html;
}

function toggleHistoryDropdown() {
  const dropdown = document.getElementById('history-dropdown');
  dropdown.classList.toggle('hidden');
}

function viewHistoryVersion(iterationId) {
  const item = findHistoryItem(iterationId);
  if (item) {
    // 临时显示历史版本内容
    renderScore(item.score);
    document.getElementById('review-report').innerHTML = parseMarkdown(item.reviewReport || '');
    document.getElementById('evaluation-report').innerHTML = parseMarkdown(item.evaluationReport || '');
    document.getElementById('optimized-prompt').textContent = item.optimizedPrompt || '';
    toggleHistoryDropdown();
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

function findHistoryItem(iterationId) {
  return (inputData.history || []).find(item => item.iterationId === iterationId);
}

// ========== 版本对比 ==========
function initCompareSelectors() {
  const leftSelect = document.getElementById('compare-left');
  const rightSelect = document.getElementById('compare-right');
  const history = inputData.history || [];

  let options = '';

  // 添加当前版本
  options += `<option value="current">v${inputData.version} (当前) - ${inputData.current.score}分</option>`;

  // 添加历史版本
  for (let i = history.length - 1; i >= 0; i--) {
    const item = history[i];
    const version = i + 1;
    options += `<option value="${item.iterationId}">v${version} - ${item.score}分</option>`;
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

  // 更新左侧面板
  document.getElementById('compare-left-header').textContent =
    `v${leftData.version} (${leftData.score}分)`;
  document.getElementById('compare-left-content').textContent =
    leftData.optimizedPrompt;

  // 更新右侧面板
  document.getElementById('compare-right-header').textContent =
    `v${rightData.version} (${rightData.score}分)`;
  document.getElementById('compare-right-content').textContent =
    rightData.optimizedPrompt;

  // 更新回滚按钮
  const rollbackBtn = document.getElementById('rollback-from-compare-btn');
  if (leftId === 'current') {
    rollbackBtn.style.display = 'none';
  } else {
    rollbackBtn.style.display = 'block';
    rollbackBtn.textContent = `基于 v${leftData.version} 继续优化`;
    rollbackBtn.onclick = () => {
      hideCompareView();
      rollbackToVersion(leftId);
    };
  }
}

function getVersionData(id) {
  if (id === 'current') {
    return {
      version: inputData.version,
      score: inputData.current.score,
      optimizedPrompt: inputData.current.optimizedPrompt
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

// ========== Tab 切换 ==========
function switchTab(tab) {
  currentTab = tab;

  // 更新按钮状态
  document.querySelectorAll('.segment').forEach(btn => {
    btn.classList.toggle('active', btn.dataset.tab === tab);
  });

  // 更新面板显示
  document.getElementById('review-report').classList.toggle('hidden', tab !== 'review');
  document.getElementById('evaluation-report').classList.toggle('hidden', tab !== 'evaluation');

  // 更新指示器位置
  updateTabIndicator();
}

function initTabIndicator() {
  updateTabIndicator();
}

function updateTabIndicator() {
  const activeBtn = document.querySelector('.segment.active');
  const indicator = document.querySelector('.segment-indicator');

  if (activeBtn && indicator) {
    const container = document.querySelector('.segmented-control');
    const btnRect = activeBtn.getBoundingClientRect();
    const containerRect = container.getBoundingClientRect();

    indicator.style.width = btnRect.width + 'px';
    indicator.style.transform = `translateX(${btnRect.left - containerRect.left - 2}px)`;
  }
}

// ========== 倒计时 ==========
function startCountdown() {
  updateCountdown();
  countdownInterval = setInterval(updateCountdown, 1000);
}

async function updateCountdown() {
  try {
    const remaining = await window.go.main.App.GetRemainingSeconds();
    const minutes = Math.floor(remaining / 60);
    const seconds = remaining % 60;

    const display = document.getElementById('countdown');
    display.textContent = `${minutes}:${seconds.toString().padStart(2, '0')}`;

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
  const prompt = inputData.current.optimizedPrompt || '';

  try {
    await navigator.clipboard.writeText(prompt);

    const btn = document.querySelector('.copy-btn');
    const textSpan = btn.querySelector('.copy-text');
    const originalText = textSpan.textContent;

    btn.classList.add('copied');
    textSpan.textContent = '已复制';

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

// ========== 工具函数 ==========
function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

function parseMarkdown(text) {
  if (!text) return '';

  // 简单的 Markdown 解析
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

// ========== 点击外部关闭下拉 ==========
document.addEventListener('click', (e) => {
  const historyWrapper = document.querySelector('.history-dropdown-wrapper');
  const historyDropdown = document.getElementById('history-dropdown');

  if (historyWrapper && !historyWrapper.contains(e.target)) {
    historyDropdown.classList.add('hidden');
  }
});
