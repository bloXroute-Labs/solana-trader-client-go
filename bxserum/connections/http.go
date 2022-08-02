package connections

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bloXroute-Labs/serum-client-go/utils"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"io/ioutil"
	"net/http"
	"time"
)

const contentType = "application/json"

var httpResponseNil = fmt.Errorf("HTTP response is nil")

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
	return HTTPGetWithClient[T](url, client, val, utils.AuthHeader)
}

func HTTPGetWithClient[T protoreflect.ProtoMessage](url string, client *http.Client, val T, authHeader string) error {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", authHeader)
	httpResp, err := client.Do(req)
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

func HTTPPostWithClient[T protoreflect.ProtoMessage](url string, client *http.Client, body interface{}, val T, authHeader string) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", contentType)
	httpResp, err := client.Do(req)
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
		return httpResponseNil
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	return errors.New(string(body))
}

func httpUnmarshal[T protoreflect.ProtoMessage](httpResp *http.Response, val T) error {
	if httpResp == nil {
		return httpResponseNil
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
