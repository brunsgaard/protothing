package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc" //newspb "google.golang.org/protobuf/internal/testprotos/news"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	bqproto "github.com/GoogleCloudPlatform/protoc-gen-bq-schema/protos"
)

func registerproto() {
	data, _ := os.ReadFile("addressbook.bin")
	fds := &descriptorpb.FileDescriptorSet{}
	proto.Unmarshal(data, fds)
	for _, f := range fds.GetFile() {
		fd, err := protodesc.NewFile(f, protoregistry.GlobalFiles)
		if err != nil {
			panic(err)
		}
		protoregistry.GlobalFiles.RegisterFile(fd)
	}

}

func BuildPIIMessageType(md protoreflect.MessageDescriptor, descriptorFqdn string) *descriptorpb.DescriptorProto {
	rv := &descriptorpb.DescriptorProto{Name: proto.String(descriptorFqdn[strings.LastIndex(descriptorFqdn, ".")+1:])}
	// Iterate over each field in message, add field to rv if field (directly or
	// recursivly) contains PII tag
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		switch fd.Kind() {
		case protoreflect.StringKind, protoreflect.BytesKind:
			// Check if field contains PII, if this is the case clear Options and add
			// field to rv
			opts, ok := proto.GetExtension(fd.Options(), bqproto.E_Bigquery).(*bqproto.BigQueryFieldOptions)
			if ok && opts != nil && opts.GetPolicyTags() == "pii" {
				field := protodesc.ToFieldDescriptorProto(fd)
				field.Options = nil
				rv.Field = append(rv.Field, field)
			}
		case protoreflect.MessageKind:
			// If field is kind message, call BuildPIIMessageType recursivly
			messageName := string(fd.Message().Name()) + strconv.Itoa(i)
			fqdnName := descriptorFqdn + "." + messageName
			messageType := BuildPIIMessageType(fd.Message(), fqdnName)
			if len(messageType.Field) != 0 {
				// Add message type as nested message to rv and add field to rv
				rv.NestedType = append(rv.NestedType, messageType)
				field := protodesc.ToFieldDescriptorProto(fd)
				field.TypeName = proto.String(fqdnName)
				rv.Field = append(rv.Field, field)
			}
		}
	}
	return rv
}

func main() {
	registerproto()
	desc, err := protoregistry.GlobalFiles.FindDescriptorByName("tutorial.Person")
	if err != nil {
		panic(err)
	}
	m := desc.(protoreflect.MessageDescriptor)
	rv := BuildPIIMessageType(m, ".tutorial.v1.Address.Metadata.PIIFields")
	fmt.Println(prototext.Format(rv))
}
