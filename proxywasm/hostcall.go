// Copyright 2020-2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxywasm

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func GetPluginConfiguration(size int) ([]byte, error) {
	return getBuffer(internal.BufferTypePluginConfiguration, 0, size)
}

func GetVMConfiguration(size int) ([]byte, error) {
	return getBuffer(internal.BufferTypeVMConfiguration, 0, size)
}

func SendHttpResponse(statusCode uint32, headers types.Headers, body []byte) error {
	shs := internal.SerializeMap(headers)
	var bp *byte
	if len(body) > 0 {
		bp = &body[0]
	}
	hp := &shs[0]
	hl := len(shs)
	return internal.StatusToError(
		internal.ProxySendLocalResponse(
			statusCode, nil, 0,
			bp, len(body), hp, hl, -1,
		),
	)
}

func SetTickPeriodMilliSeconds(millSec uint32) error {
	return internal.StatusToError(internal.ProxySetTickPeriodMilliseconds(millSec))
}

func DispatchHttpCall(upstream string,
	headers types.Headers, body string, trailers types.Trailers,
	timeoutMillisecond uint32, callBack func(numHeaders, bodySize, numTrailers int)) (calloutID uint32, err error) {
	shs := internal.SerializeMap(headers)
	hp := &shs[0]
	hl := len(shs)

	sts := internal.SerializeMap(trailers)
	tp := &sts[0]
	tl := len(sts)

	u := internal.StringBytePtr(upstream)
	switch st := internal.ProxyHttpCall(u, len(upstream),
		hp, hl, internal.StringBytePtr(body), len(body), tp, tl, timeoutMillisecond, &calloutID); st {
	case internal.StatusOK:
		internal.RegisterHttpCallout(calloutID, callBack)
		return calloutID, nil
	default:
		return 0, internal.StatusToError(st)
	}
}

func GetHttpCallResponseHeaders() (types.Headers, error) {
	return getMap(internal.MapTypeHttpCallResponseHeaders)
}

func GetHttpCallResponseBody(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeHttpCallResponseBody, start, maxSize)
}

func GetHttpCallResponseTrailers() (types.Trailers, error) {
	return getMap(internal.MapTypeHttpCallResponseTrailers)
}

func CallForeignFunction(funcName string, param []byte) (ret []byte, err error) {
	f := internal.StringBytePtr(funcName)

	var returnData *byte
	var returnSize int

	switch st := internal.ProxyCallForeignFunction(f, len(funcName), &param[0], len(param), &returnData, &returnSize); st {
	case internal.StatusOK:
		return internal.RawBytePtrToByteSlice(returnData, returnSize), nil
	default:
		return nil, internal.StatusToError(st)
	}
}

func GetDownstreamData(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeDownstreamData, start, maxSize)
}

func AppendDownstreamData(data []byte) error {
	return appendToBuffer(internal.BufferTypeDownstreamData, data)
}

func PrependDownstreamData(data []byte) error {
	return prependToBuffer(internal.BufferTypeDownstreamData, data)
}

func ReplaceDownstreamData(data []byte) error {
	return replaceBuffer(internal.BufferTypeDownstreamData, data)
}

func GetUpstreamData(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeUpstreamData, start, maxSize)
}

func AppendUpstreamData(data []byte) error {
	return appendToBuffer(internal.BufferTypeUpstreamData, data)
}

func PrependUpstreamData(data []byte) error {
	return prependToBuffer(internal.BufferTypeUpstreamData, data)
}

func ReplaceUpstreamData(data []byte) error {
	return replaceBuffer(internal.BufferTypeUpstreamData, data)
}

func ContinueDownstream() error {
	return internal.StatusToError(internal.ProxyContinueStream(internal.StreamTypeDownstream))
}

func ContinueUpstream() error {
	return internal.StatusToError(internal.ProxyContinueStream(internal.StreamTypeUpstream))
}

func CloseDownstream() error {
	return internal.StatusToError(internal.ProxyCloseStream(internal.StreamTypeDownstream))
}

func CloseUpstream() error {
	return internal.StatusToError(internal.ProxyCloseStream(internal.StreamTypeUpstream))
}

func GetHttpRequestHeaders() (types.Headers, error) {
	return getMap(internal.MapTypeHttpRequestHeaders)
}

func SetHttpRequestHeaders(headers types.Headers) error {
	return setMap(internal.MapTypeHttpRequestHeaders, headers)
}

func GetHttpRequestHeader(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpRequestHeaders, key)
}

func RemoveHttpRequestHeader(key string) error {
	return removeMapValue(internal.MapTypeHttpRequestHeaders, key)
}

func SetHttpRequestHeader(key, value string) error {
	return setMapValue(internal.MapTypeHttpRequestHeaders, key, value)
}

func AddHttpRequestHeader(key, value string) error {
	return addMapValue(internal.MapTypeHttpRequestHeaders, key, value)
}

func GetHttpRequestBody(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeHttpRequestBody, start, maxSize)
}

func AppendHttpRequestBody(data []byte) error {
	return appendToBuffer(internal.BufferTypeHttpRequestBody, data)
}

func PrependHttpRequestBody(data []byte) error {
	return prependToBuffer(internal.BufferTypeHttpRequestBody, data)
}

func ReplaceHttpRequestBody(data []byte) error {
	return replaceBuffer(internal.BufferTypeHttpRequestBody, data)
}

func GetHttpRequestTrailers() (types.Trailers, error) {
	return getMap(internal.MapTypeHttpRequestTrailers)
}

func SetHttpRequestTrailers(trailers types.Trailers) error {
	return setMap(internal.MapTypeHttpRequestTrailers, trailers)
}

func GetHttpRequestTrailer(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpRequestTrailers, key)
}

func RemoveHttpRequestTrailer(key string) error {
	return removeMapValue(internal.MapTypeHttpRequestTrailers, key)
}

func SetHttpRequestTrailer(key, value string) error {
	return setMapValue(internal.MapTypeHttpRequestTrailers, key, value)
}

func AddHttpRequestTrailer(key, value string) error {
	return addMapValue(internal.MapTypeHttpRequestTrailers, key, value)
}

func ResumeHttpRequest() error {
	return internal.StatusToError(internal.ProxyContinueStream(internal.StreamTypeRequest))
}

func GetHttpResponseHeaders() (types.Headers, error) {
	return getMap(internal.MapTypeHttpResponseHeaders)
}

func SetHttpResponseHeaders(headers types.Headers) error {
	return setMap(internal.MapTypeHttpResponseHeaders, headers)
}

func GetHttpResponseHeader(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpResponseHeaders, key)
}

func RemoveHttpResponseHeader(key string) error {
	return removeMapValue(internal.MapTypeHttpResponseHeaders, key)
}

func SetHttpResponseHeader(key, value string) error {
	return setMapValue(internal.MapTypeHttpResponseHeaders, key, value)
}

func AddHttpResponseHeader(key, value string) error {
	return addMapValue(internal.MapTypeHttpResponseHeaders, key, value)
}

func GetHttpResponseBody(start, maxSize int) ([]byte, error) {
	return getBuffer(internal.BufferTypeHttpResponseBody, start, maxSize)
}

func AppendHttpResponseBody(data []byte) error {
	return appendToBuffer(internal.BufferTypeHttpResponseBody, data)
}

func PrependHttpResponseBody(data []byte) error {
	return prependToBuffer(internal.BufferTypeHttpResponseBody, data)
}

func ReplaceHttpResponseBody(data []byte) error {
	return replaceBuffer(internal.BufferTypeHttpResponseBody, data)
}

func GetHttpResponseTrailers() (types.Trailers, error) {
	return getMap(internal.MapTypeHttpResponseTrailers)
}

func SetHttpResponseTrailers(trailers types.Trailers) error {
	return setMap(internal.MapTypeHttpResponseTrailers, trailers)
}

func GetHttpResponseTrailer(key string) (string, error) {
	return getMapValue(internal.MapTypeHttpResponseTrailers, key)
}

func RemoveHttpResponseTrailer(key string) error {
	return removeMapValue(internal.MapTypeHttpResponseTrailers, key)
}

func SetHttpResponseTrailer(key, value string) error {
	return setMapValue(internal.MapTypeHttpResponseTrailers, key, value)
}

func AddHttpResponseTrailer(key, value string) error {
	return addMapValue(internal.MapTypeHttpResponseTrailers, key, value)
}

func ResumeHttpResponse() error {
	return internal.StatusToError(internal.ProxyContinueStream(internal.StreamTypeResponse))
}

func RegisterSharedQueue(name string) (uint32, error) {
	var queueID uint32
	ptr := internal.StringBytePtr(name)
	st := internal.ProxyRegisterSharedQueue(ptr, len(name), &queueID)
	return queueID, internal.StatusToError(st)
}

func ResolveSharedQueue(vmID, queueName string) (uint32, error) {
	var ret uint32
	st := internal.ProxyResolveSharedQueue(internal.StringBytePtr(vmID),
		len(vmID), internal.StringBytePtr(queueName), len(queueName), &ret)
	return ret, internal.StatusToError(st)
}

func DequeueSharedQueue(queueID uint32) ([]byte, error) {
	var raw *byte
	var size int
	st := internal.ProxyDequeueSharedQueue(queueID, &raw, &size)
	if st != internal.StatusOK {
		return nil, internal.StatusToError(st)
	}
	return internal.RawBytePtrToByteSlice(raw, size), nil
}

func EnqueueSharedQueue(queueID uint32, data []byte) error {
	return internal.StatusToError(internal.ProxyEnqueueSharedQueue(queueID, &data[0], len(data)))
}

func GetSharedData(key string) (value []byte, cas uint32, err error) {
	var raw *byte
	var size int

	st := internal.ProxyGetSharedData(internal.StringBytePtr(key), len(key), &raw, &size, &cas)
	if st != internal.StatusOK {
		return nil, 0, internal.StatusToError(st)
	}
	return internal.RawBytePtrToByteSlice(raw, size), cas, nil
}

func SetSharedData(key string, data []byte, cas uint32) error {
	st := internal.ProxySetSharedData(internal.StringBytePtr(key),
		len(key), &data[0], len(data), cas)
	return internal.StatusToError(st)
}

func GetProperty(path []string) ([]byte, error) {
	var ret *byte
	var retSize int
	raw := internal.SerializePropertyPath(path)

	err := internal.StatusToError(internal.ProxyGetProperty(&raw[0], len(raw), &ret, &retSize))
	if err != nil {
		return nil, err
	}

	return internal.RawBytePtrToByteSlice(ret, retSize), nil

}

func SetProperty(path []string, data []byte) error {
	raw := internal.SerializePropertyPath(path)
	return internal.StatusToError(internal.ProxySetProperty(
		&raw[0], len(path), &data[0], len(data),
	))
}

func Done() {
	internal.ProxyDone()
}
