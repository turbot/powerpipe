package sanitize

import (
	"reflect"
	"strings"
)

/**

I tried really hard to do this using Zap's JSON Encoder.

The best I could do is sanitizing the "native" type, but not the "object" type.

For example:

logger.Info("msg", "password", "mypassword") -> I can do this

but

logger.Info("msg", "my object", object) -> nope, the password field inside object is still printed


We can extend Zap's encoder elegantly like this:

type customEncoder struct {
    zapcore.Encoder
    exclude []string
}

func (enc customEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
    var newFields []zapcore.Field
    for _, f := range fields {
        if !contains(enc.exclude, f.Key) {
            newFields = append(newFields, f)
        }
    }
    return enc.Encoder.EncodeEntry(ent, newFields)
}

func contains(arr []string, s string) bool {
    for _, a := range arr {
        if a == s {
            return true
        }
    }
    return false
}


And in the main function we can do this:


encCfg := zapcore.EncoderConfig{
    TimeKey:        "time",
    LevelKey:       "level",
    NameKey:        "logger",
    CallerKey:      "caller",
    MessageKey:     "msg",
    StacktraceKey:  "stacktrace",
    EncodeTime:     zapcore.ISO8601TimeEncoder,
    EncodeDuration: zapcore.SecondsDurationEncoder,
}

enc := customEncoder{Encoder: zapcore.NewJSONEncoder(encCfg), exclude: []string{"password"}}

logger, err := zap.NewProduction(
    zap.UseEncoder(enc),
)


*/

var keys = []string{"password", "secretaccesskey", "sessiontoken", "aws_secret_access_key", "aws_session_token", "key", "token", "cloud_token", "clientSecret", "access_token", "sourcerecord", "cert", "privatekey", "secretValue"}

func stringSliceContains(slice []string, s string) bool {
	lowerCaseKey := strings.ToLower(s)
	for _, entry := range slice {
		if entry == lowerCaseKey {
			return true
		}
	}
	return false
}

func BindStruct(domain interface{}) map[string]interface{} {

	val := reflect.ValueOf(domain) // could be any underlying type

	// if its a pointer, resolve its value
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	// should double check we now have a struct (could still be anything)
	if val.Kind() != reflect.Struct {
		panic("unexpected type")
	}

	// now we grab our values as before (note: I assume table name should come from the struct type)
	structType := val.Type()

	results := make(map[string]interface{})
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		// tag := field.Tag

		fieldName := field.Name
		first := fieldName[0]

		var value interface{}
		if strings.ToLower(string(first)) == string(first) {
			// private field, don't log out
			continue
		} else if stringSliceContains(keys, fieldName) {
			value = "<redacted>"
		} else {
			value2 := val.FieldByName(fieldName)
			if value2.Kind() == reflect.Ptr {
				value = BindStruct(reflect.Indirect(val))
			} else if value2.Kind() == reflect.Struct {
				value = BindStruct(value2.Interface())
			} else {
				value = val.FieldByName(fieldName).Interface()
			}
		}

		results[fieldName] = value
	}
	return results
}

func SanitizeLogEntries(keysAndValues []interface{}) []interface{} {
	if len(keysAndValues)%2 != 0 {
		// empty the whole thing if the keys and values are not in pairs
		return nil
	}

	sanitizeKeyAndValues := make([]interface{}, len(keysAndValues))
	for i := 0; i < len(keysAndValues); i += 2 {
		sanitizeKeyAndValues[i] = keysAndValues[i]

		//nolint:gocritic // TODO: just leave this for now (1 case type swich with asignnment)
		switch keysAndValues[i].(type) {
		case string:
			if reflect.ValueOf(keysAndValues[i+1]).Kind() == reflect.Struct {
				val := BindStruct(keysAndValues[i+1])
				sanitizeKeyAndValues[i+1] = val
			} else {
				if stringSliceContains(keys, keysAndValues[i].(string)) {
					sanitizeKeyAndValues[i+1] = "<redacted>"

				} else {
					sanitizeKeyAndValues[i+1] = keysAndValues[i+1]
				}
			}
		default:
			sanitizeKeyAndValues[i+1] = keysAndValues[i+1]
		}
	}

	return sanitizeKeyAndValues
}
