package utils

type ctxKeyRequestID struct{}

func CtxKeyRequestIDGet() ctxKeyRequestID {
	return ctxKeyRequestID{}
}

type ctxKeyMethod struct{}

func CtxKeyMethodGet() ctxKeyMethod {
	return ctxKeyMethod{}
}

type ctxKeyURL struct{}

func CtxKeyURLGet() ctxKeyURL {
	return ctxKeyURL{}
}
