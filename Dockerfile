FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json tsconfig.json ./
COPY src/ src/
RUN npm ci && npm run build

FROM node:22-alpine

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci --omit=dev --ignore-scripts

COPY --from=builder /app/dist/ dist/
COPY web/ web/

EXPOSE 8080
ENV PORT=8080

CMD ["node", "dist/index.js"]
