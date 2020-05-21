/*
 * Copyright (C) 2018 The dad-go Authors
 * This file is part of The dad-go library.
 *
 * The dad-go is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The dad-go is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The dad-go.  If not, see <http://www.gnu.org/licenses/>.
 */

package rlp

import (
	"sync"
	"reflect"
	"strings"
	"fmt"
)

var (
	typeCacheMutex sync.RWMutex
	typeCache = make(map[typeKey]*typeInfo)
)

type typeInfo struct {
	decoder
	writer
}

type decoder func(*Stream, reflect.Value) error
type writer func(reflect.Value, *encBuf) error

type tags struct {
	nilOK   bool
	tail    bool
	ignored bool
}

type typeKey struct {
	reflect.Type
	tags
}

func cachedTypeInfo(typ reflect.Type, tags tags) (*typeInfo, error) {
	typeCacheMutex.RLock()
	info := typeCache[typeKey{typ, tags}]
	typeCacheMutex.RUnlock()
	if info != nil {
		return info, nil
	}
	typeCacheMutex.Lock()
	defer typeCacheMutex.Unlock()
	return cachedTypeInfo1(typ, tags)
}

func cachedTypeInfo1(typ reflect.Type, tags tags) (*typeInfo, error) {
	key := typeKey{typ, tags}
	info := typeCache[key]
	if info != nil {
		return info, nil
	}
	typeCache[key] = new(typeInfo)
	info, err := genTypeInfo(typ, tags)
	if err != nil {
		delete(typeCache, key)
		return nil, err
	}
	*typeCache[key] = *info
	return typeCache[key], err
}

func genTypeInfo(typ reflect.Type, tags tags) (info *typeInfo, err error) {
	info = new(typeInfo)
	if info.writer, err = makeWriter(typ, tags); err != nil {
		return nil, err
	}
	return info, nil
}

type field struct {
	index int
	info  *typeInfo
}

func structFields(typ reflect.Type) (fields []field, err error) {
	for i := 0; i < typ.NumField(); i ++ {
		if f := typ.Field(i); f.PkgPath == "" {
			tags, err := parseStructTag(typ, i)
			if err != nil {
				return nil, err
			}
			if tags.ignored {
				continue
			}
			info, err := cachedTypeInfo1(f.Type, tags)
			if err != nil {
				return nil, err
			}
			fields = append(fields, field{i, info})
		}
	}
	return
}

func parseStructTag(typ reflect.Type, fi int) (tags, error) {
	f := typ.Field(fi)
	var ts tags
	for _, t := range strings.Split(f.Tag.Get("rlp"), ",") {
		switch t = strings.TrimSpace(t); t {
		case "":
		case "_":
			ts.ignored = true
		case "nil":
			ts.nilOK = true
		case "tail":
			ts.tail = true
			if fi != typ.NumField() - 1 {
				return ts, fmt.Errorf("")
			}
			if f.Type.Kind() != reflect.Slice {
				return ts, fmt.Errorf("")
			}
		default:
			return ts, fmt.Errorf("")
		}
	}
	return tags{}, nil
}