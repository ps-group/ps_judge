package restapi

// MethodHandler - handler which takes request and returns serializable interace or an error
type MethodHandler func(context interface{}, request Request) Response
