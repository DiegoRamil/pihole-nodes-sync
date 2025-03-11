package deserializers

import (
	"encoding/json"
	"io"
)

func JsonDeserialize(rc io.ReadCloser, v any) error {
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}
