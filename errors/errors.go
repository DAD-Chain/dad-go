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
