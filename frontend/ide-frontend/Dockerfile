# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Build-time args for Vite（Render 會從 envVars 傳入；若未傳入則用預設值）
ARG VITE_API_BASE_URL=https://feeder-ide-api-d2na.onrender.com/api/v1
ARG VITE_SIM_API_BASE_URL=https://feeder-sim-engine-fnvf.onrender.com
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
ENV VITE_SIM_API_BASE_URL=$VITE_SIM_API_BASE_URL

# Copy package files
COPY package.json ./
COPY package-lock.json* ./

# Install dependencies
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi

# Copy source code
COPY . .

# Build the application (Vite 會將 env 變數內嵌至 bundle)
RUN npm run build

# Runtime stage
FROM nginx:alpine

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

