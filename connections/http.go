package connections

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	package_info "github.com/bloXroute-Labs/solana-trader-client-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var httpResponseNil = fmt.Errorf("HTTP response is nil")

type HTTPError struct {
	Code    int         `json:"code"`
	Details interface{} `json:"details"`
	Message string      `json:"message"`
}

func (h HTTPError) Error() string {
	return h.Message
}

var GlobalCorrelationID = ""
var GlobalPostSubmitContentType = "application/json"

func HTTPGetWithClient[T protoreflect.ProtoMessage](ctx context.Context, url string, client *http.Client, val T, authHeader string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("x-sdk", package_info.Name)
	req.Header.Set("x-sdk-version", package_info.Version)
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
	protoMsg := body.(proto.Message)

	// Use protojson marshaler
	marshaler := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	b, err := marshaler.Marshal(protoMsg)
	if err != nil {
		return err
	}

	return HTTPPostWithClientRaw[T](ctx, url, client, b, val, authHeader, "application/json", GlobalCorrelationID)
}

func HTTPPostWithClientRaw[T protoreflect.ProtoMessage](ctx context.Context, url string, client *http.Client, body []byte, val T, authHeader, contentType, correlationID string) error {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-sdk", package_info.Name)
	req.Header.Set("x-sdk-version", package_info.Version)
	if len(correlationID) > 0 {
		req.Header.Set("X-CORRELATION-ID", correlationID)
	}
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
