#!/bin/bash

# University of Glasgow Postgraduate Admissions Inquiry Submission Script
# 格拉斯哥大學研究生招生諮詢提交腳本
# 每4小時自動執行並包含時間戳

# 顏色輸出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 間隔時間（4小時 = 14400秒）
INTERVAL=14400

# 日誌配置
GLA_REQUEST_LOGS="${GLA_REQUEST_LOGS:-GLA_REQUEST_LOGS}"
LOG_FILE="${GLA_REQUEST_LOGS}/$(date '+%Y-%m-%d').log"

# reCAPTCHA 配置
RECAPTCHA_TOKEN="${RECAPTCHA_TOKEN:-0.wwyqAPnNy1Vh8_QUztz3BzRrKjVUJbHVXxR2CYRMWgu9Lqx_msoZM0_Ttr0xMCTWxlN7UwHUrktrfP0iACG3GemWFSAkMQ7x3s1pY4D-wTMO61KIqa2wn_XkBkAexkfCScT2Am7Jzbw7J-KoV16ji2ktWlT3LpYQDeSIIUL7aU2xZZFpbDHyQUVhqa37a_54jeKvurcQLM_2ObDwBSlUhoidf7nwPpJsggESybp71kG79VWUVICLy2F_oUte_wr6Njnn5sHeikOirPoZKb6uRvJR0KOC-Dgqxr6nNcUe2I5tF-YYOQPqsXJAYkQbhyD1Jr2ioIZQGNIGCffdNNoJSwvWbtzBQZje-g37UPcpAyvOT7as7HRCxETw8nfIW8999ccbBFBc6yXxegWfm2blSyJIaakjnK_6fFpDUd30GvI2D1tJ8p5v2TibZ5uxaPeWTMqiGmIUsAKpR1yaqZnsWfKFWhIKpfkYtYXi2yWu77gb9pdppMfUcltcMamA6Okx0YleQkoVtRgHooQ-_IQ5M59lC3CUJRJ0m4wT7z8SVdzHrjXdToYytQvmf_eJJ48Lj26t0eBFt-aCHVn4hSdpfBtjmUca4qctEoTW66dHcYnaI5xtlvHuNgIJ5bHLnvJCNKE3g_h-MaZWs9qC-JQit15avNePuCWt1ZX7JKdT0JbL-jPoIvptYMzOHABcPoE5WBsFNUUXrkysTLSYnmWw5Y8y0ji8wUgGkGOby6lsbo55QEaQYKkPz8sXyQsLcy9zmmcABQZlRhu8X3hC5kQ0QFyzbZAtNvkHCmJddd8HjwlJ-iwp8nWQtP4Gp6KYnkfllAOVNUOI07aVt4RzNzlnAIDVDDcoO_Z90G_BRnBEXa1fI3GbgrcAvz2zgaJa0eMIOo0qYuG8Js-klygcFXvy5-cZTtNQlkaibiA3AvxI5yA.aN70bIIxbPx93moaqfMvnA.2f03d53f1616d1128f9423ce23f13dd8101d40694d6ff83ad2df2d2fa3c0be2d}"

# 創建日誌目錄
mkdir -p "$GLA_REQUEST_LOGS"

# 日誌函數
log_message() {
    local level="$1"
    local message="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S UTC')
    local log_entry="[$timestamp] [$level] $message"

    # 寫入日誌文件
    echo "$log_entry" >> "$LOG_FILE"

    # 同時輸出到控制台
    case "$level" in
        "INFO")  echo -e "${BLUE}$log_entry${NC}" ;;
        "SUCCESS") echo -e "${GREEN}$log_entry${NC}" ;;
        "ERROR") echo -e "${RED}$log_entry${NC}" ;;
        "WARN")  echo -e "${YELLOW}$log_entry${NC}" ;;
        *)       echo "$log_entry" ;;
    esac
}

# 捕獲 SIGINT (Ctrl+C) 信號以優雅退出
trap 'log_message "INFO" "Script stopped by user"; exit 0' INT

log_message "INFO" "Starting automated University of Glasgow postgraduate admissions inquiry submission (every 4 hours)..."
log_message "INFO" "Log file: $LOG_FILE"
log_message "INFO" "reCAPTCHA Token: ${RECAPTCHA_TOKEN:0:50}..."
log_message "INFO" "Press Ctrl+C to stop the script"

# 表單數據
URL="https://www.gla.ac.uk/study/enquire/send/index.html"

# 循環執行
while true; do
    # 生成當前時間戳
    TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S UTC')

    log_message "INFO" "======================================== Starting new submission cycle"
    log_message "INFO" "Execution time: ${TIMESTAMP}"

    # 準備 POST 數據（動態包含時間戳）
POST_DATA=(
    --data-urlencode "AA_RepresentativePhone="
    --data-urlencode "career=Information Studies"
    --data-urlencode "category=PGT"
    --data-urlencode "country=Application Self Service Issues"
    --data-urlencode "dateofbirth=Taiwan"
    --data-urlencode "email=17-DEC-96"
    --data-urlencode "g-recaptcha-response=${RECAPTCHA_TOKEN}"
    --data-urlencode "hasapplied=0.wwyqAPnNy1Vh8_QUztz3BzRrKjVUJbHVXxR2CYRMWgu9Lqx_msoZM0_Ttr0xMCTWxlN7UwHUrktrfP0iACG3GemWFSAkMQ7x3s1pY4D-wTMO61KIqa2wn_XkBkAexkfCScT2Am7Jzbw7J-KoV16ji2ktWlT3LpYQDeSIIUL7aU2xZZFpbDHyQUVhqa37a_54jeKvurcQLM_2ObDwBSlUhoidf7nwPpJsggESybp71kG79VWUVICLy2F_oUte_wr6Njnn5sHeikOirPoZKb6uRvJR0KOC-Dgqxr6nNcUe2I5tF-YYOQPqsXJAYkQbhyD1Jr2ioIZQGNIGCffdNNoJSwvWbtzBQZje-g37UPcpAyvOT7as7HRCxETw8nfIW8999ccbBFBc6yXxegWfm2blSyJIaakjnK_6fFpDUd30GvI2D1tJ8p5v2TibZ5uxaPeWTMqiGmIUsAKpR1yaqZnsWfKFWhIKpfkYtYXi2yWu77gb9pdppMfUcltcMamA6Okx0YleQkoVtRgHooQ-_IQ5M59lC3CUJRJ0m4wT7z8SVdzHrjXdToYytQvmf_eJJ48Lj26t0eBFt-aCHVn4hSdpfBtjmUca4qctEoTW66dHcYnaI5xtlvHuNgIJ5bHLnvJCNKE3g_h-MaZWs9qC-JQit15avNePuCWt1ZX7JKdT0JbL-jPoIvptYMzOHABcPoE5WBsFNUUXrkysTLSYnmWw5Y8y0ji8wUgGkGOby6lsbo55QEaQYKkPz8sXyQsLcy9zmmcABQZlRhu8X3hC5kQ0QFyzbZAtNvkHCmJddd8HjwlJ-iwp8nWQtP4Gp6KYnkfllAOVNUOI07aVt4RzNzlnAIDVDDcoO_Z90G_BRnBEXa1fI3GbgrcAvz2zgaJa0eMIOo0qYuG8Js-klygcFXvy5-cZTtNQlkaibiA3AvxI5yA.aN70bIIxbPx93moaqfMvnA.2f03d53f1616d1128f9423ce23f13dd8101d40694d6ff83ad2df2d2fa3c0be2d"
    --data-urlencode "id=no"
    --data-urlencode "name="
    --data-urlencode "optin=Pei Chen, Lee"
    --data-urlencode "service=false"
    --data-urlencode "subcategory=Service Desk"
    --data-urlencode "summary=Uploading documents"
    --data-urlencode "symptom=PGT - Application Self Service Issues - Uploading documents [request-time: ${TIMESTAMP}]"
    --data-urlencode "year=


Dear Postgraduate Admissions Team,

I am an applicant for MSc Software Development.

Full name: Pei-Chen Lee (李沛宸)

Application Number - 01658104

Programme: MSc Software Development

I chose the option for the University to contact my referee directly.

My referee, Yihsiu,Chen (Department of Digital Content and technologies, National Cheng Chi University), received an email with the link to the online reference form and the following user ID: UOG_OAS_RECOM_143070.

However, when he follow the link and log in, the system shows a \"Search Criteria – User ID\" page.
Following that he clicked  \"Search\", the system returns the message:
\"No matching values were found.\""
)

    log_message "INFO" "Submitting to: ${URL}"
    log_message "INFO" "Sending POST request..."

    # 發送 POST 請求並捕獲完整響應
    RESPONSE_FILE=$(mktemp)
    HTTP_CODE=$(curl -s -w "%{http_code}" -o "$RESPONSE_FILE" -X POST "${POST_DATA[@]}" "$URL")
    RESPONSE_BODY=$(cat "$RESPONSE_FILE")
    rm -f "$RESPONSE_FILE"

    # 記錄服務器響應
    log_message "INFO" "SERVER RESPONSE - HTTP Status Code: ${HTTP_CODE}"
    if [ -n "$RESPONSE_BODY" ]; then
        log_message "INFO" "SERVER RESPONSE - Body: ${RESPONSE_BODY}"

        # 檢查是否包含 reCAPTCHA 失敗訊息
        if echo "$RESPONSE_BODY" | grep -q "recaptcha fail"; then
            log_message "ERROR" "✗ reCAPTCHA validation failed!"
            log_message "WARN" "SOLUTION: Get a new reCAPTCHA token by:"
            log_message "WARN" "  1. Visit the University of Glasgow inquiry page in browser"
            log_message "WARN" "  2. Complete the reCAPTCHA challenge"
            log_message "WARN" "  3. Intercept the POST request to get new token"
            log_message "WARN" "  4. Set RECAPTCHA_TOKEN environment variable: export RECAPTCHA_TOKEN='new_token_here'"
            log_message "WARN" "  5. Or update the RECAPTCHA_TOKEN variable in the script"
            return 1
        fi
    else
        log_message "INFO" "SERVER RESPONSE - Body: (empty)"
    fi

    if [ "$HTTP_CODE" = "200" ]; then
        log_message "SUCCESS" "✓ Inquiry submission successful! HTTP Status Code: ${HTTP_CODE}"
    else
        log_message "ERROR" "✗ Inquiry submission failed. HTTP Status Code: ${HTTP_CODE}"
        log_message "WARN" "Please check network connection and form data"
        return 1
    fi

    log_message "INFO" "Submission cycle completed. Next execution in 4 hours"

    # 等待4小時後繼續下一次執行
    sleep $INTERVAL
done
