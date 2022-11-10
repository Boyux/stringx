package stringx

const (
	major = 0
	minor = 3
	patch = 0
)

var Version = Format("v%d.%d.%d", major, minor, patch)
