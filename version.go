package stringx

const (
	major = 0
	minor = 5
	patch = 5
)

var Version = Format("v%d.%d.%d", major, minor, patch)
