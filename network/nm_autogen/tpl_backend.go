package main

// file header
const fileHeader = `// This file is automatically generated, please don't edit manully.
package main

import (
	"fmt"
)
`

// get key type
const tplGetKeyType = `
// Get key type
func get{{.FieldName | ToFieldFuncBaseName}}KeyType(key string) (t ktype) {
	switch key {
	default:
		t = ktypeUnknown{{range .Keys}}
	case {{.Name}}:
		t = {{.Type}}{{end}}
	}
	return
}
`

// check is key in current field
const tplIsKeyInSettingField = `
// Check is key in current setting field
func isKeyIn{{.FieldName | ToFieldFuncBaseName}}(key string) bool {
	switch key { {{range .Keys}}{{if .UsedByBackEnd}}
	case {{.Name}}:
		return true{{end}}{{end}}
	}
	return false
}
`

// Ensure field and key exists and not empty
const tplEnsureNoEmpty = `{{$fieldFuncBaseName := .FieldName | ToFieldFuncBaseName}}{{$fieldName := .FieldName}}
// Ensure field and key exists and not empty
func ensureField{{$fieldFuncBaseName}}Exists(data connectionData, errs fieldErrors, relatedKey string) {
	if !isSettingFieldExists(data, {{.FieldName}}) {
		rememberError(errs, relatedKey, {{.FieldName}}, fmt.Sprintf(NM_KEY_ERROR_MISSING_SECTION, {{.FieldName}}))
	}
	fieldData, _ := data[{{.FieldName}}]
	if len(fieldData) == 0 {
		rememberError(errs, relatedKey, {{.FieldName}}, fmt.Sprintf(NM_KEY_ERROR_EMPTY_SECTION, {{.FieldName}}))
	}
}{{range $i, $key := .Keys}}{{if $key.UsedByBackEnd}}{{$keyFuncBaseName := $key.Name | ToKeyFuncBaseName}}
func ensure{{$keyFuncBaseName}}NoEmpty(data connectionData, errs fieldErrors) {
	if !is{{$keyFuncBaseName}}Exists(data) {
		rememberError(errs, {{$fieldName}}, {{$key.Name}}, NM_KEY_ERROR_MISSING_VALUE)
	}{{if IfNeedCheckValueLength $key.Type}}
	value := get{{$keyFuncBaseName}}(data)
	if len(value) == 0 {
		rememberError(errs, {{$fieldName}}, {{$key.Name}}, NM_KEY_ERROR_EMPTY_VALUE)
	}{{end}}
}{{end}}{{end}}
`

// get key's default json value
const tplGetDefaultValueJSON = `{{$fieldFuncBaseName := .FieldName | ToFieldFuncBaseName}}
// Get key's default value
func get{{$fieldFuncBaseName}}KeyDefaultValueJSON(key string) (valueJSON string) {
	switch key {
	default:
		logger.Error("invalid key:", key){{range .Keys}}{{if .UsedByBackEnd}}{{$default := ToKeyTypeDefaultValue .Name}}
	case {{.Name}}:
		valueJSON = ` + "`{{$default}}`" + `{{end}}{{end}}
	}
	return
}
`

// get json value generally
const tplGeneralGetterJSON = `{{$fieldFuncBaseName := .FieldName | ToFieldFuncBaseName}}
// Get JSON value generally
func generalGet{{$fieldFuncBaseName}}KeyJSON(data connectionData, key string) (value string) {
	switch key {
	default:
		logger.Error("generalGet{{$fieldFuncBaseName}}KeyJSON: invalide key", key){{range .Keys}}{{if .UsedByBackEnd}}
	case {{.Name}}:
		value = get{{.Name | ToKeyFuncBaseName}}JSON(data){{end}}{{end}}
	}
	return
}
`

// set json value generally
const tplGeneralSetterJSON = `{{$fieldFuncBaseName := .FieldName | ToFieldFuncBaseName}}
// Set JSON value generally
func generalSet{{$fieldFuncBaseName}}KeyJSON(data connectionData, key, valueJSON string) (err error) {
	switch key {
	default:
		logger.Error("generalSet{{$fieldFuncBaseName}}KeyJSON: invalide key", key){{range .Keys}}{{if .UsedByBackEnd}}
	case {{.Name}}:
		err = {{if .LogicSet}}logicSet{{else}}set{{end}}{{.Name | ToKeyFuncBaseName}}JSON(data, valueJSON){{end}}{{end}}
	}
	return
}
`

// check if key exists
const tplCheckExists = `
// Check if key exists{{$fieldName := .FieldName}}{{range $i, $key := .Keys}}{{if $key.UsedByBackEnd}}
func is{{$key.Name | ToKeyFuncBaseName}}Exists(data connectionData) bool {
	return isSettingKeyExists(data, {{$fieldName}}, {{$key.Name}})
}{{end}}{{end}}
`

// getter
const tplGetter = `
// Getter{{$fieldName := .FieldName}}{{range $i, $key := .Keys}}{{if $key.UsedByBackEnd}}
func get{{$key.Name | ToKeyFuncBaseName}}(data connectionData) (value {{$key.Type | ToKeyTypeRealData}}) {
	value, _ = getSettingKey(data, {{$fieldName}}, {{$key.Name}}).({{$key.Type | ToKeyTypeRealData}})
	return
}{{end}}{{end}}
`

// setter
const tplSetter = `
// Setter{{$fieldName := .FieldName}}{{range $i, $key := .Keys}}{{if $key.UsedByBackEnd}}
func set{{$key.Name | ToKeyFuncBaseName}}(data connectionData, value {{$key.Type | ToKeyTypeRealData}}) {
	setSettingKey(data, {{$fieldName}}, {{$key.Name}}, value)
}{{end}}{{end}}
`

// json getter
const tplJSONGetter = `
// JSON Getter{{$fieldName := .FieldName}}{{range $i, $key := .Keys}}{{if $key.UsedByBackEnd}}
func get{{$key.Name | ToKeyFuncBaseName}}JSON(data connectionData) (valueJSON string) {
	valueJSON = getSettingKeyJSON(data, {{$fieldName}}, {{$key.Name}}, get{{$fieldName | ToFieldFuncBaseName}}KeyType({{$key.Name}}))
	return
}{{end}}{{end}}
`

// json setter
const tplJSONSetter = `
// JSON Setter{{$fieldName := .FieldName}}{{range $i, $key := .Keys}}{{if $key.UsedByBackEnd}}
func set{{$key.Name | ToKeyFuncBaseName}}JSON(data connectionData, valueJSON string) (err error) {
	return setSettingKeyJSON(data, {{$fieldName}}, {{$key.Name}}, valueJSON, get{{$fieldName | ToFieldFuncBaseName}}KeyType({{$key.Name}}))
}{{end}}{{end}}
`

// logic json setter
const tplLogicJSONSetter = `
// Logic JSON Setter{{range $i, $key := .Keys}}{{if $key.LogicSet}}{{$keyFuncBaseName := $key.Name | ToKeyFuncBaseName}}
func logicSet{{$keyFuncBaseName}}JSON(data connectionData, valueJSON string) (err error) {
	err = set{{$keyFuncBaseName}}JSON(data, valueJSON)
	if err != nil {
		return
	}
	if is{{$keyFuncBaseName}}Exists(data) {
		value := get{{$keyFuncBaseName}}(data)
		err = logicSet{{$keyFuncBaseName}}(data, value)
	}
	return
}{{end}}{{end}}
`

// remover
const tplRemover = `
// Remover{{$fieldName := .FieldName}}{{range $i, $key := .Keys}}{{if $key.UsedByBackEnd}}
func remove{{$key.Name | ToKeyFuncBaseName}}(data connectionData) {
	removeSettingKey(data, {{$fieldName}}, {{$key.Name}})
}{{end}}{{end}}
`

// TODO
const tplGetAvaiableValues = `// Get avaiable values`

// general setting utils
const tplGeneralSettingUtils = `// This file is automatically generated, please don't edit manully.
package main

func generalIsKeyInSettingField(field, key string) bool {
	if isVirtualKey(field, key) {
		return true
	}
	switch field {
	default:
		logger.Warning("invalid field name", field){{range .}}
	case {{.FieldName}}:
		return isKeyIn{{.FieldName | ToFieldFuncBaseName}}(key){{end}}
	}
	return false
}

func generalGetSettingKeyType(field, key string) (t ktype) {
	if isVirtualKey(field, key) {
		t = getSettingVkKeyType(field, key)
		return
	}
	switch field {
	default:
		logger.Warning("invalid field name", field){{range .}}
	case {{.FieldName}}:
		t = get{{.FieldName | ToFieldFuncBaseName}}KeyType(key){{end}}
	}
	return
}

func generalGetSettingAvailableKeys(data connectionData, field string) (keys []string) {
	switch field { {{range .}}
	case {{.FieldName}}:
		keys = get{{.FieldName | ToFieldFuncBaseName}}AvailableKeys(data){{end}}
	}
	return
}

func generalGetSettingAvailableValues(data connectionData, field, key string) (values []kvalue) {
	if isVirtualKey(field, key) {
		values = generalGetSettingVkAvailableValues(data, field, key)
		return
	}
	switch field { {{range .}}
	case {{.FieldName}}:
		values = get{{.FieldName | ToFieldFuncBaseName}}AvailableValues(data, key){{end}}
	}
	return
}

func generalCheckSettingValues(data connectionData, field string) (errs fieldErrors) {
	switch field {
	default:
		logger.Error("invalid field name", field){{range .}}
	case {{.FieldName}}:
		errs = check{{.FieldName | ToFieldFuncBaseName}}Values(data){{end}}
	}
	return
}

func generalGetSettingKeyJSON(data connectionData, field, key string) (valueJSON string) {
	if isVirtualKey(field, key) {
		valueJSON = generalGetVirtualKeyJSON(data, field, key)
		return
	}
	switch field {
	default:
		logger.Warning("invalid field name", field){{range.}}
	case {{.FieldName}}:
		valueJSON = generalGet{{.FieldName | ToFieldFuncBaseName}}KeyJSON(data, key){{end}}
	}
	return
}

func generalSetSettingKeyJSON(data connectionData, field, key, valueJSON string) (err error) {
	if isVirtualKey(field, key) {
		err = generalSetVirtualKeyJSON(data, field, key, valueJSON)
		return
	}
	switch field {
	default:
		logger.Warning("invalid field name", field){{range .}}
	case {{.FieldName}}:
		err = generalSet{{.FieldName | ToFieldFuncBaseName}}KeyJSON(data, key, valueJSON){{end}}
	}
	return
}

func getSettingKeyDefaultValueJSON(field, key string) (valueJSON string) {
	switch field {
	default:
		logger.Warning("invalid field name", field){{range .}}
	case {{.FieldName}}:
		valueJSON = get{{.FieldName | ToFieldFuncBaseName}}KeyDefaultValueJSON(key){{end}}
	}
	return
}`

// virtual key
const tplVirtualKey = `// This file is automatically generated, please don't edit manully.
package main{{$vks := .}}

// All virtual keys data
var virtualKeys = []virtualKey{ {{range .}}
	virtualKey{ Name:{{.Name}}, Type:{{.Type}}, RelatedField:{{.RelatedField}}, RelatedKey:{{.RelatedKey}}, EnableWrapper:{{.EnableWrapper}}, Available:{{.UsedByFrontEnd}}, Optional:{{.Optional}} },{{end}}
}

// Get JSON value generally
func generalGetVirtualKeyJSON(data connectionData, field, key string) (valueJSON string) {
	switch field { {{range $i, $field := GetAllVkFields $vks}}
	case {{$field}}:
		switch key { {{range $i, $key := GetAllVkFieldKeys $vks $field}}
		case {{$key}}:
			return get{{$key | ToKeyFuncBaseName}}JSON(data){{end}}
		}{{end}}
	}
	logger.Error("invalid virtual key:", field, key)
	return
}

// Set JSON value generally
func generalSetVirtualKeyJSON(data connectionData, field, key string, valueJSON string) (err error) {
	// each virtual key has a logic setter
	switch field { {{range $i, $field := GetAllVkFields $vks}}
	case {{$field}}:
		switch key { {{range $i, $key := GetAllVkFieldKeys $vks $field}}
		case {{$key}}:
			err = logicSet{{$key | ToKeyFuncBaseName}}JSON(data, valueJSON)
			return{{end}}
		}{{end}}
	}
	logger.Error("invalid virtual key:", field, key)
	return
}

// JSON getter{{range $i, $vk := $vks}}{{$keyBaseFuncName := $vk.Name | ToKeyFuncBaseName}}
func get{{$keyBaseFuncName}}JSON(data connectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(get{{$keyBaseFuncName}}(data))
	return
}{{end}}

// Logic JSON setter{{range $i, $vk := $vks}}{{$keyBaseFuncName := $vk.Name | ToKeyFuncBaseName}}
func logicSet{{$keyBaseFuncName}}JSON(data connectionData, valueJSON string) (err error) {
	value, _ := jsonToKeyValue{{$vk.Type | ToKeyTypeShortName}}(valueJSON)
	return logicSet{{$keyBaseFuncName}}(data, value)
}{{end}}

// Getter for enable key wrapper{{range $i, $vk := $vks}}{{if $vk.EnableWrapper}}{{$keyBaseFuncName := $vk.Name | ToKeyFuncBaseName}}
func get{{$keyBaseFuncName}}(data connectionData) (value bool) {
	if is{{$vk.RelatedKey | ToKeyFuncBaseName}}Exists(data) {
		return true
	}
	return false
}{{end}}{{end}}

// Setter for enable key wrapper{{range $i, $vk := $vks}}{{if $vk.EnableWrapper}}{{$keyBaseFuncName := $vk.Name | ToKeyFuncBaseName}}
func logicSet{{$keyBaseFuncName}}(data connectionData, value bool) (err error) {
	if value {
		set{{$vk.RelatedKey | ToKeyFuncBaseName}}(data, {{$vk.RelatedKey | ToKeyTypeDefaultValue}})
	} else {
		remove{{$vk.RelatedKey | ToKeyFuncBaseName}}(data)
	}
	return
}{{end}}{{end}}

`
