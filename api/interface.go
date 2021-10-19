package api

type dataLogger interface {
	LogData() map[string]interface{}
}

type coder interface {
	Code() int
}

type responser interface {
	Response() string
}
