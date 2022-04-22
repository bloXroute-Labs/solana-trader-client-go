package connections

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPError struct {
	Code    int         `json:"code"`
	Details interface{} `json:"details"`
	Message string      `json:"message"`
}

func (h HTTPError) Error() string {
	return h.Message
}

// HTTP response for GET request
func HTTPGet[T protoreflect.ProtoMessage](url string, val T) error {
	client := &http.Client{Timeout: time.Second * 7}
	return HTTPGetWithClient[T](url, client, val)
}

func HTTPGetWithClient[T protoreflect.ProtoMessage](url string, client *http.Client, val T) error {
	httpResp, err := client.Get(url)
	if err != nil {
		return err
	}

	if httpResp.StatusCode != http.StatusOK {
		return httpUnmarshalError(httpResp)
	}

	if err := httpUnmarshal[T](httpResp, val); err != nil {
		return err
	}

	return nil
}

func httpUnmarshalError(httpResp *http.Response) error {
	if httpResp == nil {
		return fmt.Errorf("HTTP response is nil")
	}

	var httpError HTTPError
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("error unmarshalling response to HTTPError")
	}

	err = json.Unmarshal(body, &httpError)
	if err != nil {
		return err
	}

	return httpError
}

func httpUnmarshal[T protoreflect.ProtoMessage](httpResp *http.Response, val T) error {
	if httpResp == nil {
		return fmt.Errorf("HTTP response is nil")
	}

	b, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	if err := protojson.Unmarshal(b, val); err != nil {
		return err
	}

	return nil
}
