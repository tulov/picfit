package constants

const (
	Version = "0.9.1"
)

var (
	// Branch is the compiled branch
	Branch string

	// Revision is the compiled revision
	Revision string

	// LatestCommitMessage is the latest commit message
	LatestCommitMessage string

	// BuildTime is the compiled build time
	BuildTime string

	// Compiler is the compiler used during build
	Compiler string
)

const (
	BottomLeft  = "bottom-left"
	BottomRight = "bottom-right"
	TopLeft     = "top-left"
	TopRight    = "top-right"
)

var StickPositions = []string{
	BottomLeft,
	BottomRight,
	TopLeft,
	TopRight,
}
