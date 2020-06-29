module github.com/fuskovic/networker

go 1.13

require (
	github.com/google/gopacket v1.1.17
	github.com/jackpal/gateway v1.0.6
	github.com/spf13/pflag v1.0.3
	github.com/tatsushid/go-fastping v0.0.0-20160109021039-d7bb493dee3e
	go.coder.com/cli v0.4.0
	go.coder.com/flog v0.0.0-20190906214207-47dd47ea0512
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092 // indirect
)

replace (
	github.com/fuskovic/networker/cmd => ./cmd
	github.com/fuskovic/networker/pkg/capture => ./pkg/capture
	github.com/fuskovic/networker/pkg/list => ./pkg/list
	github.com/fuskovic/networker/pkg/lookup => ./pkg/lookup
	github.com/fuskovic/networker/pkg/scan => ./pkg/scan

)
