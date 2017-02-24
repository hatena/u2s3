package io

type BufferedReader interface {
	Readln() ([]byte, error)
	Close()
}
