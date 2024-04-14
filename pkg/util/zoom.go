package util

import (
	"fmt"
	"strconv"
	"time"
)

type FallbackMarshaler struct{}

func (f FallbackMarshaler) Marshal(v interface{}) ([]byte, error) {
	switch x := v.(type) {
	case time.Time:
		return []byte(strconv.Itoa(int(x.Unix()))), nil
	}

	return nil, fmt.Errorf("unsupported type %T", v)
}

func (f FallbackMarshaler) Unmarshal(data []byte, v interface{}) error {
	switch x := v.(type) {
	case *time.Time:
		i, err := strconv.Atoi(string(data))
		if err != nil {
			return err
		}
		*x = time.Unix(int64(i), 0)
		return nil
	}

	return fmt.Errorf("unsupported type %T", v)
}
