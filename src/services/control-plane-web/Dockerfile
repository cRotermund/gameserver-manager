FROM node:22-alpine AS base

WORKDIR /app
COPY package.json package-lock.json tsconfig.json ./

FROM base AS builder
RUN npm ci
COPY src/ src/
RUN npm run build
RUN npm ci --omit=dev && npm cache clean --force

FROM node:22-alpine AS runner

WORKDIR /app
COPY --from=builder /app/dist/ dist/
COPY --from=builder /app/node_modules/ node_modules/
COPY --from=builder /app/package.json ./

EXPOSE 3000
USER node
CMD ["node", "dist/index.js"]
