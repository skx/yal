// reflected.go - Call native golang code via reflection.


package eval

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/skx/yal/primitive"
)

// Holder for mappings between native and lisp
var Reflected map[string]interface{}

// init sets up the bindings between golang functions we've made
// available and their names.
func init() {
	Reflected = make(map[string]interface{})

	Reflected["math.Abs"] = math.Abs
	Reflected["math.Intn"] = rand.Intn
	Reflected["math.Sin"] = math.Sin

	Reflected["os.Getenv"] = os.Getenv
	Reflected["os.Stat"] = os.Stat

	Reflected["time.Now"] = time.Now

	Reflected["path.Base"] = filepath.Base
	Reflected["path.Glob"] = filepath.Glob

	Reflected["strings.Contains"] = strings.Contains
	Reflected["strings.Cut"] = strings.Cut
	Reflected["strings.Split"] = strings.Split
	Reflected["strings.ToLower"] = strings.ToLower
	Reflected["strings.ToUpper"] = strings.ToUpper
}



// Call a function, dynamically.
//
// This is a bit horrid, but it tries to map from the Yal types to the
// appropriate types which the golang functions will require.
//
// Arrays are not handled in the general case, though string arrays
// will work.
func Call(funcName string, params interface{}) (result interface{}, err error) {

	f := reflect.ValueOf(Reflected[funcName])

	p := params.([]primitive.Primitive)

	// TODO: Cope with varargs, etc.
	if len(p) != f.Type().NumIn() {
		err = errors.New(fmt.Sprintf("Function %s expects %d arguments %d supplied", funcName, f.Type().NumIn(), len(p)))
		return
	}

	// create holder for values
	in := make([]reflect.Value, len(p))

	for i := 0; i < f.Type().NumIn(); i++ {

		// This should be better
		ty := f.Type().In(i)

		// This is also horribly wrong
		switch fmt.Sprintf("%s", ty) {

		case "string":
			s := string(p[i].(primitive.String))
			in[i] = reflect.ValueOf(s)

		case "[]string":
			tmp := []string{}

			lst, ok := p[i].(primitive.List)
			if !ok {
				return primitive.Error(fmt.Sprintf("argument %d should have been []string, but we were given %T", i, p[i])), nil
			}
			for _, x := range lst {
				tmp = append(tmp, x.ToString())
			}
			in[i] = reflect.ValueOf(tmp)

		case "int", "int32", "int64", "uint8", "uint16", "uint32", "uint64":
			n := int(p[i].(primitive.Number))
			in[i] = reflect.ValueOf(n)

		case "float32", "float64":
			n := float64(p[i].(primitive.Number))
			in[i] = reflect.ValueOf(n)

		default:
			return primitive.Error(fmt.Sprintf("unknown param-type %s", ty)), nil
		}
	}

	// Call the function, and get the result
	res := f.Call(in)

	// No results?  Return nil
	if len(res) == 0 {
		return primitive.Nil{}, nil
	}

	// One result?  Return it
	if len(res) == 1 {
		return (goToLisp(res[0].Interface())), nil
	}

	// Otherwise return all values as a list
	var r primitive.List

	for _, e := range res {
		tmp := e.Interface()
		r = append(r, goToLisp(tmp))
	}
	return r, nil
}

// goToLisp will attempt to turn a native golang type to a suitable
// YAL type.  So an integer will become a Number, etc.
//
// This is non-exhaustive, and does not at all work with arrays of mixed
// types, or structures.
func goToLisp(result any) primitive.Primitive {

	switch result.(type) {

	// number / array of numbers
	case int:
		return primitive.Number(result.(int))

	case int16:
		return primitive.Number(result.(int16))

	case int32:
		return primitive.Number(result.(int32))

	case int64:
		return primitive.Number(result.(int64))

	case uint8:
		return primitive.Number(result.(uint8))

	case uint16:
		return primitive.Number(result.(uint16))

	case uint32:
		return primitive.Number(result.(uint32))

	case uint64:
		return primitive.Number(result.(uint64))

	case []uint8:
		r := primitive.List{}
		for _, x := range result.([]uint8) {
			r = append(r, primitive.Number(x))
		}
		return r

	case []uint16:
		r := primitive.List{}
		for _, x := range result.([]uint16) {
			r = append(r, primitive.Number(x))
		}
		return r

	case []uint32:
		r := primitive.List{}
		for _, x := range result.([]uint32) {
			r = append(r, primitive.Number(x))
		}
		return r

	case []uint64:
		r := primitive.List{}
		for _, x := range result.([]uint64) {
			r = append(r, primitive.Number(x))
		}
		return r

	case []int:
		r := primitive.List{}
		for _, x := range result.([]int) {
			r = append(r, primitive.Number(x))
		}
		return r

	case  []int32:
		r := primitive.List{}
		for _, x := range result.([]int32) {
			r = append(r, primitive.Number(x))
		}
		return r

	case []int64:
		r := primitive.List{}
		for _, x := range result.([]int64) {
			r = append(r, primitive.Number(x))
		}
		return r

	// float / array of floats
	case float32:
		return primitive.Number(result.(float32))

	case float64:
		return primitive.Number(result.(float64))

	case []float32:
		r := primitive.List{}
		for _, x := range result.([]float32) {
			r = append(r, primitive.Number(x))
		}
		return r

	case []float64:
		r := primitive.List{}
		for _, x := range result.([]float64) {
			r = append(r, primitive.Number(x))
		}
		return r

	// string / array of strings
	case string:
		return primitive.String(result.(string))

	case []string:
		r := primitive.List{}
		for _, x := range result.([]string) {
			r = append(r, primitive.String(x))
		}
		return r

	// misc
	case bool:
		return primitive.Bool(result.(bool))

	case error:
		return primitive.Error(result.(error).Error())

	case nil:
		return primitive.Nil{}

	case time.Time:
		time, ok := result.(time.Time)
		if ok {
			return primitive.Number(time.Unix())
		}
		return primitive.Nil{}

	default:
		return primitive.Error(fmt.Sprintf("unknown return type %T", result))
	}
}
