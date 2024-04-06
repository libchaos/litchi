package bind

import (
	"fmt"
	"net/http"
	"net/textproto"
	"reflect"
	"time"
)

// HeaderField binds a field from the request's header into a value of type T. T can be either a
// primitive type or a [time.Time].
//
// HeaderField consider header as case-insensitive.
//
// If the value can't be bound into T, HeaderField returns Error.
func HeaderField[T primitiveType | time.Time](r *http.Request, header string) (T, error) {
	var (
		target      T
		key         = textproto.CanonicalMIMEHeaderKey(header)
		headerValue = r.Header.Get(key)
		targetValue = reflect.ValueOf(&target).Elem()
	)

	if err := bind(headerValue, targetValue); err != nil {
		return target, fmt.Errorf("%s: %w", header, err)
	}

	return target, nil
}
