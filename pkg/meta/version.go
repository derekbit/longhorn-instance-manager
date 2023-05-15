package meta

const (
	// InstanceManagerAPIVersion is used for compatibility check for longhorn-manager
	InstanceManagerAPIVersion    = 4
	InstanceManagerAPIMinVersion = 1

	// InstanceManagerProxyAPIVersion is used for compatibility check for longhorn-manager
	InstanceManagerProxyAPIVersion    = 4
	InstanceManagerProxyAPIMinVersion = 1

	// InstanceManagerDiskServiceAPIVersion used to communicate with the user e.g. longhorn-manager
	InstanceManagerDiskServiceAPIVersion    = 1
	InstanceManagerDiskServiceAPIMinVersion = 1

	// InstanceManagerInstanceServiceAPIVersion used to communicate with the user e.g. longhorn-manager
	InstanceManagerInstanceServiceAPIVersion    = 1
	InstanceManagerInstanceServiceAPIMinVersion = 1
)

// Following variables are filled in by main.go
var (
	Version   string
	GitCommit string
	BuildDate string
)

type VersionOutput struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`

	InstanceManagerAPIVersion    int `json:"instanceManagerAPIVersion"`
	InstanceManagerAPIMinVersion int `json:"instanceManagerAPIMinVersion"`

	InstanceManagerProxyAPIVersion    int `json:"instanceManagerProxyAPIVersion"`
	InstanceManagerProxyAPIMinVersion int `json:"instanceManagerProxyAPIMinVersion"`

	InstanceManagerDiskServiceAPIVersion    int `json:"instanceManagerDiskServiceAPIVersion"`
	InstanceManagerDiskServiceAPIMinVersion int `json:"instanceManagerDiskServiceAPIMinVersion"`

	InstanceManagerInstanceServiceAPIVersion    int `json:"instanceManagerInstanceServiceAPIVersion"`
	InstanceManagerInstanceServiceAPIMinVersion int `json:"instanceManagerInstanceServiceAPIMinVersion"`
}

func GetVersion() VersionOutput {
	return VersionOutput{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,

		InstanceManagerAPIVersion:    InstanceManagerAPIVersion,
		InstanceManagerAPIMinVersion: InstanceManagerAPIMinVersion,

		InstanceManagerProxyAPIVersion:    InstanceManagerProxyAPIVersion,
		InstanceManagerProxyAPIMinVersion: InstanceManagerProxyAPIMinVersion,

		InstanceManagerDiskServiceAPIVersion:    InstanceManagerDiskServiceAPIVersion,
		InstanceManagerDiskServiceAPIMinVersion: InstanceManagerDiskServiceAPIMinVersion,

		InstanceManagerInstanceServiceAPIVersion:    InstanceManagerInstanceServiceAPIVersion,
		InstanceManagerInstanceServiceAPIMinVersion: InstanceManagerInstanceServiceAPIMinVersion,
	}
}
