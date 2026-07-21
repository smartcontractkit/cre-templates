# oven/bun:1.3.5 (x64 / amd64 digest)
FROM oven/bun@sha256:f48ca3e042de3c9da0afb25455414db376c873bf88f79ffb6f26aec82bdb6da2 AS builder

RUN apt-get update && apt-get install -y make && rm -rf /var/lib/apt/lists/*

WORKDIR /app

ENV CRE_DOCKER_BUILD_IMAGE=true

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .

RUN make build

# --- EXPORT STAGE ---
FROM scratch AS exporter
COPY --from=builder /app/wasm /wasm
