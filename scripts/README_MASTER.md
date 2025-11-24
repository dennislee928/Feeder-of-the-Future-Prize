# University of Glasgow Inquiry Management System

## 統一管理腳本使用指南

這個系統包含多個腳本來處理格拉斯哥大學研究生招生諮詢的完整流程。使用 `gla-inquiry-manager.sh` 作為統一入口，可以輕鬆管理所有功能。

## 🚀 快速開始

### 1. 獲取 reCAPTCHA Token
```bash
# 手動模式（推薦新手）
./gla-inquiry-manager.sh manual

# 自動化模式（需要 Python 和 Selenium）
./gla-inquiry-manager.sh auto
```

### 2. 開始自動提交
```bash
# 確保已設置 RECAPTCHA_TOKEN
export RECAPTCHA_TOKEN="your_token_here"

# 開始每4小時自動提交
./gla-inquiry-manager.sh submit
```

## 📋 可用命令

### 主要命令

| 命令 | 描述 | 用法 |
|------|------|------|
| `manual` | 手動獲取 reCAPTCHA token | `./gla-inquiry-manager.sh manual` |
| `auto` | 自動化獲取 reCAPTCHA token | `./gla-inquiry-manager.sh auto [--firefox] [--headless]` |
| `submit` | 開始自動提交諮詢 | `./gla-inquiry-manager.sh submit` |
| `status` | 檢查系統狀態 | `./gla-inquiry-manager.sh status` |
| `logs` | 查看日誌 | `./gla-inquiry-manager.sh logs` |
| `help` | 顯示幫助信息 | `./gla-inquiry-manager.sh help` |

### 自動模式選項

- `--firefox`: 使用 Firefox 瀏覽器（默認 Chrome）
- `--headless`: 無頭模式運行（無 GUI）

## 🔄 典型工作流程

### 首次設置
```bash
# 1. 檢查系統狀態
./gla-inquiry-manager.sh status

# 2. 獲取 reCAPTCHA token
./gla-inquiry-manager.sh manual  # 或 auto

# 3. 開始自動提交
./gla-inquiry-manager.sh submit
```

### 日常使用
```bash
# 檢查日誌
./gla-inquiry-manager.sh logs

# 檢查狀態
./gla-inquiry-manager.sh status
```

### Token 過期時
```bash
# 停止當前提交（如果在運行）
# 然後重新獲取 token
./gla-inquiry-manager.sh auto

# 重新開始提交
./gla-inquiry-manager.sh submit
```

## 📁 腳本說明

| 腳本文件 | 功能 | 調用方式 |
|----------|------|----------|
| `gla-inquiry-manager.sh` | **主管理腳本** | 直接運行 |
| `get_recaptcha_token.sh` | 手動 token 獲取助手 | 通過主腳本調用 |
| `run_recaptcha_getter.sh` | 自動化 token 獲取 | 通過主腳本調用 |
| `submit-inquiry.sh` | 自動提交腳本 | 通過主腳本調用 |

## ⚙️ 配置

### 環境變數

```bash
# reCAPTCHA token
export RECAPTCHA_TOKEN="your_token_here"

# 日誌目錄（可選）
export GLA_REQUEST_LOGS="/path/to/logs"
```

### 依賴項安裝

對於自動化模式：
```bash
pip install -r requirements_recaptcha.txt
```

## 📊 監控

### 查看提交狀態
```bash
./gla-inquiry-manager.sh status
```

### 查看日誌
```bash
./gla-inquiry-manager.sh logs
```

### 手動檢查文件
```bash
# 查看 token
cat recaptcha_token.txt

# 查看日誌
ls -la GLA_REQUEST_LOGS/
tail GLA_REQUEST_LOGS/*.log
```

## 🛠️ 故障排除

### Token 相關問題
```bash
# 檢查 token 是否設置
./gla-inquiry-manager.sh status

# 重新獲取 token
./gla-inquiry-manager.sh auto
```

### 腳本無法運行
```bash
# 檢查權限
ls -la scripts/gla-inquiry-manager.sh

# 添加執行權限
chmod +x scripts/gla-inquiry-manager.sh
```

### 自動化模式失敗
```bash
# 檢查 Python 和 Selenium
python --version
pip list | grep selenium

# 安裝依賴項
pip install -r requirements_recaptcha.txt

# 檢查瀏覽器驅動
# ChromeDriver: https://chromedriver.chromium.org/
# GeckoDriver: https://github.com/mozilla/geckodriver/releases
```

## 📝 日誌和數據

### 日誌位置
- 默認: `GLA_REQUEST_LOGS/YYYY-MM-DD.log`
- 自定義: 通過 `GLA_REQUEST_LOGS` 環境變數設置

### Token 存儲
- 文件: `recaptcha_token.txt`
- 環境變數: `RECAPTCHA_TOKEN`

### 數據格式
- 日誌: `[timestamp] [LEVEL] message`
- Token: 長字符串，以 `0.` 開頭

## 🔒 安全注意事項

- reCAPTCHA token 有時效性
- 定期檢查 token 是否仍然有效
- 不要分享您的 token
- 腳本僅用於個人學術申請目的

## 🎯 最佳實踐

1. **定期檢查狀態**: 使用 `status` 命令監控系統
2. **監控日誌**: 使用 `logs` 命令查看提交歷史
3. **及時更新 token**: 當遇到 reCAPTCHA 錯誤時立即更新
4. **備份配置**: 保存重要的環境變數設置

## 📞 獲取幫助

```bash
# 顯示完整幫助
./gla-inquiry-manager.sh help

# 檢查詳細狀態
./gla-inquiry-manager.sh status
```

---

**注意**: 這個系統設計用於格拉斯哥大學研究生招生諮詢。請確保遵守學校的政策和服務條款。
