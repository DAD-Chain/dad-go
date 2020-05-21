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

package errors

import (
	"errors"
)

const callStackDepth = 10

type DetailError interface {
	error
	ErrCoder
	CallStacker
	GetRoot() error
}

func NewErr(errmsg string) error {
	return errors.New(errmsg)
}

func NewDetailErr(err error, errcode ErrCode, errmsg string) DetailError {
	if err == nil {
		return nil
	}

	dad-goerr, ok := err.(dad-goError)
	if !ok {
		dad-goerr.root = err
		dad-goerr.errmsg = err.Error()
		dad-goerr.callstack = getCallStack(0, callStackDepth)
		dad-goerr.code = errcode

	}
	if errmsg != "" {
		dad-goerr.errmsg = errmsg + ": " + dad-goerr.errmsg
	}

	return dad-goerr
}

func RootErr(err error) error {
	if err, ok := err.(DetailError); ok {
		return err.GetRoot()
	}
	return err
}
