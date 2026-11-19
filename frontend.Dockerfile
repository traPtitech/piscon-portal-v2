# https://bun.sh/guides/ecosystem/docker

# use the official Bun image
# see all versions at https://hub.docker.com/r/oven/bun/tags
FROM oven/bun:1@sha256:1d22548f177f4ff38e0c15baa3c2c46c6c2060ba2d1bc573bb4720b407548f1d AS base
WORKDIR /usr/src/app

# install dependencies into temp directory
# this will cache them and speed up future builds
FROM base AS install
RUN mkdir -p /temp/dev
COPY ./client/package.json ./client/bun.lockb /temp/dev/
RUN cd /temp/dev && bun install --frozen-lockfile

FROM base AS prerelease

ENV NODE_ENV=production
RUN --mount=type=bind,source=./client/,target=/usr/src/app,readwrite \
  --mount=type=bind,from=install,source=/temp/dev/node_modules,target=/usr/src/app/node_modules,readwrite \
  mkdir -p /usr/src/dist && \
  bun run build-only --outDir /usr/src/dist --emptyOutDir

FROM nginx:1-alpine@sha256:0563c36f94a6d3ec3584161b77859f740f7e20b4853f52c796a864ce6908e8ea AS production

COPY --from=prerelease /usr/src/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
