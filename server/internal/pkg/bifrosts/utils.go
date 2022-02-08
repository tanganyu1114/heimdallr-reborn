package bifrosts

import (
	"gin-vue-admin/pkg/sort_map"
	"github.com/marmotedu/errors"
)

func uintKey(k interface{}) (uint, error) {
	key, ok := k.(uint)
	if !ok {
		return 0, errors.Errorf("k type(%T) is not uint", k)
	}
	return key, nil
}

func checkUINTKeyer(k interface{}) (sort_map.Keyer, error) {
	switch keyer := k.(type) {
	case *uintKeyer:
		return keyer, nil
	case sort_map.Keyer:
		key, err := uintKey(keyer.Key())
		if err != nil {
			return nil, errors.Wrap(err, "k change to uint Keyer failed")
		}

		return NewUINTKeyer(key, keyer.GetOrder()), nil
	default:
		return nil, errors.Errorf("k type(%T) is not Keyer", k)
	}
}
