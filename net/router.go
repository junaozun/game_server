package net

type group struct {
	prefix     string
	handlerMap map[string]HandlerFunc
}

type HandlerFunc func()

type router struct {
}