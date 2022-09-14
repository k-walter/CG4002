package common

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"log"
)

func PbToJson(m protoreflect.Message) []byte {
	msg, err := protojson.MarshalOptions{EmitUnpopulated: true}.
		Marshal(m.Interface())
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
