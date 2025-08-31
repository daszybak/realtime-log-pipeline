tgt_dir:=target
tgt_temp_dir:=$(tgt_dir)/tmp
tgt_artfs_dir:=$(tgt_dir)/artfs
pkg:=github.com/daszybak/realtime-log-pipeline
srcs:=$(shell find backend -name "*.go")

# Build the API.
.PHONY: api
api: $(tgt_artfs_dir)/api

$(tgt_artfs_dir)/api: $(srcs) | $(tgt_artfs_dir)
	( \
		cd backend/cmd/api \
			&& go get -v \
			&& CGO_ENABLED=0 go build \
				-o '../../../$@' \
	)

# Build the worker.
.PHONY: worker
worker: $(tgt_artfs_dir)/worker

$(tgt_artfs_dir)/worker: $(srcs) | $(tgt_artfs_dir)
	( \
		cd backend/cmd/worker \
			&& go get -v \
			&& CGO_ENABLED=0 go build \
				-o '../../../$@' \
	)
	
# Build the aggregator.
.PHONY: aggregator
aggregator: $(tgt_artfs_dir)/aggregator

$(tgt_artfs_dir)/aggregator: $(srcs) | $(tgt_artfs_dir)
	( \
		cd backend/cmd/aggregator \
			&& go get -v \
			&& CGO_ENABLED=0 go build \
				-o '../../../$@' \
	)

# Build the streamer.
.PHONY: streamer
streamer: $(tgt_artfs_dir)/streamer

$(tgt_artfs_dir)/streamer: $(srcs) | $(tgt_artfs_dir)
	( \
		cd backend/cmd/streamer \
			&& go get -v \
			&& CGO_ENABLED=0 go build \
				-o '../../../$@' \
	)

$(tgt_dir):
	mkdir '$@'

$(tgt_artfs_dir): | $(tgt_dir)
	mkdir '$@'

$(tgt_temp_dir): | $(tgt_dir)
	mkdir '$@'

$(tgt_temp_dir)/air_api: | $(tgt_temp_dir)
	mkdir '$@'

$(tgt_temp_dir)/air_api_testdata: | $(tgt_temp_dir)
	mkdir '$@'

$(tgt_temp_dir)/air_worker: | $(tgt_temp_dir)
	mkdir '$@'

$(tgt_temp_dir)/air_worker_testdata: | $(tgt_temp_dir)
	mkdir '$@'

$(tgt_temp_dir)/air_streamer: | $(tgt_temp_dir)
	mkdir '$@'

$(tgt_temp_dir)/air_streamer_testdata: | $(tgt_temp_dir)
	mkdir '$@'

$(tgt_temp_dir)/air_aggregator: | $(tgt_temp_dir)
	mkdir '$@'

$(tgt_temp_dir)/air_aggregator_testdata: | $(tgt_temp_dir)
	mkdir '$@'
