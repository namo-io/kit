package buildinfo

var (
	Tag string = "development"
)

func GetVersion() string {
	return Tag
}
