FROM cgr.dev/chainguard/wolfi-base AS base

RUN apk add --no-cache go make git npm pnpm

FROM base AS bin
WORKDIR /app
COPY . .
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/uv \
    --mount=type=cache,target=/root/go/pkg/mod \
    make all

FROM base AS tools
RUN apk add --no-cache curl python-3.13 py3.13-pip
WORKDIR /app
COPY . .
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/uv \
    --mount=type=cache,target=/root/go/pkg/mod \
    UV_LINK_MODE=copy BIN_DIR=/bin make package-tools

FROM cgr.dev/chainguard/postgres:latest-dev AS build-pgvector
RUN apk add build-base git postgresql-dev
RUN git clone --branch v0.8.0 https://github.com/pgvector/pgvector.git && \
    cd pgvector && \
    make clean && \
    make OPTFLAGS="" && \
    make install && \
    cd .. && \
    rm -rf pgvector

FROM cgr.dev/chainguard/postgres:latest-dev AS final
ENV POSTGRES_USER=orion101
ENV POSTGRES_PASSWORD=orion101
ENV POSTGRES_DB=orion101
ENV PGDATA=/data/postgresql

COPY --from=build-pgvector /usr/lib/postgresql17/vector.so /usr/lib/postgresql17/
COPY --from=build-pgvector /usr/share/postgresql17/extension/vector* /usr/share/postgresql17/extension/

RUN apk add --no-cache git python-3.13 py3.13-pip openssh-server npm bash tini procps libreoffice
COPY --chmod=0755 /tools/package-chrome.sh /
RUN /package-chrome.sh && rm /package-chrome.sh
RUN sed -E 's/^#(PermitRootLogin)no/\1yes/' /etc/ssh/sshd_config -i
RUN ssh-keygen -A
RUN mkdir /run/sshd && /usr/sbin/sshd
COPY encryption.yaml /
COPY --chmod=0755 run.sh /bin/run.sh

COPY --link --from=tools /app/orion101-tools /orion101-tools
COPY --from=bin /app/bin/orion101 /bin/

EXPOSE 22
# libreoffice executables
ENV PATH=/orion101-tools/venv/bin:$PATH:/usr/lib/libreoffice/program
ENV HOME=/data
ENV XDG_CACHE_HOME=/data/cache
ENV GPTSCRIPT_SYSTEM_TOOLS_DIR=/orion101-tools/
ENV ORION101_SERVER_WORKSPACE_TOOL=/orion101-tools/workspace-provider
ENV ORION101_SERVER_DATASETS_TOOL=/orion101-tools/datasets
ENV ORION101_SERVER_TOOL_REGISTRY=/orion101-tools
ENV ORION101_SERVER_ENCRYPTION_CONFIG_FILE=/encryption.yaml
ENV GOMEMLIMIT=1GiB
ENV BAAAH_THREADINESS=20
ENV TERM=vt100
WORKDIR /data
VOLUME /data
ENTRYPOINT ["run.sh"]