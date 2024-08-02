FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

RUN apt update && apt install -y build-essential git curl

# RUN curl -O https://dl.google.com/go/go1.22.0.linux-amd64.tar.gz
# RUN tar -xvzf go1.22.0.linux-amd64.tar.gz -C /usr/local
# ENV GOROOT=/usr/local/go
# ENV GOPATH=$HOME/go
# ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# RUN git clone https://github.com/avezina-ubik/migrate.git
# RUN cd migrate/cmd/migrate && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ../../cli/build/migrate -ldflags='-extldflags "-static"' -tags 'postgres mysql redshift cassandra spanner cockroachdb yugabytedb clickhouse mongodb sqlserver firebird neo4j pgx pgx5 rqlite libsql file go_bindata github github_ee bitbucket aws_s3 google_cloud_storage godoc_vfs gitlab' .
# RUN cp migrate/cli/build/migrate /usr/local/bin/migrate

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>
# USER root

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
