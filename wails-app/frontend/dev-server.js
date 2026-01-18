/**
 * 开发服务器 - 用于在浏览器中预览和调试 UI
 *
 * 使用方式:
 *   node dev-server.js [--port 8080] [--input test-input.json]
 *
 * 功能:
 *   - 启动静态文件服务器
 *   - 注入 Mock Wails API
 *   - 支持自定义测试数据
 */

import http from 'http';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// 解析命令行参数
const args = process.argv.slice(2);
let port = 8080;
let inputFile = path.join(__dirname, '..', 'test-input.json');

for (let i = 0; i < args.length; i++) {
  if (args[i] === '--port' && args[i + 1]) {
    port = parseInt(args[i + 1], 10);
    i++;
  } else if (args[i] === '--input' && args[i + 1]) {
    inputFile = path.resolve(args[i + 1]);
    i++;
  }
}

// MIME 类型映射
const mimeTypes = {
  '.html': 'text/html',
  '.js': 'application/javascript',
  '.css': 'text/css',
  '.json': 'application/json',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.svg': 'image/svg+xml',
  '.ico': 'image/x-icon',
};

// Mock Wails API 脚本
function getMockScript(inputData) {
  return `
<script>
// ========== Mock Wails API ==========
(function() {
  const mockData = ${JSON.stringify(inputData, null, 2)};
  let remainingSeconds = 600; // 10 分钟

  // 模拟倒计时
  setInterval(() => {
    if (remainingSeconds > 0) remainingSeconds--;
  }, 1000);

  // Mock window.go.main.App
  window.go = {
    main: {
      App: {
        GetInputData: async function() {
          console.log('[Mock] GetInputData called');
          return mockData;
        },

        GetRemainingSeconds: async function() {
          return remainingSeconds;
        },

        Submit: async function(directions, userInput) {
          console.log('[Mock] Submit called:', { directions, userInput });
          alert('提交成功！\\n\\n选择的方向: ' + directions.join(', ') + '\\n用户输入: ' + userInput);
        },

        Cancel: async function() {
          console.log('[Mock] Cancel called');
          alert('已取消');
        },

        Rollback: async function(iterationId, directions, userInput) {
          console.log('[Mock] Rollback called:', { iterationId, directions, userInput });
          alert('回滚到版本: ' + iterationId + '\\n\\n选择的方向: ' + directions.join(', ') + '\\n用户输入: ' + userInput);
        }
      }
    }
  };

  console.log('[Mock] Wails API 已注入');
  console.log('[Mock] 测试数据:', mockData);
})();
</script>
`;
}

// 创建服务器
const server = http.createServer((req, res) => {
  let filePath = path.join(__dirname, req.url === '/' ? 'index.html' : req.url);

  // 安全检查：防止目录遍历
  if (!filePath.startsWith(__dirname)) {
    res.writeHead(403);
    res.end('Forbidden');
    return;
  }

  const ext = path.extname(filePath);
  const contentType = mimeTypes[ext] || 'application/octet-stream';

  fs.readFile(filePath, (err, content) => {
    if (err) {
      if (err.code === 'ENOENT') {
        res.writeHead(404);
        res.end('File not found: ' + req.url);
      } else {
        res.writeHead(500);
        res.end('Server error: ' + err.code);
      }
      return;
    }

    // 对 HTML 文件注入 Mock 脚本
    if (ext === '.html') {
      try {
        const inputData = JSON.parse(fs.readFileSync(inputFile, 'utf8'));
        const mockScript = getMockScript(inputData);
        // 在 </head> 前注入
        content = content.toString().replace('</head>', mockScript + '</head>');
      } catch (e) {
        console.error('加载测试数据失败:', e.message);
      }
    }

    res.writeHead(200, { 'Content-Type': contentType });
    res.end(content);
  });
});

server.listen(port, () => {
  console.log('');
  console.log('╔════════════════════════════════════════════════════════════╗');
  console.log('║       Prompt Optimizer UI 开发服务器                       ║');
  console.log('╠════════════════════════════════════════════════════════════╣');
  console.log('║                                                            ║');
  console.log(`║  🌐 浏览器访问: http://localhost:${port}/                      ║`);
  console.log('║                                                            ║');
  console.log(`║  📁 测试数据: ${path.basename(inputFile).padEnd(40)}║`);
  console.log('║                                                            ║');
  console.log('║  💡 提示:                                                  ║');
  console.log('║     - Wails API 已自动 Mock                                ║');
  console.log('║     - 修改文件后刷新浏览器即可                             ║');
  console.log('║     - 按 Ctrl+C 停止服务器                                 ║');
  console.log('║                                                            ║');
  console.log('╚════════════════════════════════════════════════════════════╝');
  console.log('');
});
