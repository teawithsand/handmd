package scripting

import "net"

// Returns free port to use by creating and closing listener on 0 port.
//
// Note: this creates race condition, which malicious software may use, since there
// is lag between picking port and allocating it again.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
