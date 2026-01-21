# Prompt Optimizer Skill - 前端实现方案

## 设计理念

| 属性 | 值 |
|------|-----|
| 风格 | macOS Sequoia/Sonoma 原生应用风格 |
| 核心体验 | 物理感、深度感、高级半透明质感 |
| 设计哲学 | "System-Level Realism" — 拒绝扁平 Web 设计 |
| 最大宽度 | 1200px |

---

## 目录结构

```
webui/
├── templates/
│   └── optimize.html       # Jinja2 模板
└── static/
    ├── css/
    │   └── styles.css
    └── js/
        └── app.js
```

---

## 页面布局

### 主界面布局

```
┌─────────────────────────────────────────────────────────────────┐
│ ● ● ●  Prompt Optimizer                                  [9:45] │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  [当前版本 v3]  [历史记录 ▼]  [版本对比]                          │
│                                                                 │
├────────────────────────────┬────────────────────────────────────┤
│                            │                                    │
│     评分环形图 (SVG)        │     评审报告 + 评估报告 (Tab)       │
│     Score: 85/100          │     [Review | Evaluation]          │
│     Grade: Excellent       │     (Markdown, 可滚动, 320px)      │
│                            │                                    │
├────────────────────────────┴────────────────────────────────────┤
│                                                                 │
│     优化后的提示词 (Markdown + 代码高亮)                         │
│     [复制按钮]                                                  │
│     (max-height: 400px, 可滚动)                                 │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  评分: ████████░░ 85/100                                        │
│                                                                 │
│  原始 Prompt: [查看原始输入]                                     │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│ 优化方向选择 (动态渲染 3-4 个卡片)                                │
│ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐             │
│ │ ☐ 添加示例   │ │ ☐ 增强约束   │ │ ☐ 优化结构   │             │
│ │ 补充使用案例 │ │ 明确边界条件 │ │ 改进段落组织 │             │
│ └──────────────┘ └──────────────┘ └──────────────┘             │
├─────────────────────────────────────────────────────────────────┤
│ 补充说明 (可选)                                    [0 / 2000]   │
│ ┌───────────────────────────────────────────────────────────┐  │
│ │                                                           │  │
│ └───────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                                         [ 取消 ]  [ 确定 ]     │
└─────────────────────────────────────────────────────────────────┘
```

### 历史记录下拉面板

```
┌─ 历史记录 ─────────────────────────────────────────────────────┐
│                                                                │
│  ● v3 (当前) - 85分                                            │
│    └─ (当前版本，待提交)                                        │
│                                                                │
│  ○ v2 - 75分                                                   │
│    ├─ 选择方向: 添加示例, 增强约束                              │
│    └─ 用户输入: "需要更多示例和约束"                            │
│    [查看详情]  [基于此版本继续优化]                             │
│                                                                │
│  ○ v1 - 65分                                                   │
│    ├─ 选择方向: 优化结构                                       │
│    └─ 用户输入: "希望结构更清晰"                                │
│    [查看详情]  [基于此版本继续优化]                             │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

### 版本对比视图

```
┌─────────────────────────────────────────────────────────────────┐
│  版本对比                                          [← 返回]     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  对比: [v1 ▼]  ↔  [v3 ▼]                                        │
│                                                                 │
├────────────────────────────┬────────────────────────────────────┤
│         v1 (65分)          │           v3 (85分)                │
├────────────────────────────┼────────────────────────────────────┤
│                            │                                    │
│ # Role: 助手               │ # Role: 专家                       │
│                            │                                    │
│ 你是一个AI助手...          │ ## Profile                         │
│                            │ - Author: Expert                   │
│                            │ - Version: 2.0                     │
│                            │                                    │
│                            │ ## Background                      │
│                            │ 你是一个专业的...                   │
│                            │                                    │
└────────────────────────────┴────────────────────────────────────┘
│                                                                 │
│                        [基于 v1 继续优化]                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

# 第一部分：视觉物理引擎

## 1. 高级毛玻璃效果 (Vibrancy 3.0)

### 材质公式

```css
backdrop-filter: blur(VAR) saturate(VAR);
background: rgba(r, g, b, opacity);
```

### 材质类型

| 材质 | 用途 | Light Mode | Dark Mode | Blur | Saturate |
|------|------|------------|-----------|------|----------|
| Thin | 侧边栏 | `bg-gray-100/60` | `bg-[#1e1e1e]/60` | `blur-2xl` | `saturate-150` |
| Thick | 主窗口 | `bg-white/80` | `bg-[#282828]/70` | `blur-3xl` | `saturate-180` |
| Ultra | 弹出层 | `bg-white/90` | `bg-[#323232]/90` | `blur-xl` | `saturate-200` |

### CSS 变量

```css
:root {
  /* 材质 - Light */
  --material-thin: rgba(246, 246, 246, 0.6);
  --material-thick: rgba(255, 255, 255, 0.8);
  --material-ultra: rgba(255, 255, 255, 0.9);

  /* 材质 - Dark */
  --material-thin-dark: rgba(30, 30, 30, 0.6);
  --material-thick-dark: rgba(40, 40, 40, 0.7);
  --material-ultra-dark: rgba(50, 50, 50, 0.9);

  /* 模糊值 */
  --blur-thin: 40px;
  --blur-thick: 64px;
  --blur-ultra: 24px;
}
```

### 噪点纹理 (必需)

```css
.noise-overlay::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image: url("data:image/svg+xml,...");
  opacity: 0.015;
  pointer-events: none;
  mix-blend-mode: overlay;
}
```

---

## 2. 光照与 Retina 边框

### 0.5px 规则

| 模式 | Tailwind | CSS |
|------|----------|-----|
| Light | `border-black/5` | `box-shadow: 0 0 0 0.5px rgba(0,0,0,0.05)` |
| Dark | `border-white/10` | `box-shadow: 0 0 0 0.5px rgba(255,255,255,0.1)` |

### 顶部高光 (The Bezel)

```css
/* 每个浮动容器必须有 */
box-shadow: inset 0 1px 0 0 rgba(255, 255, 255, 0.4);

/* Dark Mode */
box-shadow: inset 0 1px 0 0 rgba(255, 255, 255, 0.1);
```

---

## 3. 阴影系统

```css
:root {
  /* 窗口阴影 */
  --shadow-window: 0px 0px 1px rgba(0,0,0,0.4), 0px 16px 36px -8px rgba(0,0,0,0.2);

  /* 卡片阴影 */
  --shadow-card: 0px 0px 0.5px rgba(0,0,0,0.1), 0px 4px 12px -2px rgba(0,0,0,0.08);

  /* 弹出层阴影 */
  --shadow-popover: 0px 0px 1px rgba(0,0,0,0.3), 0px 24px 48px -12px rgba(0,0,0,0.25);

  /* 按钮阴影 */
  --shadow-button: 0px 0.5px 1px rgba(0,0,0,0.1), 0px 1px 2px rgba(0,0,0,0.06);
}
```

---

# 第二部分：字体与图标

## 字体系统

```css
:root {
  --font-system: -apple-system, BlinkMacSystemFont, "SF Pro Text", "Inter", sans-serif;
  --font-mono: "SF Mono", "JetBrains Mono", ui-monospace, monospace;
}

* {
  -webkit-font-smoothing: antialiased;
}
```

## 字体层级

| 用途 | 大小 | 字重 | 字距 |
|------|------|------|------|
| 窗口标题 | 13px | 600 | normal |
| 区块标题 | 15px | 600 | -0.01em |
| 正文 | 13px | 400 | normal |
| 辅助文字 | 11px | 400 | 0.02em |
| 代码 | 12px | 400 | normal |
| 评分数字 | 28px | 700 | -0.02em |

## 字距规则 (Letter Spacing)

根据 Apple HIG，字距应按字体大小动态调整：

| 字体大小 | Tailwind | CSS | 说明 |
|----------|----------|-----|------|
| < 14px | `tracking-wide` | `letter-spacing: 0.025em` | 小字放宽，提升可读性 |
| 14-20px | `tracking-normal` | `letter-spacing: 0` | 正常字距 |
| > 20px | `tracking-tight` | `letter-spacing: -0.02em` | 大字紧凑，Display 风格 |

```css
/* 实现示例 */
.text-xs, .text-sm { letter-spacing: 0.025em; }
.text-base { letter-spacing: 0; }
.text-xl, .text-2xl, .text-3xl { letter-spacing: -0.02em; }
```

## 图标规范

- **图标库**: Lucide React / Heroicons
- **线条粗细**: 1.5px
- **尺寸**: 按钮内 16-18px

---

# 第三部分：组件规范

## 1. 窗口头部

```
● ● ●  Prompt Optimizer                                  [9:45]
```

### Traffic Lights

| 属性 | 值 |
|------|-----|
| 尺寸 | 12px |
| 间距 | 8px |
| 关闭 | `#FF5F57` |
| 最小化 | `#FFBD2E` |
| 最大化 | `#28C840` |

### 头部样式

| 属性 | Light | Dark |
|------|-------|------|
| 高度 | 52px | 52px |
| 背景 | `bg-white/80` | `bg-[#282828]/70` |
| 模糊 | `backdrop-blur-3xl` | `backdrop-blur-3xl` |
| 底边框 | `border-black/5` | `border-white/8` |

### 倒计时显示

位于头部右侧，显示剩余超时时间。

| 属性 | 值 |
|------|-----|
| 字体 | SF Mono / monospace |
| 大小 | 13px |
| 颜色 | `text-gray-500` |
| 警告颜色 | `text-orange-500` (< 60s) |
| 危急颜色 | `text-red-500` (< 30s) |

```css
.countdown {
  font-family: var(--font-mono);
  font-size: 13px;
  font-variant-numeric: tabular-nums;
  color: #6B7280;
  transition: color 0.3s;
}

.countdown.warning {
  color: #FF9500;
}

.countdown.critical {
  color: #FF3B30;
  animation: pulse 1s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
```

---

## 2. 评分环形图

| 属性 | 值 |
|------|-----|
| 尺寸 | 160px × 160px |
| 环宽 | 10px |
| 端点 | `stroke-linecap: round` |
| 动画 | spring, 1s |

### 评分颜色

| 分数 | 等级 | 颜色 |
|------|------|------|
| 0-40 | Poor | `#FF3B30` |
| 41-60 | Fair | `#FF9500` |
| 61-80 | Good | `#007AFF` |
| 81-100 | Excellent | `#34C759` |

---

## 3. 分段控件 (Tab)

### 容器

```css
.segmented-control {
  display: inline-flex;
  padding: 2px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 8px;
}
```

### Tab 状态

| 状态 | 背景 | 文字 | 阴影 |
|------|------|------|------|
| 非活跃 | 透明 | `text-gray-600` | 无 |
| 悬浮 | `bg-black/3` | `text-gray-700` | 无 |
| 活跃 | `bg-white` | `text-gray-900` | `shadow-sm` |
| 活跃(Dark) | `bg-gray-600` | `text-white` | `shadow-sm` |

### 滑动动画

使用 Framer Motion `layoutId` 实现平滑过渡：
- spring stiffness: 400
- spring damping: 30

---

## 4. 提示词展示区

### 容器

| 属性 | Light | Dark |
|------|-------|------|
| 背景 | `bg-white/90` | `bg-[#323232]/90` |
| 模糊 | `backdrop-blur-xl` | `backdrop-blur-xl` |
| 圆角 | 12px | 12px |
| 内阴影 | `inset 0 1px rgba(255,255,255,0.4)` | `inset 0 1px rgba(255,255,255,0.1)` |

### 复制按钮

```css
.copy-btn {
  padding: 6px;
  border-radius: 6px;
  color: var(--gray-500);
  transition: all 0.15s;
}
.copy-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  color: var(--gray-700);
}
.copy-btn:active {
  transform: scale(0.95);
}
```

---

## 5. 方向卡片 (Bento Grid)

### 布局

```css
.direction-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}
```

### 卡片样式

```css
.direction-card {
  background: var(--material-ultra);
  backdrop-filter: blur(16px) saturate(1.8);
  border-radius: 12px;
  padding: 14px;
  box-shadow:
    inset 0 1px 0 0 rgba(255,255,255,0.4),
    0 0 0 0.5px rgba(0,0,0,0.05),
    0 2px 8px -2px rgba(0,0,0,0.06);
  transition: all 0.2s ease-out;
}
```

### 卡片状态

| 状态 | 背景 | 边框 | 变换 |
|------|------|------|------|
| 默认 | `bg-white/50` | 发丝阴影 | none |
| 悬浮 | `bg-white/70` | 增强阴影 | `translateY(-2px) scale(1.02)` |
| 选中 | `bg-[#007AFF]/10` | `ring-2 ring-[#007AFF]` | none |

### 复选框

| 状态 | 样式 |
|------|------|
| 未选 | 18×18px, 圆角4px, 边框1.5px `rgba(0,0,0,0.2)` |
| 选中 | 背景 `#007AFF`, 白色勾 |

---

## 6. 输入框

### 样式

```css
.textarea {
  background: white;
  border-radius: 8px;
  padding: 12px 14px;
  box-shadow:
    inset 0 1px 2px rgba(0,0,0,0.06),
    0 0 0 0.5px rgba(0,0,0,0.08);
}

.textarea:focus {
  box-shadow:
    inset 0 1px 2px rgba(0,0,0,0.06),
    0 0 0 0.5px rgba(0,122,255,0.5),
    0 0 0 4px rgba(0,122,255,0.15);
}
```

---

## 7. Switch 开关 (Toggle)

Apple "Capsule" 风格的开关组件。

### 尺寸规格

| 属性 | 值 |
|------|-----|
| 宽度 | 26px |
| 高度 | 16px |
| 圆角 | 8px (完全圆角) |
| 滑块尺寸 | 12px × 12px |
| 滑块边距 | 2px |

### 样式

```css
.switch {
  width: 26px;
  height: 16px;
  border-radius: 8px;
  background: rgba(120, 120, 128, 0.16);
  position: relative;
  transition: background 0.2s ease-out;
}

.switch.active {
  background: #34C759;
}

.switch-thumb {
  width: 12px;
  height: 12px;
  border-radius: 6px;
  background: white;
  position: absolute;
  top: 2px;
  left: 2px;
  box-shadow:
    0 0 0 0.5px rgba(0, 0, 0, 0.04),
    0 3px 8px rgba(0, 0, 0, 0.15),
    0 3px 1px rgba(0, 0, 0, 0.06);
}

.switch.active .switch-thumb {
  transform: translateX(10px);
}
```

### 动画

使用 spring 物理动画实现平滑切换：

```javascript
const switchSpring = {
  type: "spring",
  stiffness: 500,
  damping: 30
};
```

### 状态

| 状态 | 背景色 | 滑块位置 |
|------|--------|----------|
| 关闭 | `rgba(120, 120, 128, 0.16)` | `left: 2px` |
| 开启 | `#34C759` (绿色) | `left: 12px` |
| 禁用-关 | `rgba(120, 120, 128, 0.08)` | `left: 2px`, opacity 0.5 |
| 禁用-开 | `rgba(52, 199, 89, 0.5)` | `left: 12px`, opacity 0.5 |

---

## 8. 按钮系统

### 主按钮

```css
.btn-primary {
  height: 32px;
  padding: 0 16px;
  border-radius: 6px;
  background: linear-gradient(180deg, #3B82F6, #2563EB);
  color: white;
  box-shadow:
    inset 0 0.5px 0 0 rgba(255,255,255,0.2),
    0 1px 2px rgba(0,0,0,0.1);
}
.btn-primary:active {
  transform: scale(0.96);
}
```

### 次按钮

```css
.btn-secondary {
  height: 32px;
  padding: 0 16px;
  border-radius: 6px;
  background: rgba(0, 0, 0, 0.05);
  color: #374151;
  box-shadow:
    inset 0 0.5px 0 0 rgba(255,255,255,0.8),
    0 0 0 0.5px rgba(0,0,0,0.1);
}
```

---

## 9. Context Menu 右键菜单

macOS 风格的上下文菜单/下拉菜单。

### 容器样式

```css
.context-menu {
  min-width: 180px;
  padding: 4px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(24px) saturate(1.8);
  border-radius: 10px;
  border: 0.5px solid rgba(0, 0, 0, 0.1);
  box-shadow:
    0 0 0 0.5px rgba(0, 0, 0, 0.05),
    0 10px 38px -10px rgba(22, 23, 24, 0.35),
    0 10px 20px -15px rgba(22, 23, 24, 0.2);
}

/* Dark Mode */
.context-menu.dark {
  background: rgba(50, 50, 50, 0.9);
  border: 0.5px solid rgba(255, 255, 255, 0.1);
}
```

### 菜单项

```css
.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 10px;
  border-radius: 6px;
  font-size: 13px;
  color: #1d1d1f;
  cursor: default;
  transition: background 0.1s;
}

.menu-item:hover {
  background: #007AFF;
  color: white;
}

.menu-item:hover .menu-icon {
  color: white;
}

.menu-item .menu-icon {
  width: 16px;
  height: 16px;
  color: #636366;
}

.menu-item .shortcut {
  margin-left: auto;
  font-size: 12px;
  color: #8E8E93;
}

.menu-item:hover .shortcut {
  color: rgba(255, 255, 255, 0.7);
}
```

### 分隔线

```css
.menu-separator {
  height: 1px;
  background: rgba(0, 0, 0, 0.05);
  margin: 4px 10px;
}

/* Dark Mode */
.menu-separator.dark {
  background: rgba(255, 255, 255, 0.08);
}
```

### 菜单项状态

| 状态 | 背景 | 文字颜色 | 图标颜色 |
|------|------|----------|----------|
| 默认 | 透明 | `#1d1d1f` | `#636366` |
| 悬浮 | `#007AFF` | `white` | `white` |
| 禁用 | 透明 | `#8E8E93` | `#C7C7CC` |

### 动画

```javascript
// 菜单出现动画
const menuAnimation = {
  initial: { opacity: 0, scale: 0.95, y: -4 },
  animate: { opacity: 1, scale: 1, y: 0 },
  exit: { opacity: 0, scale: 0.95, y: -4 },
  transition: {
    type: "spring",
    stiffness: 500,
    damping: 30,
    mass: 0.5
  }
};
```

---

## 10. 历史记录下拉面板

点击"历史记录"按钮时展开的下拉面板，展示所有迭代版本。

### 容器样式

```css
.history-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  min-width: 360px;
  max-height: 400px;
  overflow-y: auto;
  padding: 8px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(24px) saturate(1.8);
  border-radius: 12px;
  border: 0.5px solid rgba(0, 0, 0, 0.1);
  box-shadow:
    0 0 0 0.5px rgba(0, 0, 0, 0.05),
    0 10px 38px -10px rgba(22, 23, 24, 0.35),
    0 10px 20px -15px rgba(22, 23, 24, 0.2);
}

/* Dark Mode */
.history-dropdown.dark {
  background: rgba(40, 40, 40, 0.95);
  border: 0.5px solid rgba(255, 255, 255, 0.1);
}
```

### 版本项样式

```css
.history-item {
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 4px;
  transition: background 0.15s;
}

.history-item:hover {
  background: rgba(0, 0, 0, 0.04);
}

.history-item.current {
  background: rgba(0, 122, 255, 0.08);
  border: 1px solid rgba(0, 122, 255, 0.2);
}

.history-item-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.history-version-badge {
  font-size: 13px;
  font-weight: 600;
  color: #1d1d1f;
}

.history-score {
  font-size: 12px;
  color: #8E8E93;
}

.history-feedback {
  font-size: 11px;
  color: #636366;
  margin-left: 20px;
}

.history-actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
  margin-left: 20px;
}
```

### 版本指示器

| 状态 | 图标 | 颜色 |
|------|------|------|
| 当前版本 | ● (实心圆) | `#007AFF` |
| 历史版本 | ○ (空心圆) | `#8E8E93` |

---

## 11. 版本对比视图

并排对比两个版本的差异。

### 容器样式

```css
.comparison-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--material-thick);
  backdrop-filter: blur(64px) saturate(1.8);
}

.comparison-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.05);
}

.comparison-selector {
  display: flex;
  align-items: center;
  gap: 12px;
}

.comparison-arrow {
  color: #8E8E93;
  font-size: 16px;
}
```

### 对比面板

```css
.comparison-panels {
  display: grid;
  grid-template-columns: 1fr 1fr;
  flex: 1;
  overflow: hidden;
}

.comparison-panel {
  display: flex;
  flex-direction: column;
  border-right: 0.5px solid rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.comparison-panel:last-child {
  border-right: none;
}

.comparison-panel-header {
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.02);
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.05);
  font-size: 13px;
  font-weight: 600;
}

.comparison-panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}
```

### 版本选择下拉

使用与分段控件相同的样式，但以下拉形式呈现。

---

## 12. Loading 骨架屏

页面加载时显示的骨架屏，模拟最终内容布局。

### 骨架元素样式

```css
.skeleton {
  background: linear-gradient(
    90deg,
    rgba(0, 0, 0, 0.06) 25%,
    rgba(0, 0, 0, 0.10) 50%,
    rgba(0, 0, 0, 0.06) 75%
  );
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s ease-in-out infinite;
  border-radius: 6px;
}

/* Dark Mode */
.skeleton.dark {
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.06) 25%,
    rgba(255, 255, 255, 0.10) 50%,
    rgba(255, 255, 255, 0.06) 75%
  );
  background-size: 200% 100%;
}

@keyframes skeleton-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
```

### 骨架布局

```css
.skeleton-score-ring {
  width: 160px;
  height: 160px;
  border-radius: 50%;
}

.skeleton-text-line {
  height: 14px;
  margin-bottom: 8px;
}

.skeleton-text-line.short {
  width: 60%;
}

.skeleton-text-line.medium {
  width: 80%;
}

.skeleton-text-line.full {
  width: 100%;
}

.skeleton-card {
  height: 80px;
  border-radius: 12px;
}

.skeleton-button {
  width: 80px;
  height: 32px;
  border-radius: 6px;
}
```

---

## 13. 成功动画 (Checkmark)

用户点击确定/回滚后显示的成功反馈动画。

### 容器样式

```css
.success-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(8px);
  z-index: 100;
}

/* Dark Mode */
.success-overlay.dark {
  background: rgba(0, 0, 0, 0.6);
}
```

### Checkmark 动画

```css
.success-checkmark {
  width: 80px;
  height: 80px;
}

.checkmark-circle {
  stroke: #34C759;
  stroke-width: 2;
  fill: none;
  stroke-dasharray: 166;
  stroke-dashoffset: 166;
  animation: checkmark-circle 0.6s ease-in-out forwards;
}

.checkmark-check {
  stroke: #34C759;
  stroke-width: 3;
  stroke-linecap: round;
  stroke-linejoin: round;
  fill: none;
  stroke-dasharray: 48;
  stroke-dashoffset: 48;
  animation: checkmark-check 0.3s 0.4s ease-in-out forwards;
}

@keyframes checkmark-circle {
  100% { stroke-dashoffset: 0; }
}

@keyframes checkmark-check {
  100% { stroke-dashoffset: 0; }
}
```

### SVG 结构

```html
<svg class="success-checkmark" viewBox="0 0 52 52">
  <circle class="checkmark-circle" cx="26" cy="26" r="25"/>
  <path class="checkmark-check" d="M14.1 27.2l7.1 7.2 16.7-16.8"/>
</svg>
```

### 动画时序

| 阶段 | 时长 | 说明 |
|------|------|------|
| 圆圈绘制 | 0.6s | 从顶部顺时针绘制 |
| 勾选绘制 | 0.3s | 延迟 0.4s 后开始 |
| 总时长 | 0.7s | 完成后 0.3s 触发关闭 |

---

## 14. 原始 Prompt 弹窗

点击"查看原始输入"时显示的弹窗，展示用户最初输入的原始 prompt。

### 弹窗样式

```css
.original-prompt-modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(24px) saturate(1.8);
  border-radius: 14px;
  border: 0.5px solid rgba(0, 0, 0, 0.1);
  box-shadow:
    0 0 0 0.5px rgba(0, 0, 0, 0.05),
    0 24px 48px -12px rgba(0, 0, 0, 0.25);
  overflow: hidden;
}

.original-prompt-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.05);
}

.original-prompt-title {
  font-size: 15px;
  font-weight: 600;
}

.original-prompt-content {
  padding: 20px;
  max-height: 60vh;
  overflow-y: auto;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  font-family: var(--font-mono);
}
```

### 背景遮罩

```css
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(2px);
  z-index: 50;
}
```

---

# 第四部分：动画系统

## Apple Spring 配置

```javascript
const appleSpring = {
  type: "spring",
  stiffness: 300,
  damping: 30
};
```

## 微交互

| 交互 | 动画 | 时长 |
|------|------|------|
| 按钮按下 | `scale(0.96)` | 0.1s |
| 悬浮上浮 | `translateY(-2px)` | 0.2s |
| Tab 切换 | layoutId 滑动 | spring |
| 评分环 | stroke-dashoffset | 1s spring |
| 卡片选中 | ring 扩散 | 0.15s |

## 减少动画

```css
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

# 第五部分：响应式设计

## 断点

| 断点 | 宽度 | 布局 |
|------|------|------|
| Mobile | < 640px | 单列 |
| Tablet | 640-1024px | 两列 |
| Desktop | > 1024px | 完整 |

## 移动端调整

- 评分区和报告区垂直堆叠
- 方向卡片单列
- 按钮区 sticky 底部
- 内边距 16px

---

# 第六部分：无障碍

## 对比度 (WCAG AA)

| 元素 | 对比度 |
|------|--------|
| 正文 | 14.5:1 ✓ |
| 次要文字 | 5.7:1 ✓ |
| 链接 | 4.5:1 ✓ |

## 焦点环

```css
:focus-visible {
  box-shadow:
    0 0 0 2px white,
    0 0 0 4px #007AFF;
}
```

## 键盘导航

| 按键 | 操作 |
|------|------|
| Tab | 焦点移动 |
| Space/Enter | 激活 |
| Escape | 取消 |

---

# 第七部分：颜色系统

## macOS HIG 颜色

```css
:root {
  /* Accent */
  --color-blue: #007AFF;
  --color-green: #34C759;
  --color-orange: #FF9500;
  --color-red: #FF3B30;

  /* Gray - Light */
  --gray-1: #8E8E93;
  --gray-2: #AEAEB2;
  --gray-3: #C7C7CC;
  --gray-4: #D1D1D6;
  --gray-5: #E5E5EA;
  --gray-6: #F2F2F7;

  /* Gray - Dark */
  --gray-1-dark: #8E8E93;
  --gray-2-dark: #636366;
  --gray-3-dark: #48484A;
  --gray-4-dark: #3A3A3C;
  --gray-5-dark: #2C2C2E;
  --gray-6-dark: #1C1C1E;
}
```

---

# 第八部分：实现检查清单

## 视觉物理
- [ ] 浮动容器有顶部高光 (bezel)
- [ ] 边框 0.5px 发丝线
- [ ] 大面积背景有噪点
- [ ] 毛玻璃含 saturate
- [ ] 分层阴影

## 交互动画
- [ ] 使用 spring 物理
- [ ] 按钮 scale(0.96)
- [ ] Tab layoutId 滑动
- [ ] 支持 reduced-motion
- [ ] 成功动画 (checkmark)

## 组件
- [ ] Traffic Lights
- [ ] 倒计时显示 (含警告状态)
- [ ] 分段控件滑动
- [ ] Bento 卡片悬浮
- [ ] 输入框发光焦点
- [ ] 按钮渐变高光
- [ ] 历史记录下拉面板
- [ ] 版本对比视图
- [ ] Loading 骨架屏
- [ ] 原始 Prompt 弹窗

## 数据交互
- [ ] 调用 GetInputData() 获取数据
- [ ] 调用 GetRemainingSeconds() 更新倒计时
- [ ] Submit() 提交当前版本
- [ ] Rollback() 回滚到历史版本
- [ ] Cancel() 取消操作

## 响应式
- [ ] 375px 正常
- [ ] 768px 正常
- [ ] 1024px+ 正常

## 无障碍
- [ ] WCAG AA 对比度
- [ ] ARIA 标签
- [ ] 键盘可操作
- [ ] 焦点环可见

---

# 第九部分：实现优先级

| P0 | P1 | P2 |
|----|----|----|
| 材质系统 | 深色模式 | 无障碍完善 |
| 窗口头部 + 倒计时 | 响应式 | 动画微调 |
| 评分环形图 | 版本对比视图 | |
| 分段控件 | | |
| Bento 卡片 | | |
| 按钮系统 | | |
| 历史记录下拉 | | |
| Loading 骨架屏 | | |
| 成功动画 | | |
