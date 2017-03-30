package content

type BufferedReader interface {
	Readln() ([]byte, error)
	Close()
}
