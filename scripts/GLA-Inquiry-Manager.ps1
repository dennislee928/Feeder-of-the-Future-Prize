# University of Glasgow Inquiry Management System - PowerShell Edition
# 格拉斯哥大學諮詢管理系統 - PowerShell 版本

param(
    [Parameter(Mandatory=$false)]
    [string]$Command = "help",

    [Parameter(Mandatory=$false)]
    [switch]$Firefox,

    [Parameter(Mandatory=$false)]
    [switch]$Headless
)

# 顏色定義
$Colors = @{
    "Green" = "Green"
    "Yellow" = "Yellow"
    "Red" = "Red"
    "Blue" = "Cyan"
    "Purple" = "Magenta"
    "Cyan" = "Cyan"
}

# 腳本路徑
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$GetTokenScript = Join-Path $ScriptDir "get_recaptcha_token.sh"
$RunGetterScript = Join-Path $ScriptDir "run_recaptcha_getter.sh"
$SubmitScript = Join-Path $ScriptDir "submit-inquiry.sh"

# 顯示橫幅
function Show-Banner {
    Write-Host "====================================================" -ForegroundColor Green
    Write-Host "  University of Glasgow Inquiry Management System" -ForegroundColor Green
    Write-Host "====================================================" -ForegroundColor Green
    Write-Host ""
}

# 顯示幫助信息
function Show-Help {
    Show-Banner
    Write-Host "Usage: .\GLA-Inquiry-Manager.ps1 -Command <command> [options]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Commands:" -ForegroundColor Cyan
    Write-Host "  manual                  Manual reCAPTCHA token retrieval (with instructions)" -ForegroundColor Green
    Write-Host "  auto                    Automated reCAPTCHA token retrieval (Selenium)" -ForegroundColor Green
    Write-Host "  submit                  Start automated inquiry submission (every 4 hours)" -ForegroundColor Green
    Write-Host "  status                  Check current system status" -ForegroundColor Green
    Write-Host "  logs                    Show recent log entries" -ForegroundColor Green
    Write-Host "  help                    Show this help message" -ForegroundColor Green
    Write-Host ""
    Write-Host "Auto Command Options:" -ForegroundColor Cyan
    Write-Host "  -Firefox                Use Firefox instead of Chrome" -ForegroundColor Yellow
    Write-Host "  -Headless               Run in headless mode" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Cyan
    Write-Host "  .\GLA-Inquiry-Manager.ps1 -Command manual" -ForegroundColor White
    Write-Host "  .\GLA-Inquiry-Manager.ps1 -Command auto" -ForegroundColor White
    Write-Host "  .\GLA-Inquiry-Manager.ps1 -Command auto -Firefox" -ForegroundColor White
    Write-Host "  .\GLA-Inquiry-Manager.ps1 -Command submit" -ForegroundColor White
    Write-Host "  .\GLA-Inquiry-Manager.ps1 -Command status" -ForegroundColor White
    Write-Host "  .\GLA-Inquiry-Manager.ps1 -Command logs" -ForegroundColor White
    Write-Host ""
}

# 檢查腳本是否存在
function Test-Scripts {
    $missingScripts = @()

    if (-not (Test-Path $GetTokenScript)) {
        $missingScripts += "get_recaptcha_token.sh"
    }

    if (-not (Test-Path $RunGetterScript)) {
        $missingScripts += "run_recaptcha_getter.sh"
    }

    if (-not (Test-Path $SubmitScript)) {
        $missingScripts += "submit-inquiry.sh"
    }

    if ($missingScripts.Count -gt 0) {
        Write-Host "✗ Missing required scripts:" -ForegroundColor Red
        foreach ($script in $missingScripts) {
            Write-Host "  - $script" -ForegroundColor Red
        }
        Write-Host ""
        Write-Host "Please ensure all scripts are in the same directory." -ForegroundColor Yellow
        exit 1
    }
}

# 檢查 bash 是否可用
function Test-Bash {
    try {
        $bashVersion = & bash --version 2>$null | Select-Object -First 1
        return $true
    }
    catch {
        Write-Host "✗ Bash is not available. Please install Git Bash or WSL." -ForegroundColor Red
        Write-Host "Download from: https://git-scm.com/downloads" -ForegroundColor Yellow
        return $false
    }
}

# 運行 bash 腳本
function Invoke-BashScript {
    param(
        [string]$ScriptPath,
        [string[]]$Arguments = @()
    )

    if (-not (Test-Bash)) {
        return
    }

    $argString = $Arguments -join " "
    $command = "bash `"$ScriptPath`" $argString"

    try {
        Invoke-Expression $command
    }
    catch {
        Write-Host "✗ Failed to execute script: $ScriptPath" -ForegroundColor Red
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    }
}

# 手動模式
function Invoke-ManualMode {
    Write-Host "Starting manual reCAPTCHA token retrieval..." -ForegroundColor Blue
    Write-Host ""

    if (Test-Path $GetTokenScript) {
        Invoke-BashScript -ScriptPath $GetTokenScript
    }
    else {
        Write-Host "✗ Manual token getter script not found: $GetTokenScript" -ForegroundColor Red
    }
}

# 自動模式
function Invoke-AutoMode {
    Write-Host "Starting automated reCAPTCHA token retrieval..." -ForegroundColor Blue
    Write-Host ""

    $arguments = @()
    if ($Firefox) {
        $arguments += "--firefox"
    }
    if ($Headless) {
        $arguments += "--headless"
    }

    if (Test-Path $RunGetterScript) {
        Invoke-BashScript -ScriptPath $RunGetterScript -Arguments $arguments
    }
    else {
        Write-Host "✗ Automated token getter script not found: $RunGetterScript" -ForegroundColor Red
    }
}

# 提交模式
function Invoke-SubmitMode {
    Write-Host "Starting automated inquiry submission..." -ForegroundColor Blue
    Write-Host ""

    # 檢查 reCAPTCHA token
    $recaptchaToken = $env:RECAPTCHA_TOKEN
    $tokenFile = Join-Path $ScriptDir "recaptcha_token.txt"

    if ([string]::IsNullOrEmpty($recaptchaToken)) {
        Write-Host "⚠ Warning: RECAPTCHA_TOKEN environment variable is not set" -ForegroundColor Yellow

        if (Test-Path $tokenFile) {
            Write-Host "Found token file: recaptcha_token.txt" -ForegroundColor Yellow
            Write-Host "Loading token from file..." -ForegroundColor Yellow
            $recaptchaToken = Get-Content $tokenFile -Raw
            $env:RECAPTCHA_TOKEN = $recaptchaToken
            Write-Host "✓ Token loaded from file" -ForegroundColor Green
        }
        else {
            Write-Host "✗ No RECAPTCHA_TOKEN set and no token file found" -ForegroundColor Red
            Write-Host ""
            Write-Host "Please get a reCAPTCHA token first:" -ForegroundColor Yellow
            Write-Host "  .\GLA-Inquiry-Manager.ps1 -Command manual    # Manual instructions" -ForegroundColor White
            Write-Host "  .\GLA-Inquiry-Manager.ps1 -Command auto      # Automated retrieval" -ForegroundColor White
            Write-Host ""
            exit 1
        }
    }
    else {
        Write-Host "✓ RECAPTCHA_TOKEN is set" -ForegroundColor Green
    }

    Write-Host ""

    if (Test-Path $SubmitScript) {
        # 設置環境變數並運行腳本
        $env:RECAPTCHA_TOKEN = $recaptchaToken
        Invoke-BashScript -ScriptPath $SubmitScript
    }
    else {
        Write-Host "✗ Submit script not found: $SubmitScript" -ForegroundColor Red
    }
}

# 檢查狀態
function Get-SystemStatus {
    Write-Host "System Status Check" -ForegroundColor Blue
    Write-Host "==================" -ForegroundColor Blue
    Write-Host ""

    # 檢查腳本
    Write-Host "Scripts:" -ForegroundColor Cyan
    if (Test-Path $GetTokenScript) {
        Write-Host "✓ get_recaptcha_token.sh (manual mode)" -ForegroundColor Green
    }
    else {
        Write-Host "✗ get_recaptcha_token.sh (manual mode)" -ForegroundColor Red
    }

    if (Test-Path $RunGetterScript) {
        Write-Host "✓ run_recaptcha_getter.sh (auto mode)" -ForegroundColor Green
    }
    else {
        Write-Host "✗ run_recaptcha_getter.sh (auto mode)" -ForegroundColor Red
    }

    if (Test-Path $SubmitScript) {
        Write-Host "✓ submit-inquiry.sh (submit mode)" -ForegroundColor Green
    }
    else {
        Write-Host "✗ submit-inquiry.sh (submit mode)" -ForegroundColor Red
    }

    Write-Host ""

    # 檢查環境
    Write-Host "Environment:" -ForegroundColor Cyan
    if ($env:RECAPTCHA_TOKEN) {
        Write-Host "✓ RECAPTCHA_TOKEN is set ($($env:RECAPTCHA_TOKEN.Length) characters)" -ForegroundColor Green
    }
    else {
        Write-Host "✗ RECAPTCHA_TOKEN is not set" -ForegroundColor Red
    }

    if ($env:GLA_REQUEST_LOGS) {
        Write-Host "✓ GLA_REQUEST_LOGS is set to: $($env:GLA_REQUEST_LOGS)" -ForegroundColor Green
    }
    else {
        Write-Host "⚠ GLA_REQUEST_LOGS not set, using default: GLA_REQUEST_LOGS" -ForegroundColor Yellow
    }

    Write-Host ""

    # 檢查文件
    Write-Host "Files:" -ForegroundColor Cyan
    $tokenFile = Join-Path $ScriptDir "recaptcha_token.txt"
    if (Test-Path $tokenFile) {
        $tokenLength = (Get-Item $tokenFile).Length
        Write-Host "✓ recaptcha_token.txt exists ($tokenLength characters)" -ForegroundColor Green
    }
    else {
        Write-Host "✗ recaptcha_token.txt not found" -ForegroundColor Red
    }

    $logDir = $env:GLA_REQUEST_LOGS
    if ([string]::IsNullOrEmpty($logDir)) {
        $logDir = Join-Path $ScriptDir "GLA_REQUEST_LOGS"
    }

    if (Test-Path $logDir) {
        $logCount = (Get-ChildItem $logDir -Filter "*.log" -File).Count
        Write-Host "✓ GLA_REQUEST_LOGS directory exists ($logCount log files)" -ForegroundColor Green
    }
    else {
        Write-Host "✗ GLA_REQUEST_LOGS directory not found" -ForegroundColor Red
    }

    Write-Host ""

    # 檢查依賴項
    Write-Host "Dependencies:" -ForegroundColor Cyan
    if (Get-Command python -ErrorAction SilentlyContinue) {
        $pythonVersion = & python --version 2>$null
        Write-Host "✓ Python available: $pythonVersion" -ForegroundColor Green
    }
    elseif (Get-Command python3 -ErrorAction SilentlyContinue) {
        $pythonVersion = & python3 --version 2>$null
        Write-Host "✓ Python available: $pythonVersion" -ForegroundColor Green
    }
    else {
        Write-Host "✗ Python not found" -ForegroundColor Red
    }

    if (Test-Bash) {
        Write-Host "✓ Bash available" -ForegroundColor Green
    }
    else {
        Write-Host "✗ Bash not available" -ForegroundColor Red
    }

    Write-Host ""

    # 總結
    Write-Host "Summary:" -ForegroundColor Cyan
    $hasToken = ($env:RECAPTCHA_TOKEN -or (Test-Path $tokenFile))
    if ($hasToken) {
        Write-Host "✓ System is ready for submission" -ForegroundColor Green
    }
    else {
        Write-Host "⚠ System needs reCAPTCHA token" -ForegroundColor Yellow
        Write-Host "  Run: .\GLA-Inquiry-Manager.ps1 -Command manual  or  -Command auto" -ForegroundColor Yellow
    }
}

# 顯示日誌
function Show-Logs {
    Write-Host "Recent Log Entries" -ForegroundColor Blue
    Write-Host "==================" -ForegroundColor Blue
    Write-Host ""

    $logDir = $env:GLA_REQUEST_LOGS
    if ([string]::IsNullOrEmpty($logDir)) {
        $logDir = Join-Path $ScriptDir "GLA_REQUEST_LOGS"
    }

    if (-not (Test-Path $logDir)) {
        Write-Host "✗ Log directory not found: $logDir" -ForegroundColor Red
        Write-Host "No logs available yet. Run the submit command first." -ForegroundColor Yellow
        return
    }

    # 找到最新的日誌文件
    $latestLog = Get-ChildItem $logDir -Filter "*.log" -File |
                 Sort-Object LastWriteTime -Descending |
                 Select-Object -First 1

    if (-not $latestLog) {
        Write-Host "No log files found in $logDir" -ForegroundColor Yellow
        return
    }

    Write-Host "Latest log file: $($latestLog.FullName)" -ForegroundColor Cyan
    Write-Host ""

    # 顯示最後10行
    Write-Host "Last 10 entries:" -ForegroundColor Cyan
    Write-Host ""
    Get-Content $latestLog.FullName -Tail 10
}

# 主函數
function Main {
    Test-Scripts

    switch ($Command.ToLower()) {
        "help" {
            Show-Help
        }
        "manual" {
            Invoke-ManualMode
        }
        "auto" {
            Invoke-AutoMode
        }
        "submit" {
            Invoke-SubmitMode
        }
        "status" {
            Get-SystemStatus
        }
        "logs" {
            Show-Logs
        }
        default {
            Write-Host "✗ Unknown command: $Command" -ForegroundColor Red
            Write-Host ""
            Show-Help
            exit 1
        }
    }
}

# 執行主函數
Main
