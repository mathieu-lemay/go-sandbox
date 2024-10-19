package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/rs/zerolog/log"

	"github.com/go-playground/validator/v10"
	"go-sandbox/logging"
)

type Animal struct {
	Name string `json:"name" validate:"required"`
}

type Person struct {
	Name     string  `json:"name" validate:"required"`
	Nickname *string `json:"nickname" validate:"required"`
	Age      int     `json:"age" validate:"required"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func main() {
	logging.ConfigureLogger(
	//logging.WithLevel(logging.LevelDebug),
	)

	//jsonData := []byte(`{"name": "John Smith", "age": 0}`)
	jsonData := []byte(`{"name": "John"}`)
	var person Person
	if err := parseAndValidate(&person, jsonData); err != nil {
		log.Fatal().Err(err).Msg("Failed to validate struct")
	}
}

var structsCache = map[reflect.Type]reflect.Type{}

func createPointerStruct(src interface{}) reflect.Type {
	srcType := reflect.TypeOf(src)
	if ptrType, ok := structsCache[srcType]; ok {
		return ptrType
	}

	var sfs []reflect.StructField
	var t reflect.Type

	vt := reflect.ValueOf(src)
	if vt.Kind() == reflect.Ptr {
		t = vt.Elem().Type()
	} else {
		t = vt.Type()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		switch field.Type.Kind() {
		case reflect.Struct:
			interfaceType := createPointerStruct(reflect.New(field.Type).Elem().Interface())
			sf := reflect.StructField{
				Name: field.Name,
				Type: reflect.PointerTo(interfaceType),
				Tag:  field.Tag,
			}
			sfs = append(sfs, sf)
		case reflect.Ptr:
			sf := reflect.StructField{
				Name: field.Name,
				Type: field.Type,
				Tag:  field.Tag,
			}
			sfs = append(sfs, sf)
		default:
			sf := reflect.StructField{
				Name: field.Name,
				Type: reflect.PointerTo(field.Type),
				Tag:  field.Tag,
			}
			sfs = append(sfs, sf)
		}

	}

	ptrType := reflect.StructOf(sfs)
	structsCache[srcType] = ptrType

	return ptrType
}

func copyStruct2(st reflect.Type, src reflect.Value) interface{} {
	so := reflect.New(st)
	dVal := src
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		log.Debug().Interface("field", field.Name).Msg("field")
		vValue := dVal.FieldByName(field.Name)
		tValue := vValue.Type()
		vType := field.Type
		if tValue.AssignableTo(vType) {
			log.Debug().Interface("vValue", vValue).Msg("AssignableTo")
			so.Elem().FieldByName(field.Name).Set(vValue)
		} else if vValue.CanConvert(vType) {
			log.Debug().Interface("vValue", vValue).Msg("CanConvert")
			so.Elem().FieldByName(field.Name).Set(vValue.Convert(vType))
		} else {
			log.Debug().Interface("vValue", vValue).Msg("Pointer")
			if vValue.Kind() == reflect.Struct {
				if field.Type.Kind() == reflect.Ptr {
					vStruct := copyStruct2(field.Type.Elem(), vValue)
					so.Elem().FieldByName(field.Name).Set(reflect.ValueOf(&vStruct))
					continue
				}
				vStruct := copyStruct2(field.Type, vValue)
				so.Elem().FieldByName(field.Name).Set(reflect.ValueOf(vStruct))
				continue
			}
			so.Elem().FieldByName(field.Name).Set(reflect.NewAt(vType.Elem(), unsafe.Pointer(vValue.UnsafeAddr())))
		}

	}
	log.Debug().Interface("so", so.Interface()).Msg("so")
	return so.Interface()
}

func copyStruct(src interface{}, dst interface{}) {
	//so := reflect.New(st)
	//dVal := src
	//for i := 0; i < st.NumField(); i++ {
	//	field := st.Field(i)
	//	log.Debug().Interface("field", field.Name).Msg("field")
	//	vValue := dVal.FieldByName(field.Name)
	//	tValue := vValue.Type()
	//	vType := field.Type
	//	if tValue.AssignableTo(vType) {
	//		log.Debug().Interface("vValue", vValue).Msg("AssignableTo")
	//		so.Elem().FieldByName(field.Name).Set(vValue)
	//	} else if vValue.CanConvert(vType) {
	//		log.Debug().Interface("vValue", vValue).Msg("CanConvert")
	//		so.Elem().FieldByName(field.Name).Set(vValue.Convert(vType))
	//	} else {
	//		log.Debug().Interface("vValue", vValue).Msg("Pointer")
	//		if vValue.Kind() == reflect.Struct {
	//			if field.Type.Kind() == reflect.Ptr {
	//				vStruct := copyStruct(field.Type.Elem(), vValue)
	//				so.Elem().FieldByName(field.Name).Set(reflect.ValueOf(&vStruct))
	//				continue
	//			}
	//			vStruct := copyStruct(field.Type, vValue)
	//			so.Elem().FieldByName(field.Name).Set(reflect.ValueOf(vStruct))
	//			continue
	//		}
	//		so.Elem().FieldByName(field.Name).Set(reflect.NewAt(vType.Elem(), unsafe.Pointer(vValue.UnsafeAddr())))
	//	}
	//
	//}
	//log.Debug().Interface("so", so.Interface()).Msg("so")

	return
}

func parseAndValidate(target any, data []byte) error {
	// Create copy of the given type, but with all fields as pointers
	ptrType := createPointerStruct(target)

	// Something something interface, something something golang
	// https://stackoverflow.com/a/45680060
	ptrInstance := reflect.New(ptrType).Interface()

	if err := json.Unmarshal(data, &ptrInstance); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if err := validate.Struct(ptrInstance); err != nil {
		return fmt.Errorf("invalid data: %w", err)
	}

	// Fill the original struct with validated values
	copyStruct(ptrInstance, &target)

	return nil
}
