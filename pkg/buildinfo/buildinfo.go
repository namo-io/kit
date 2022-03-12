package buildinfo

var (
	Tag string = "dev"
)

func GetVersion() string {
	return Tag
}
