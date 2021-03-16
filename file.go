package remote_logger

type file struct {
	bytes   []byte
	name    string
	caption string
}

func newFileBytes(name, caption string, fileBytes []byte) *file {
	return &file{
		bytes:   fileBytes,
		name:    name,
		caption: caption,
	}
}
