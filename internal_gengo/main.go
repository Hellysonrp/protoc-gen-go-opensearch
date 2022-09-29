// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package internal_gengo is internal to the protobuf module.
package internal_gengo

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"google.golang.org/protobuf/types/pluginpb"
)

// SupportedFeatures reports the set of supported protobuf language features.
var SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

// GenerateVersionMarkers specifies whether to generate version markers.
var GenerateVersionMarkers = true

const (
	opensearchTypesPackage = protogen.GoImportPath("github.com/Hellysonrp/protoc-gen-go-opensearch/types")
)

var mapMessageProcessed = map[protogen.GoIdent]struct{}{}

// GenerateFile generates the contents of a .pb.go file.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) {
	f := newFileInfo(file)

	filename := file.GeneratedFilenamePrefix + ".pb.opensearch.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P("package ", f.GoPackageName)
	g.P()

	for _, m := range f.Messages {
		processMessage(g, f, m)
	}
}

func processMessage(g *protogen.GeneratedFile, f *fileInfo, m *protogen.Message) {
	if _, ok := mapMessageProcessed[m.GoIdent]; ok {
		return
	}

	if m.Desc.IsMapEntry() {
		return
	}

	opensearchMappingType := g.QualifiedGoIdent(opensearchTypesPackage.Ident("OpensearchMapping"))

	for _, mm := range m.Messages {
		processMessage(g, f, mm)
	}

	// TODO json field name or message field name

	g.P("func (*", m.GoIdent, ") GetOpensearchMappings() map[string]", opensearchMappingType, " {")

	initializedFalseBool := false

	g.P("mapping := make(map[string]", opensearchMappingType, ", ", len(m.Fields), ")")
	for _, ff := range m.Fields {
		if ff.Desc.IsMap() {
			continue
		}

		// too painful to map manually
		if ff.Message != nil && ff.Message.GoIdent.String() == `"\"google.golang.org/protobuf/types/known/structpb\"".Value` {
			continue
		}

		switch ff.Desc.Kind() {
		case protoreflect.BoolKind:
		case protoreflect.EnumKind:
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
			// integer
			g.P("mapping[\"", ff.Desc.JSONName(), "\"] = ", opensearchMappingType, "{")
			g.P("Type: \"integer\",")
			g.P("}")
		case protoreflect.Uint32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
			// long
			g.P("mapping[\"", ff.Desc.JSONName(), "\"] = ", opensearchMappingType, "{")
			g.P("Type: \"long\",")
			g.P("}")
		case protoreflect.Uint64Kind:
			// unsigned_long
			g.P("mapping[\"", ff.Desc.JSONName(), "\"] = ", opensearchMappingType, "{")
			g.P("Type: \"unsigned_long\",")
			g.P("}")
		case protoreflect.FloatKind:
			// float
			g.P("mapping[\"", ff.Desc.JSONName(), "\"] = ", opensearchMappingType, "{")
			g.P("Type: \"float\",")
			g.P("}")
		case protoreflect.DoubleKind:
			// double
			g.P("mapping[\"", ff.Desc.JSONName(), "\"] = ", opensearchMappingType, "{")
			g.P("Type: \"double\",")
			g.P("}")
		case protoreflect.StringKind:
			// text with keyword 256
			g.P("mapping[\"", ff.Desc.JSONName(), "\"] = ", opensearchMappingType, "{")
			g.P("Type: \"text\",")
			g.P("Fields: map[string]", opensearchMappingType, "{")
			g.P("\"keyword\": ", opensearchMappingType, "{")
			g.P("Type: \"keyword\",")
			g.P("IgnoreAbove: 256,")
			g.P("},")
			g.P("},")
			g.P("}")
		case protoreflect.BytesKind:
			// byte with index false
			if !initializedFalseBool {
				g.P("falseBool := false")
				initializedFalseBool = true
			}

			g.P("mapping[\"", ff.Desc.JSONName(), "\"] = ", opensearchMappingType, "{")
			g.P("Type: \"byte\",")
			g.P("Index: &falseBool,")
			g.P("}")
		case protoreflect.MessageKind, protoreflect.GroupKind:
			// nested

			if ff.Message.GoIdent.String() == `"\"google.golang.org/protobuf/types/known/anypb\"".Any` && !initializedFalseBool {
				g.P("falseBool := false")
				initializedFalseBool = true
			}

			g.P("mapping[\"", ff.Desc.JSONName(), "\"] = ", opensearchMappingType, "{")
			g.P("Type: \"nested\",")
			if ff.Message.GoIdent.String() == `"\"google.golang.org/protobuf/types/known/timestamppb\"".Timestamp` {
				g.P("Properties: map[string]", opensearchMappingType, "{")
				g.P("\"seconds\": ", opensearchMappingType, "{")
				g.P("Type: \"long\",")
				g.P("},")
				g.P("\"nanos\": ", opensearchMappingType, "{")
				g.P("Type: \"integer\",")
				g.P("},")
				g.P("},")
			} else if ff.Message.GoIdent.String() == `"\"google.golang.org/protobuf/types/known/anypb\"".Any` {
				g.P("Properties: map[string]", opensearchMappingType, "{")
				g.P("\"type_url\": ", opensearchMappingType, "{")
				g.P("Type: \"text\",")
				g.P("Fields: map[string]", opensearchMappingType, "{")
				g.P("\"keyword\": ", opensearchMappingType, "{")
				g.P("Type: \"keyword\",")
				g.P("IgnoreAbove: 256,")
				g.P("},")
				g.P("},")
				g.P("},")
				g.P("\"value\": ", opensearchMappingType, "{")
				g.P("Type: \"byte\",")
				g.P("Index: &falseBool,")
				g.P("},")
				g.P("},")
			} else {
				g.P("Properties: (&", ff.Message.GoIdent, "{}).GetOpensearchMappings(),")
			}
			g.P("}")
		}
	}
	g.P("return mapping")
	g.P("}")
	g.P()

	mapMessageProcessed[m.GoIdent] = struct{}{}
}
