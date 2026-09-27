# Build stage
FROM node:23 AS builder

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install

COPY . .
RUN npm run build

# Runtime stage
FROM nginx:1.25

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
