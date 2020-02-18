package errors

type dad-goError struct {
	errmsg    string
	callstack *CallStack
	root      error
	code      ErrCode
}

func (e dad-goError) Error() string {
	return e.errmsg
}

func (e dad-goError) GetErrCode() ErrCode {
	return e.code
}

func (e dad-goError) GetRoot() error {
	return e.root
}

func (e dad-goError) GetCallStack() *CallStack {
	return e.callstack
}
