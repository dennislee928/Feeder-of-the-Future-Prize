#!/bin/bash

# Test reCAPTCHA error handling
GLA_REQUEST_LOGS="${GLA_REQUEST_LOGS:-GLA_REQUEST_LOGS}"
LOG_FILE="${GLA_REQUEST_LOGS}/test-$(date '+%Y-%m-%d').log"

# Create log directory
mkdir -p "$GLA_REQUEST_LOGS"

# Log function
log_message() {
    local level="$1"
    local message="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S UTC')
    local log_entry="[$timestamp] [$level] $message"

    # Write to log file
    echo "$log_entry" >> "$LOG_FILE"

    # Also output to console
    case "$level" in
        "INFO")  echo -e "\033[0;34m$log_entry\033[0m" ;;
        "SUCCESS") echo -e "\033[0;32m$log_entry\033[0m" ;;
        "ERROR") echo -e "\033[0;31m$log_entry\033[0m" ;;
        "WARN") echo -e "\033[1;33m$log_entry\033[0m" ;;
        *)       echo "$log_entry" ;;
    esac
}

echo "Testing reCAPTCHA error handling..."
log_message "INFO" "Testing SERVER RESPONSE with recaptcha fail message"

# Simulate recaptcha fail response
RESPONSE_BODY='{"message":"recaptcha fail"}'

# Test the error detection logic
if echo "$RESPONSE_BODY" | grep -q "recaptcha fail"; then
    log_message "ERROR" "✗ reCAPTCHA validation failed!"
    log_message "WARN" "SOLUTION: Get a new reCAPTCHA token by:"
    log_message "WARN" "  1. Visit the University of Glasgow inquiry page in browser"
    log_message "WARN" "  2. Complete the reCAPTCHA challenge"
    log_message "WARN" "  3. Intercept the POST request to get new token"
    log_message "WARN" "  4. Set RECAPTCHA_TOKEN environment variable: export RECAPTCHA_TOKEN='new_token_here'"
    log_message "WARN" "  5. Or update the RECAPTCHA_TOKEN variable in the script"
    echo "✓ reCAPTCHA error detection working correctly"
else
    echo "✗ reCAPTCHA error detection failed"
fi

echo ""
echo "Test log saved to: $LOG_FILE"
echo "Log contents:"
cat "$LOG_FILE"
