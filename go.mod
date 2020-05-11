module github.com/fuskovic/networker

go 1.13

require (
	github.com/google/gopacket v1.1.17
	github.com/jackpal/gateway v1.0.6
	github.com/spf13/cobra v1.0.0
)

replace (
	github.com/fuskovic/networker/cmd => ./cmd
	github.com/fuskovic/networker/pkg/capture => ./pkg/capture
	github.com/fuskovic/networker/pkg/list => ./pkg/list
	github.com/fuskovic/networker/pkg/lookup => ./pkg/lookup
	github.com/fuskovic/networker/pkg/scan => ./pkg/scan

)
