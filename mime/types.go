package mime

type MediaType string

const (
	GZIP         = MediaType("application/gzip")
	JSON         = MediaType("application/json")
	OCTET_STREAM = MediaType("application/octet-stream")
	PLAIN_TEXT   = MediaType("text/plain")
	YAML         = MediaType("application/yaml")
)

func (t MediaType) String() string {
	return string(t)
}
