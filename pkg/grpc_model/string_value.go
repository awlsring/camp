package grpcmodel

import "google.golang.org/protobuf/types/known/wrapperspb"

func NewStringValue(in *string) *wrapperspb.StringValue {
	if in == nil {
		return nil
	}
	return &wrapperspb.StringValue{Value: *in}
}
