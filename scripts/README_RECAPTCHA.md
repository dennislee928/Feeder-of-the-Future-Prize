# University of Glasgow Inquiry Submission - reCAPTCHA Solution

## Problem: reCAPTCHA Validation Failed

The script encountered a reCAPTCHA validation failure with the response: `{"message":"recaptcha fail"}`

## Root Cause

The reCAPTCHA token in the script has expired or is invalid. reCAPTCHA tokens are time-limited and must be freshly generated for each submission.

## Solutions

### Option 1: Update Environment Variable (Recommended)

1. **Get a fresh reCAPTCHA token:**
   ```bash
   # Visit the University of Glasgow inquiry page in your browser
   # Complete the reCAPTCHA challenge manually
   # Use browser developer tools to intercept the POST request
   # Copy the new g-recaptcha-response value
   ```

2. **Set the environment variable:**
   ```bash
   export RECAPTCHA_TOKEN="your_new_fresh_recaptcha_token_here"
   ```

3. **Run the script:**
   ```bash
   ./scripts/submit-inquiry.sh
   ```

### Option 2: Update Script Directly

1. Open `scripts/submit-inquiry.sh`
2. Find the `RECAPTCHA_TOKEN` variable (around line 22)
3. Replace the old token with your new fresh token
4. Save and run the script

### Option 3: Browser-based Solution

If the automated script continues to fail, consider submitting the inquiry manually through the browser:

1. Visit: https://www.gla.ac.uk/study/enquire/send/index.html
2. Fill out the form with your details
3. Complete the reCAPTCHA challenge
4. Submit the form

## How to Get a Fresh reCAPTCHA Token

### Using Browser Developer Tools:

1. **Open Chrome/Edge Developer Tools** (F12)
2. **Go to Network tab**
3. **Visit:** https://www.gla.ac.uk/study/enquire/send/index.html
4. **Fill out the form** (use the same data as in the script)
5. **Complete reCAPTCHA** challenge
6. **Click submit**
7. **In Network tab**, find the POST request to `/study/enquire/send/index.html`
8. **Check Request Headers/Payload** to find `g-recaptcha-response` parameter
9. **Copy the token** value

### Alternative: Browser Extensions

Use browser extensions like:
- "reCAPTCHA Solver" (automates reCAPTCHA solving)
- "EditThisCookie" (to inspect cookies if needed)

## Token Format

A valid reCAPTCHA token looks like:
```
0.ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.
ABCDEF... (very long string, usually 1000+ characters)
```

## Prevention

To prevent this issue in the future:
- Monitor the log files regularly for reCAPTCHA failures
- Set up automated alerts when submissions fail
- Consider implementing reCAPTCHA token refresh logic in the script

## Troubleshooting

If you continue to have issues:
1. Check your internet connection
2. Verify the form data is correct
3. Ensure you're not being rate-limited
4. Try submitting during off-peak hours
5. Consider using a VPN if IP blocking is suspected
