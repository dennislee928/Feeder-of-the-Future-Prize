# University of Glasgow Inquiry Management System

## 完整的自動化解決方案

這個系統整合了多個腳本，提供從獲取 reCAPTCHA token 到自動提交諮詢的完整解決方案。

## 📦 系統架構

```
GLA Inquiry System/
├── 📄 GLA_INQUIRY_SYSTEM_README.md          # 本文檔
├── 📁 scripts/
│   ├── 🔧 gla-inquiry-manager.sh           # Linux/Mac 主管理腳本
│   ├── 🔧 GLA-Inquiry-Manager.ps1          # Windows PowerShell 版本
│   ├── 🔧 get_recaptcha_token.sh           # 手動 token 獲取助手
│   ├── 🔧 run_recaptcha_getter.sh          # 自動化 token 獲取包裝器
│   ├── 🔧 submit-inquiry.sh                # 自動提交腳本
│   ├── 🤖 get_recaptcha_automated.py       # Selenium 自動化腳本
│   ├── 📋 requirements_recaptcha.txt       # Python 依賴項
│   ├── 📋 test_log.sh                      # 日誌測試腳本
│   ├── 📋 test_recaptcha.sh                # reCAPTCHA 測試腳本
│   ├── 📖 README_MASTER.md                 # 主腳本詳細說明
│   ├── 📖 README_RECAPTCHA.md              # reCAPTCHA 問題解決
│   └── 📖 README_RECAPTCHA_AUTOMATION.md   # 自動化說明
├── 📁 GLA_REQUEST_LOGS/                    # 日誌目錄
│   └── 📄 2025-11-19.log                   # 每日日誌
└── 📄 recaptcha_token.txt                  # Token 存儲文件
```

## 🚀 快速開始（推薦）

### Windows 用戶
```powershell
# 進入腳本目錄
cd scripts

# 顯示幫助
.\GLA-Inquiry-Manager.ps1 -Command help

# 檢查系統狀態
.\GLA-Inquiry-Manager.ps1 -Command status

# 獲取 reCAPTCHA token（選擇一種方式）
.\GLA-Inquiry-Manager.ps1 -Command manual    # 手動說明
.\GLA-Inquiry-Manager.ps1 -Command auto      # 自動化（需要 Python）

# 開始自動提交
.\GLA-Inquiry-Manager.ps1 -Command submit

# 查看日誌
.\GLA-Inquiry-Manager.ps1 -Command logs
```

### Linux/Mac 用戶
```bash
# 進入腳本目錄
cd scripts

# 添加執行權限
chmod +x *.sh

# 顯示幫助
./gla-inquiry-manager.sh help

# 檢查系統狀態
./gla-inquiry-manager.sh status

# 獲取 reCAPTCHA token（選擇一種方式）
./gla-inquiry-manager.sh manual    # 手動說明
./gla-inquiry-manager.sh auto      # 自動化（需要 Python）

# 開始自動提交
./gla-inquiry-manager.sh submit

# 查看日誌
./gla-inquiry-manager.sh logs
```

## 🔄 完整工作流程

### 1. 初始設置
```bash
# 檢查所有依賴項
./gla-inquiry-manager.sh status

# 安裝 Python 依賴項（如果使用自動化）
pip install -r requirements_recaptcha.txt
```

### 2. 獲取 reCAPTCHA Token
**選項 A: 手動方法（推薦新手）**
```bash
./gla-inquiry-manager.sh manual
# 按照屏幕指示操作
```

**選項 B: 自動化方法（需要 Python + Selenium）**
```bash
# 使用 Chrome
./gla-inquiry-manager.sh auto

# 使用 Firefox
./gla-inquiry-manager.sh auto --firefox

# 無頭模式
./gla-inquiry-manager.sh auto --headless
```

### 3. 開始自動提交
```bash
./gla-inquiry-manager.sh submit
# 腳本會每4小時自動提交一次
```

### 4. 監控和維護
```bash
# 查看最新日誌
./gla-inquiry-manager.sh logs

# 檢查系統狀態
./gla-inquiry-manager.sh status

# 當 reCAPTCHA 過期時，重新獲取
./gla-inquiry-manager.sh auto
```

## 📋 腳本功能對照表

| 腳本名稱 | 功能 | 適用場景 | 依賴項 |
|----------|------|----------|--------|
| `gla-inquiry-manager.sh` | 主管理腳本 | Linux/Mac | Bash |
| `GLA-Inquiry-Manager.ps1` | 主管理腳本 | Windows | PowerShell |
| `get_recaptcha_token.sh` | 手動 token 助手 | 所有系統 | Curl |
| `run_recaptcha_getter.sh` | 自動化包裝器 | 所有系統 | Python, Selenium |
| `submit-inquiry.sh` | 自動提交 | 所有系統 | Curl |
| `get_recaptcha_automated.py` | Selenium 核心 | 所有系統 | Python, Selenium |

## ⚙️ 配置和環境變數

### 必須設置
```bash
export RECAPTCHA_TOKEN="your_token_here"
```

### 可選設置
```bash
export GLA_REQUEST_LOGS="/custom/log/directory"
```

## 🔧 故障排除

### 常見問題

#### 1. "reCAPTCHA validation failed"
**解決方案**:
```bash
./gla-inquiry-manager.sh auto
```

#### 2. "Python not found"
**解決方案**:
```bash
# 安裝 Python
# Windows: https://python.org
# Linux: sudo apt install python3
# Mac: brew install python3

pip install -r requirements_recaptcha.txt
```

#### 3. "WebDriver not found"
**解決方案**:
```bash
# ChromeDriver: https://chromedriver.chromium.org/
# GeckoDriver: https://github.com/mozilla/geckodriver/releases
# 將驅動程序放在 PATH 中
```

#### 4. 腳本無執行權限
**解決方案**:
```bash
chmod +x scripts/*.sh
```

### 檢查所有狀態
```bash
./gla-inquiry-manager.sh status
```

## 📊 日誌和監控

### 日誌位置
- **默認**: `GLA_REQUEST_LOGS/YYYY-MM-DD.log`
- **自定義**: 通過 `GLA_REQUEST_LOGS` 變數設置

### 日誌內容示例
```
[2025-11-19 11:36:31 UTC] [INFO] Starting automated University of Glasgow postgraduate admissions inquiry submission (every 4 hours)...
[2025-11-19 11:36:31 UTC] [INFO] SERVER RESPONSE - HTTP Status Code: 200
[2025-11-19 11:36:31 UTC] [SUCCESS] ✓ Inquiry submission successful! HTTP Status Code: 200
```

### 監控腳本運行狀態
```bash
# 查看最新提交
tail -f GLA_REQUEST_LOGS/$(date +%Y-%m-%d).log

# 檢查 token 狀態
./gla-inquiry-manager.sh status
```

## 🔄 自動化程度說明

| 方法 | 手動操作 | 自動化程度 | 適用用戶 |
|------|----------|------------|----------|
| 純手動 | 高 | 0% | 所有用戶 |
| 手動 token + 自動提交 | 中 | 75% | 大多數用戶 |
| 全自動化 | 低 | 90% | 技術用戶 |

## 📝 使用技巧

### 1. 定期檢查
```bash
# 每天檢查一次狀態
./gla-inquiry-manager.sh status
```

### 2. 及時更新 Token
當看到 "recaptcha fail" 錯誤時，立即更新 token。

### 3. 備份重要文件
- `recaptcha_token.txt`
- `GLA_REQUEST_LOGS/` 目錄

### 4. 環境變數持久化
將 `RECAPTCHA_TOKEN` 添加到 shell 配置中。

## 🎯 最佳實踐

1. **使用主管理腳本**: 不要直接調用單個腳本
2. **定期檢查狀態**: 使用 `status` 命令
3. **監控日誌**: 使用 `logs` 命令
4. **及時更新 token**: 當出現錯誤時立即處理
5. **備份配置**: 保存重要的環境變數

## 🔒 安全和隱私

- Token 僅用於個人學術申請
- 不要分享您的 reCAPTCHA token
- 腳本不會上傳任何敏感信息
- 所有通信都通過 HTTPS

## 📞 技術支持

如果遇到問題：

1. 運行 `status` 檢查系統狀態
2. 查看 `logs` 了解詳細錯誤
3. 參考對應的 README 文件
4. 檢查依賴項是否正確安裝

---

**系統設計理念**: 提供從完全手動到高度自動化的完整解決方案，讓不同技術水平的用戶都能成功使用。
