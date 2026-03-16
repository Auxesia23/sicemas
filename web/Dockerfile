FROM docker.io/oven/bun:alpine AS builder

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .

RUN bun --bun run build

FROM oven/bun:alpine AS runner

WORKDIR /app

COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json .

ENV NODE_ENV=production

EXPOSE 3000

CMD ["bun", "./build/index.js"]
