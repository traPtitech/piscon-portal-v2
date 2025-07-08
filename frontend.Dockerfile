# https://bun.sh/guides/ecosystem/docker

# use the official Bun image
# see all versions at https://hub.docker.com/r/oven/bun/tags
FROM oven/bun:1 AS base
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

FROM nginx:1-alpine AS production

COPY --from=prerelease /usr/src/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
