// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package remote

import (
	"github.com/m3db/m3/src/query/errors"
	rpc "github.com/m3db/m3/src/query/generated/proto/rpcpb"
	"github.com/m3db/m3/src/query/storage"
)

func decodeTagNamesOnly(
	response *rpc.TagNames,
) *storage.CompleteTagsResult {
	names := response.GetNames()
	tags := make([]storage.CompletedTag, len(names))
	for i, name := range names {
		tags[i] = storage.CompletedTag{Name: name}
	}

	return &storage.CompleteTagsResult{
		CompleteNameOnly: true,
		CompletedTags:    tags,
	}
}

func decodeTagProperties(
	response *rpc.TagValues,
) *storage.CompleteTagsResult {
	values := response.GetValues()
	tags := make([]storage.CompletedTag, len(values))
	for i, value := range values {
		tags[i] = storage.CompletedTag{
			Name:   value.GetKey(),
			Values: value.GetValues(),
		}
	}

	return &storage.CompleteTagsResult{
		CompleteNameOnly: false,
		CompletedTags:    tags,
	}
}

func decodeCompleteTagsResponse(
	response *rpc.CompleteTagsResponse,
) (*storage.CompleteTagsResult, error) {
	if names := response.GetNamesOnly(); names != nil {
		return decodeTagNamesOnly(names), nil
	}

	if props := response.GetDefault(); props != nil {
		return decodeTagProperties(props), nil
	}

	return nil, errors.ErrUnexpectedGRPCResponseType
}

func encodeCompleteTagsRequest(
	query *storage.CompleteTagsQuery,
) (*rpc.CompleteTagsRequest, error) {
	completionType := rpc.CompleteTagsType_DEFAULT
	if query.CompleteNameOnly {
		completionType = rpc.CompleteTagsType_TAGNAME
	}

	matchers, err := encodeTagMatchers(query.TagMatchers)
	if err != nil {
		return nil, err
	}

	return &rpc.CompleteTagsRequest{
		Matchers: &rpc.CompleteTagsRequest_TagMatchers{
			TagMatchers: matchers,
		},
		Options: &rpc.CompleteTagsRequestOptions{
			Type:           completionType,
			FilterNameTags: query.FilterNameTags,
		},
	}, nil
}
