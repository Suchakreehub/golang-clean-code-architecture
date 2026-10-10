package sharedcrud

import (
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AvroType(t reflect.Type) interface{} {
	if t == reflect.TypeOf(primitive.ObjectID{}) {
		return "string"
	} else if t == reflect.TypeOf(primitive.DateTime(0)) {
		return "datetime"
	} else if t == reflect.TypeOf(time.Time{}) {
		return "datetime"
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint:
		return "int"
	case reflect.Uint64:
		return "long"
	case reflect.Bool:
		return "boolean"
	case reflect.Ptr:
		// Recursive call for pointer types
		return AvroType(t.Elem())
	case reflect.Struct:
		// return GenerateAvroSchema(reflect.New(t).Interface())
		return ConvertNameToCRUDSlug(t.Name())
	case reflect.Map:
		return "map"
	default:
		return "unknown"
	}
}

func GenerateAvroSchema(entityModel interface{}) map[string]interface{} {
	var typeModel = reflect.TypeOf(entityModel)

	if typeModel.Kind() == reflect.Ptr {
		typeModel = typeModel.Elem()
	}

	schema := map[string]interface{}{
		"type":   "record",
		"name":   ConvertNameToCRUDSlug(typeModel.Name()),
		"fields": []map[string]interface{}{},
	}

	for i := 0; i < typeModel.NumField(); i++ {
		var field = typeModel.Field(i)
		var jsonTag = field.Tag.Get("json")
		if strings.Contains(jsonTag, "omitempty") {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}
		if jsonTag == "-" {
			continue
		}
		// var avroType = AvroType(field.Type)

		// if field.Type.Kind() == reflect.Struct {
		// 	fieldSchema = GenerateAvroSchema(reflect.New(field.Type).Interface())
		// } else {
		var fieldSchema = map[string]interface{}{
			"name": jsonTag,
			"type": AvroType(field.Type),
		}

		var fieldModel = reflect.New(field.Type).Elem().Interface()

		if _, ok := fieldModel.(interface{ TableName() string }); ok {
			fieldSchema["fields"] = GenerateAvroSchema(fieldModel)
		}
		schema["fields"] = append(schema["fields"].([]map[string]interface{}), fieldSchema)
	}

	return schema
}
