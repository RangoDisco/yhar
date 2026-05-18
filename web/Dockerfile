# TODO: slim image
FROM node:24-alpine AS web-base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV CI=true
WORKDIR /app
COPY . .

FROM web-base AS web-dev
RUN corepack enable
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install
EXPOSE 5173
ENV NODE_ENV=development
CMD ["pnpm", "run", "dev", "--host"]

FROM web-base AS web-prod-build
RUN corepack enable
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build
RUN pnpm prune --prod

FROM node:24-alpine AS web-prod
LABEL org.opencontainers.image.source="https://github.com/rangodisco/yhar"
WORKDIR /app
COPY --from=web-prod-build /app/node_modules ./node_modules
COPY --from=web-prod-build /app/build ./build
EXPOSE 3000
ENV NODE_ENV=production
ENV ORIGIN=http://127.0.0.1:3000
CMD [ "node", "build/index.js" ]