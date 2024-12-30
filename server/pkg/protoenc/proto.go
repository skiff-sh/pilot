package protoenc

import "google.golang.org/protobuf/encoding/protojson"

var (
	ProtoMarshaller = &protojson.MarshalOptions{
		AllowPartial: true,
	}

	ProtoUnmarshaller = &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
)
