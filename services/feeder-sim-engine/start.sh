#!/bin/sh
# 啟動腳本 - 使用 PORT 環境變數（Render 兼容）

PORT=${PORT:-8081}
exec uvicorn src.app:app --host 0.0.0.0 --port $PORT

