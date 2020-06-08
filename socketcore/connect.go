package gluasocket_socketcore

import (
	"bufio"
	"fmt"
	"net"
	"golang.org/x/net/proxy"
	"github.com/yuin/gopher-lua"
)

func connectFn(L *lua.LState) int {
	hostname := L.ToString(1)
	port := L.ToInt(2)
	socks5_host := L.OptString(3, "")
	socks5_usr := L.OptString(4, "")
	socks5_pwd := L.OptString(5, "")

	var conn net.Conn
	var err error

	if socks5_host != "" {
		socks5_auth := &proxy.Auth{}
		if socks5_usr != ""{
			socks5_auth = &proxy.Auth{
				User:     socks5_usr,
				Password: socks5_pwd,
			}
		} else {
			socks5_auth = nil
		}
		sk5, err := proxy.SOCKS5("tcp", socks5_host, socks5_auth, proxy.Direct)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		conn, err = sk5.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))

		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	} else {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))

		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	}

	reader := bufio.NewReader(conn)
	client := &Client{Conn: conn, Reader: reader}
	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable(CLIENT_TYPENAME))
	L.Push(ud)
	return 1
}
