FROM python:3.11-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY src/ ./src/
COPY start.sh ./

# Make start script executable
RUN chmod +x start.sh

# Expose port (Render will set PORT env var)
EXPOSE 8084

# Run the application (use PORT env var for Render compatibility)
CMD ["./start.sh"]

