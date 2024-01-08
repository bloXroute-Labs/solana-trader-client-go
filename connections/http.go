package connections

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
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

func HTTPGetWithClient[T protoreflect.ProtoMessage](ctx context.Context, url string, client *http.Client, val T, authHeader string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("X-SDK", "solana-trader-client-go")
	req.Header.Set("X-SDK-Version", utils.Version())
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

func HTTPPostWithClient[T protoreflect.ProtoMessage](ctx context.Context, url string, client *http.Client, body interface{}, val T, authHeader string) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(b))
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
