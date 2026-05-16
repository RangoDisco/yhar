# TODO: slim image
FROM node:24-slim AS web-base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV CI=true
COPY . /app
WORKDIR /app

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


FROM web-base as web-prod
COPY --from=web-prod-build /app/node_modules /app/node_modules
COPY --from=web-prod-build /app/build /app/build
EXPOSE 3000
ENV NODE_ENV=production
ENV ORIGIN=http://127.0.0.1:3000
CMD [ "node", "build/index.js" ]