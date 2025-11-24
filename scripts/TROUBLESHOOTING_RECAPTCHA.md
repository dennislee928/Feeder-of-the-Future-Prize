# reCAPTCHA 問題排查指南

## 問題：頁面只顯示 "recaptcha fail"，沒有實際表單

如果自動化腳本遇到此問題，這通常意味著伺服器記住了之前的失敗狀態。

## 解決方案

### 方案 1: 等待並重試（推薦）

伺服器可能暫時封鎖了您的 IP 或會話。建議：

1. **等待 5-10 分鐘**
2. **清除瀏覽器緩存和 cookies**
3. **重新運行腳本**

```bash
# 等待後重試
./gla-inquiry-manager.sh auto --edge
```

### 方案 2: 手動獲取 reCAPTCHA Token（最可靠）

由於自動化可能觸發反機器人檢測，建議手動獲取 token：

#### 步驟：

1. **打開瀏覽器開發者工具**
   - 按 `F12` 或右鍵 → 檢查

2. **切換到 Network（網絡）標籤**
   - 確保正在記錄網絡活動

3. **訪問表單頁面**
   ```
   https://www.gla.ac.uk/study/enquire/send/index.html
   ```

4. **填寫表單**
   - Name: Test User
   - Email: test@example.com
   - 等等...

5. **完成 reCAPTCHA 挑戰**
   - 勾選 "I'm not a robot"
   - 完成圖片驗證（如果需要）

6. **點擊提交按鈕**

7. **在 Network 標籤中找到 POST 請求**
   - 尋找 `send/index.html` 請求
   - 點擊該請求

8. **查看 Payload（有效負載）**
   - 切換到 "Payload" 或 "Request" 標籤
   - 找到 `g-recaptcha-response` 參數
   - 複製整個 token 值

9. **設置環境變數**
   ```bash
   export RECAPTCHA_TOKEN="你複製的token值"
   ```

10. **運行提交腳本**
    ```bash
    ./gla-inquiry-manager.sh submit
    ```

### 方案 3: 使用不同的瀏覽器

Edge 可能有特定的問題，嘗試使用其他瀏覽器：

```bash
# 使用 Chrome（如果已安裝）
./gla-inquiry-manager.sh auto

# 使用 Firefox（如果已安裝）
./gla-inquiry-manager.sh auto --firefox
```

### 方案 4: 使用 InPrivate/Incognito 模式

手動在瀏覽器的隱私模式下訪問：

1. 打開 Edge 的 InPrivate 視窗（Ctrl+Shift+N）
2. 訪問表單頁面
3. 按照方案 2 的步驟手動獲取 token

### 方案 5: 更改網絡環境

如果可能，嘗試：

1. **切換到不同的網絡**（例如：手機熱點）
2. **使用 VPN**（如果允許）
3. **等待更長時間**（可能需要幾小時）

## 為什麼會發生這種情況？

1. **頻繁請求**: 短時間內多次提交表單
2. **自動化檢測**: Google reCAPTCHA 檢測到自動化行為
3. **IP 限制**: 伺服器暫時封鎖了您的 IP
4. **會話狀態**: 瀏覽器保留了失敗的會話狀態

## 預防措施

1. **不要頻繁運行腳本**
   - 建議間隔至少 4 小時

2. **使用新的瀏覽器配置文件**
   - 腳本現在會自動使用臨時配置文件

3. **定期更新 token**
   - reCAPTCHA token 會過期
   - 建議每次使用前獲取新 token

4. **監控日誌**
   - 檢查 `GLA_REQUEST_LOGS` 目錄中的日誌
   - 查看是否有重複的失敗

## 快速診斷

運行以下命令檢查當前狀態：

```bash
# 檢查最近的日誌
tail -n 50 GLA_REQUEST_LOGS/$(date +%Y-%m-%d).log

# 測試網絡連接
curl -I https://www.gla.ac.uk/study/enquire/send/index.html

# 檢查是否有 Edge 進程殘留
tasklist | grep -i edge
```

## 需要立即提交？

如果需要立即提交表單，最可靠的方法是：

1. **完全手動提交**
   - 直接在瀏覽器中訪問並提交表單
   - 不使用任何自動化

2. **使用手動 token 獲取 + 自動提交**
   - 按照方案 2 手動獲取 token
   - 使用 `./gla-inquiry-manager.sh submit` 提交

## 聯繫支持

如果問題持續存在：

1. 檢查 University of Glasgow 的系統狀態
2. 確認表單頁面是否正常運作
3. 考慮直接聯繫大學的 IT 支持

## 成功指標

當問題解決時，您應該看到：

- ✓ 瀏覽器打開並顯示完整的表單
- ✓ 可以看到所有輸入欄位
- ✓ reCAPTCHA 挑戰框出現
- ✓ 沒有 "recaptcha fail" 錯誤訊息

