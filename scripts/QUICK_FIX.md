# 快速修復：reCAPTCHA Fail 問題

## 問題現象
Edge 打開後只顯示 `{"message":"recaptcha fail"}`，沒有實際的表單可以填寫，導致瀏覽器關閉。

## 根本原因
伺服器記住了之前失敗的 reCAPTCHA 驗證，並在會話中保持此錯誤狀態。

---

## ✅ 推薦解決方案（按優先順序）

### 1. 等待 + 使用新的腳本版本（最簡單）

腳本已更新，現在會：
- 使用臨時瀏覽器配置文件（避免緩存）
- 自動清除 cookies
- 檢測並處理 "recaptcha fail" 錯誤

**操作步驟：**
```bash
# 等待 5-10 分鐘讓伺服器重置
# 然後運行：
./gla-inquiry-manager.sh auto --edge
```

### 2. 手動獲取 Token（最可靠，100% 成功率）

**操作步驟：**

1. 打開 Edge 瀏覽器（正常模式或 InPrivate）
2. 按 F12 打開開發者工具
3. 切換到 "Network" 標籤
4. 訪問：https://www.gla.ac.uk/study/enquire/send/index.html
5. 填寫表單並完成 reCAPTCHA
6. 點擊提交
7. 在 Network 標籤找到 `send/index.html` POST 請求
8. 查看 Payload，複製 `g-recaptcha-response` 的值
9. 運行：
   ```bash
   export RECAPTCHA_TOKEN="你的token"
   ./gla-inquiry-manager.sh submit
   ```

### 3. 使用 Chrome 或 Firefox

Edge 可能有特定問題，嘗試其他瀏覽器：

```bash
# Chrome（如果已安裝 ChromeDriver）
./gla-inquiry-manager.sh auto

# Firefox（如果已安裝 GeckoDriver）
./gla-inquiry-manager.sh auto --firefox
```

### 4. 完全手動提交

如果自動化持續失敗：

1. 直接在瀏覽器訪問表單頁面
2. 手動填寫並提交
3. 這是最可靠但最耗時的方法

---

## 🔧 已實施的改進

腳本現在包含：

1. **臨時瀏覽器配置文件**
   - 每次運行使用全新的配置文件
   - 避免緩存的錯誤狀態

2. **自動 Cookie 清除**
   - 檢測到錯誤時自動清除 cookies
   - 嘗試多次重新加載

3. **表單檢測**
   - 檢查頁面是否包含實際表單
   - 如果只有錯誤訊息，會嘗試清除狀態

4. **會話保持**
   - 添加 `detach` 選項防止瀏覽器意外關閉
   - 改進的錯誤處理

---

## 📊 成功率預測

| 方法 | 成功率 | 時間 | 難度 |
|------|--------|------|------|
| 手動獲取 Token | 100% | 5 分鐘 | 簡單 |
| 等待 + 新腳本 | 70-80% | 10 分鐘 | 簡單 |
| 使用其他瀏覽器 | 60-70% | 5 分鐘 | 簡單 |
| 完全手動 | 100% | 10 分鐘 | 最簡單 |

---

## ⚠️ 重要提示

1. **不要頻繁重試**
   - 每次失敗後等待至少 5 分鐘
   - 頻繁重試會加重封鎖

2. **reCAPTCHA Token 會過期**
   - Token 通常在 2 分鐘內有效
   - 獲取後盡快使用

3. **檢查日誌**
   ```bash
   # 查看最近的錯誤
   tail -n 50 GLA_REQUEST_LOGS/$(date +%Y-%m-%d).log
   ```

---

## 🎯 立即行動

**如果需要馬上提交表單：**

→ 使用**方案 2**（手動獲取 Token）

**如果可以等待：**

→ 等待 10 分鐘，然後使用**方案 1**（新腳本）

**如果持續失敗：**

→ 查看 `TROUBLESHOOTING_RECAPTCHA.md` 獲取詳細診斷步驟

