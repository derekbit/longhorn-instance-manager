module github.com/longhorn/longhorn-instance-manager

go 1.20

require (
	github.com/RoaringBitmap/roaring v1.2.3
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.3.0
	github.com/longhorn/backupstore v0.0.0-20230324160313-e1d0b33c2f82
	github.com/longhorn/go-iscsi-helper v0.0.0-20230215054929-acb305e1031b
	github.com/longhorn/go-spdk-helper v0.0.0-20230428184001-b185b2f90a72
	github.com/longhorn/longhorn-engine v1.4.0-rc1.0.20230413074927-aeff7b57a163
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0
	github.com/urfave/cli v1.22.12
	golang.org/x/net v0.9.0
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.7.0
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f
)

require (
	github.com/aws/aws-sdk-go v1.34.2 // indirect
	github.com/bits-and-blooms/bitset v1.2.0 // indirect
	github.com/c9s/goprocinfo v0.0.0-20210130143923-c95fcf8c64a8 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/google/fscrypt v0.3.3 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/honestbee/jobq v1.0.2 // indirect
	github.com/jmespath/go-jmespath v0.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/longhorn/nsfilelock v0.0.0-20200723175406-fa7c83ad0003 // indirect
	github.com/longhorn/sparse-tools v0.0.0-20230408015858-c849def39d3c // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/rancher/go-fibmap v0.0.0-20160418233256-5fc9f8c1ed47 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	k8s.io/apimachinery v0.26.0 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/mount-utils v0.26.0 // indirect
	k8s.io/utils v0.0.0-20221107191617-1a15be271d1d // indirect
)

replace golang.org/x/text v0.3.2 => golang.org/x/text v0.3.3
