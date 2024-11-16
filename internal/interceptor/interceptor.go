package interceptor

type Interceptor interface {
	SetNext(Interceptor)
	Handle()
}

type BaseInterceptor struct {
	next Interceptor
}

func (interceptor *BaseInterceptor) SetNext(next Interceptor) {
	interceptor.next = next
}
