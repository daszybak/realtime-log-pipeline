project := 'realtime-log-pipeline'
img := 'trustap' / project
tgt_dir := 'target'
tgt_artfs_dir := tgt_dir / 'artfs'
tgt_tmp_dir := tgt_dir / 'tmp'
tgt_api := tgt_artfs_dir / 'api'
tgt_worker := tgt_artfs_dir / 'worker'
tgt_aggregator := tgt_artfs_dir / 'aggregator'
tgt_streamer := tgt_artfs_dir / 'streamer'

src_dirs := 'backend/cmd backend/pkg'

dev_all: dev_api dev_worker dev_aggregator dev_streamer dev_app
run_all: run_api run_worker run_aggregator run_streamer

# Run `app` in dev mode. In production, the app is statically served.
dev_app port='8080':
    ( \
        cd app \
            && npx vite \
                --port='{{port}}' \
                --strictPort \
                --host \
    )

# Run `api`.
run_api conf='backend/configs/api.yaml' addr='0.0.0.0:8081':
    make '{{tgt_api}}'
    '{{tgt_api}}' \
        '{{conf}}' \
        '{{addr}}'

# Run `api` in hot-reload mode.
dev_api conf='backend/configs/api.yaml' addr='0.0.0.0:8081':
    make \
        '{{tgt_tmp_dir}}/air_api' \
        '{{tgt_tmp_dir}}/air_api_testdata'
    air \
        -c 'backend/configs/build/api.air.toml' \
        '{{conf}}' \
        '{{addr}}'

# Run `worker`.
run_worker conf='backend/configs/worker.yaml' addr='0.0.0.0:8082':
    make '{{tgt_worker}}'
    '{{tgt_worker}}' \
        '{{conf}}' \
        '{{addr}}'

# Run `worker` in hot-reload mode.
dev_worker conf='backend/configs/worker.yaml' addr='0.0.0.0:8082':
    make \
        '{{tgt_tmp_dir}}/air_worker' \
        '{{tgt_tmp_dir}}/air_worker_testdata'
    air \
        -c 'backend/configs/build/worker.air.toml' \
        '{{conf}}' \
        '{{addr}}'

# Run `aggregator`.
run_aggregator conf='backend/configs/aggregator.yaml' addr='0.0.0.0:8083':
    make '{{tgt_aggregator}}'
    '{{tgt_aggregator}}' \
        '{{conf}}' \
        '{{addr}}'

# Run `aggregator` in hot-reload mode.
dev_aggregator conf='backend/configs/aggregator.yaml' addr='0.0.0.0:8083':
    make \
        '{{tgt_tmp_dir}}/air_aggregator' \
        '{{tgt_tmp_dir}}/air_aggregator_testdata'
    air \
        -c 'backend/configs/build/aggregator.air.toml' \
        '{{conf}}' \
        '{{addr}}'

# Run `streamer`.
run_streamer conf='backend/configs/streamer.yaml' addr='0.0.0.0:8084':
    make '{{tgt_streamer}}'
    '{{tgt_streamer}}' \
        '{{conf}}' \
        '{{addr}}'

# Run `streamer` in hot-reload mode.
dev_streamer conf='backend/configs/streamer.yaml' addr='0.0.0.0:8084':
    make \
        '{{tgt_tmp_dir}}/air_streamer' \
        '{{tgt_tmp_dir}}/air_streamer_testdata'
    air \
        -c 'backend/configs/build/streamer.air.toml' \
        '{{conf}}' \
        '{{addr}}'

# Run style checks.
check_style: check_go_style check_yaml_style check_js_lint

# Check for semantic issues in Go files.
check_go_lint:
    make api
    revive \
        -config='backend/configs/build/revive.toml' \
        -formatter=plain \
        backend/cmd/... \
        backend/pkg/...
    ( \
        cd api \
            && golangci-lint run \
                --config='configs/build/golangci.yaml' \
                cmd/... \
                internal/... \
                pkg/... \
    )

# Check for semantic issues in JS files.
check_js_lint:
    ( \
        cd app \
            && npm run lint \
    )

# Run style checks for Go files.
check_go_style:
    ! (gofumpt -d {{src_dirs}} | grep '')
    ! (goimports -d {{src_dirs}} | grep '')
    ! (gofmt -s -d {{src_dirs}} | grep '')

# Run style checks for YAML files.
check_yaml_style:
    prettier \
        --ignore-path='.gitignore' \
        --ignore-path='.prettierignore' \
        --list-different \
        '**/*.yaml'

# Format source files.
fmt: fmt_go fmt_yaml fmt_app

# Format Go files.
fmt_go:
    gofumpt -w backend

fmt_app:
    ( \
        cd app \
            && npm run lint:fix \
    )

# Format YAML files.
fmt_yaml:
    prettier \
        --ignore-path='.gitignore' \
        --ignore-path='.prettierignore' \
        --write \
        '**/*.yaml'
