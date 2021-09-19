package utils

import (
	"github.com/emirpasic/gods/sets/hashset"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

func BsonMEqual(v1, v2 bson.M) (ok bool) {
	//ok = reflect.DeepEqual(v1, v2)
	ok = bsonMEqual(v1, v2)
	return ok
}

func bsonMEqual(v1, v2 bson.M) (ok bool) {
	// all keys
	allKeys := hashset.New()
	for key := range v1 {
		allKeys.Add(key)
	}
	for key := range v2 {
		allKeys.Add(key)
	}

	for _, keyRes := range allKeys.Values() {
		key := keyRes.(string)
		v1Value, ok := v1[key]
		if !ok {
			return false
		}
		v2Value, ok := v2[key]
		if !ok {
			return false
		}

		mode := 0

		var v1ValueBsonM bson.M
		var v1ValueBsonA bson.A
		switch v1Value.(type) {
		case bson.M:
			mode = 1
			v1ValueBsonM = v1Value.(bson.M)
		case bson.A:
			mode = 2
			v1ValueBsonA = v1Value.(bson.A)
		}

		var v2ValueBsonM bson.M
		var v2ValueBsonA bson.A
		switch v2Value.(type) {
		case bson.M:
			if mode != 1 {
				return false
			}
			v2ValueBsonM = v2Value.(bson.M)
		case bson.A:
			if mode != 2 {
				return false
			}
			v2ValueBsonA = v2Value.(bson.A)
		}

		switch mode {
		case 0:
			if v1Value != v2Value {
				return false
			}
		case 1:
			if !bsonMEqual(v1ValueBsonM, v2ValueBsonM) {
				return false
			}
		case 2:
			if !reflect.DeepEqual(v1ValueBsonA, v2ValueBsonA) {
				return false
			}
		default:
			// not reachable
			return false
		}
	}

	return true
}
