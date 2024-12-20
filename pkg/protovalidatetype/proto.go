package protovalidatetype

import "google.golang.org/protobuf/proto"

type Validator interface {
	Validate(m proto.Message) error
}
