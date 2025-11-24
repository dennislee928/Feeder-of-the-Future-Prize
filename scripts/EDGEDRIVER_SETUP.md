# EdgeDriver 設置指南

## 問題：網絡連接失敗，無法自動下載 EdgeDriver

如果遇到 "Could not reach host" 或 "Are you offline?" 錯誤，可以手動下載和設置 EdgeDriver。

## 方法 1: 手動下載 EdgeDriver（推薦）

### 步驟：

1. **下載 EdgeDriver**
   - 訪問：https://developer.microsoft.com/en-us/microsoft-edge/tools/webdriver/
   - 選擇與您的 Edge 版本匹配的 EdgeDriver
   - 下載 ZIP 文件

2. **解壓縮**
   - 解壓下載的 ZIP 文件
   - 找到 `msedgedriver.exe` 文件

3. **放置驅動程序**
   
   **選項 A: 放在腳本目錄（最簡單）**
   ```bash
   # 將 msedgedriver.exe 複製到 scripts 目錄
   cp msedgedriver.exe scripts/
   ```

   **選項 B: 添加到系統 PATH**
   ```bash
   # 將 msedgedriver.exe 放在一個固定位置，例如：
   C:\tools\msedgedriver.exe
   
   # 然後添加到 PATH 環境變數
   ```

   **選項 C: 使用環境變數**
   ```bash
   # 設置 EDGEDRIVER_PATH 環境變數
   export EDGEDRIVER_PATH="C:\path\to\msedgedriver.exe"
   ```

4. **驗證**
   ```bash
   # 運行腳本，應該能自動找到驅動程序
   ./gla-inquiry-manager.sh auto --edge
   ```

## 方法 2: 使用 Chrome 或 Firefox（替代方案）

如果 EdgeDriver 設置有困難，可以使用其他瀏覽器：

```bash
# 使用 Chrome（默認）
./gla-inquiry-manager.sh auto

# 使用 Firefox
./gla-inquiry-manager.sh auto --firefox
```

## 方法 3: 檢查網絡連接

如果網絡正常，webdriver-manager 應該能自動下載：

```bash
# 檢查網絡連接
ping google.com

# 如果有代理，可能需要設置代理環境變數
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=http://proxy.example.com:8080
```

## 查找 Edge 版本

要下載正確版本的 EdgeDriver，需要知道您的 Edge 版本：

1. 打開 Edge 瀏覽器
2. 點擊右上角的三個點（...）
3. 選擇 "設定" > "關於 Microsoft Edge"
4. 查看版本號（例如：120.0.2210.91）

## 常見問題

### Q: 如何知道 EdgeDriver 是否正確設置？
A: 運行腳本時，如果看到 "✓ EdgeDriver initialized from..." 訊息，表示設置成功。

### Q: 可以同時使用多個瀏覽器嗎？
A: 可以，但每次運行只能使用一個瀏覽器。分別運行：
- `./gla-inquiry-manager.sh auto` (Chrome)
- `./gla-inquiry-manager.sh auto --edge` (Edge)
- `./gla-inquiry-manager.sh auto --firefox` (Firefox)

### Q: EdgeDriver 需要更新嗎？
A: 當 Edge 瀏覽器更新時，建議更新 EdgeDriver 以匹配版本。

## 快速設置腳本

如果您已經下載了 EdgeDriver，可以快速設置：

```bash
# 假設 EdgeDriver 在當前目錄
cp msedgedriver.exe scripts/

# 或者設置環境變數
export EDGEDRIVER_PATH="$(pwd)/msedgedriver.exe"
```

## 需要幫助？

如果仍然遇到問題：
1. 檢查錯誤訊息中的詳細信息
2. 確認 Edge 瀏覽器已安裝
3. 確認 msedgedriver.exe 有執行權限
4. 嘗試使用 Chrome 或 Firefox 作為替代
