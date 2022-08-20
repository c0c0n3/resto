package mime

type MediaType string

const (
	JSON = MediaType("application/json")
	YAML = MediaType("application/yaml")
	GZIP = MediaType("application/gzip")
)

func (t MediaType) String() string {
	return string(t)
}
