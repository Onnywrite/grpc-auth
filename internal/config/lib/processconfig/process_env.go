package processconfig

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func PrioritizeEnvs(cfgStructPtr interface{}) {
	v := reflect.ValueOf(cfgStructPtr)
	t := reflect.TypeOf(cfgStructPtr)

	for i := range v.NumField() {
		f := v.Field(i)

		if f.Kind() == reflect.Struct {
			PrioritizeEnvs(&f)
		}

		value, ok := t.Field(i).Tag.Lookup("env")
		env, ok2 := os.LookupEnv(value)

		if !(ok && ok2 && f.CanSet()) {
			continue
		}

		if parseInto(f, f.Type(), env) != nil {
			continue
		}
	}
}

func parseEnvSlice(envType reflect.Type, env string) (*reflect.Value, error) {
	envSlice := strings.Split(env, ":")

	slice := reflect.MakeSlice(envType, 0, len(envSlice))

	for _, e := range envSlice {
		var eValue reflect.Value
		if err := parseInto(eValue, envType, e); err != nil {
			return nil, err
		}
		reflect.Append(slice, eValue)
	}

	return &slice, nil
}

func parseInto(f reflect.Value, fType reflect.Type, env string) error {
	switch f.Kind() {

	case reflect.String:
		f.SetString(env)

	case reflect.Bool:
		b, err := strconv.ParseBool(env)
		if err != nil {
			return fmt.Errorf("could not parse env")
		}
		f.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		number, err := strconv.ParseInt(env, 0, fType.Bits())
		if err != nil {
			return fmt.Errorf("could not parse env")
		}
		f.SetInt(number)

	case reflect.Int64:
		if fType == reflect.TypeOf(time.Duration(0)) {
			// try to parse time
			d, err := time.ParseDuration(env)
			if err != nil {
				return fmt.Errorf("could not parse env")
			}
			f.SetInt(int64(d))
		} else {
			// parse regular integer
			number, err := strconv.ParseInt(env, 0, fType.Bits())
			if err != nil {
				return fmt.Errorf("could not parse env")
			}
			f.SetInt(number)
		}

	// parse unsigned integer env
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		number, err := strconv.ParseUint(env, 0, fType.Bits())
		if err != nil {
			return fmt.Errorf("could not parse env")
		}
		f.SetUint(number)

	// parse floating point env
	case reflect.Float32, reflect.Float64:
		number, err := strconv.ParseFloat(env, fType.Bits())
		if err != nil {
			return fmt.Errorf("could not parse env")
		}
		f.SetFloat(number)

	// parse sliced env
	case reflect.Slice:
		sliceValue, err := parseEnvSlice(fType, env)
		if err != nil {
			return fmt.Errorf("could not parse env")
		}
		f.Set(*sliceValue)
	}

	return nil
}
