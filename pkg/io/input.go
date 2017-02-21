package io

type BufferedReader interface {
	Readln() (string, error)
	Close()
}
