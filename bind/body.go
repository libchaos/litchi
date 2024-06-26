package bind

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	formTag = "form"
	fileTag = "file"
)

var ErrUnsupportedContentType = errors.New("unsupported Content-Type")

// Body binds the request's body into the fields of a struct of type T.
//
// It checks the Content-Type header to select an appropriated parsing method:
//   - "application/json" for JSON parsing
//   - "application/xml" or "text/xml" for XML parsing
//   - "application/x-yaml" for YAML parsing
//   - "application/x-www-form-urlencoded" or "multipart/form-data" for form parsing
//
// Tags from encoding packages, such as "json", "xml" and "yaml" tags, can be used appropriately. For form parsing, use
// the tag "form".
//
// For files inside multipart forms, use the tag "file". Target fields should also be of type
// [*mime/multipart.FileHeader] or [][*mime/multipart.FileHeader].
// The maximum number of bytes stored in memory is 32MB, while the rest is stored in temporary files.
//
// If the Content-Type header is not set, Body defaults to JSON parsing. If it is not supported, it returns
// ErrUnsupportedContentType.
//
// If *T implements [validate.Validatable] (with a pointer receiver), Body calls [validate.Fields] on the result
// and can return [validate.Error].
//
// If T is not a struct type, Body panics.
func Body[T any](r *http.Request) (T, error) {
	var target T

	targetValue := reflect.ValueOf(&target).Elem()

	if targetValue.Kind() != reflect.Struct {
		panic(nonStructTypeParameter)
	}

	if err := bindBody(r, &target, targetValue); err != nil {
		return target, err
	}

	if err := validateFields(&target); err != nil {
		return target, err
	}

	return target, nil
}

func bindBody(r *http.Request, target any, targetValue reflect.Value) error {
	var err error

	contentType := removeFlags(r.Header.Get("Content-Type"))

	switch contentType {
	case "application/xml", "text/xml":
		err = decodeXML(r.Body, target)
	case "application/x-yaml", "text/yaml":
		err = decodeYAML(r.Body, target)
	case "application/x-www-form-urlencoded":
		err = bindForm(r, targetValue)
	case "multipart/form-data":
		err = bindMultipartForm(r, targetValue)
	case "application/json", "":
		err = decodeJSON(r.Body, target)
	default:
		return ErrUnsupportedContentType
	}

	if errors.Is(err, io.EOF) {
		return nil
	}

	return err
}

func removeFlags(contentType string) string {
	i := strings.IndexAny(contentType, "; ")
	if i > 0 {
		return contentType[:i]
	}

	return contentType
}

func bindForm(r *http.Request, targetValue reflect.Value) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	fields := reflect.VisibleFields(targetValue.Type())

	return bindFields(r.Form, formTag, targetValue, fields, bindAll)
}

func bindMultipartForm(r *http.Request, targetValue reflect.Value) error {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return err
	}

	multipartForm := r.MultipartForm

	fields := reflect.VisibleFields(targetValue.Type())

	if err := bindFields(multipartForm.Value, formTag, targetValue, fields, bindAll); err != nil {
		return err
	}

	return bindFields(multipartForm.File, fileTag, targetValue, fields, bindFiles)
}

func decodeJSON(body io.ReadCloser, target any) error {
	err := json.NewDecoder(body).Decode(target)

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return fmt.Errorf("%s: %w", unmarshalTypeError.Field,
			Error{unmarshalTypeError.Value, unmarshalTypeError.Type, nil},
		)
	}

	return err
}

func decodeYAML(body io.ReadCloser, target any) error {
	return yaml.NewDecoder(body).Decode(target)
}

func decodeXML(body io.ReadCloser, target any) error {
	return xml.NewDecoder(body).Decode(target)
}
