package main

import (
	"fmt"
	"os"

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
	fmt.Println(prototext.Format(fds))
	println()
	println()
	for _, f := range fds.GetFile() {
		fd, err := protodesc.NewFile(f, protoregistry.GlobalFiles)
		if err != nil {
			panic(err)
		}
		protoregistry.GlobalFiles.RegisterFile(fd)
	}

}

func rec(md protoreflect.MessageDescriptor) *descriptorpb.DescriptorProto {

	// TODO create message proto
	rv := &descriptorpb.DescriptorProto{}

	// iterate over fields only add to message proto if pii or
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)

		switch fd.Kind() {
		case protoreflect.StringKind:
			// Check if pii, if Pii add to return type
			opts, ok := proto.GetExtension(fd.Options(), bqproto.E_Bigquery).(*bqproto.BigQueryFieldOptions)
			if ok && opts != nil && opts.GetPolicyTags() == "pii" {

				// Fix up field
				field := protodesc.ToFieldDescriptorProto(fd)
				field.Options = nil
				rv.Field = append(rv.Field, field)

			}
		case protoreflect.MessageKind:
			nestedMessage := rec(fd.Message())
			if nestedMessage != nil {
				name := fd.Message().Name() + "PII"

				// Add nested field
				nestedMessage.Name = proto.String(string(name))
				rv.NestedType = append(rv.NestedType, nestedMessage)

				// Add field
				field := protodesc.ToFieldDescriptorProto(fd)
				field.TypeName = proto.String("." + string(fd.Parent().FullName()+"."+protoreflect.FullName(name)))
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
	rv := rec(m)
	fmt.Println("New PII message that can be added as nestedmessage and field to original message type")
	fmt.Println(prototext.Format(rv))

}
