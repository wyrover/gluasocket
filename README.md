你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# LuaSocket for GopherLua

A native Go implementation of [LuaSocket](https://github.com/diegonehab/luasocket) for the [GopherLua](https://github.com/yuin/gopher-lua) VM.

## Using

### Loading Modules

```go
import (
	"github.com/nubix-io/gluasocket"
)

// Bring up a GopherLua VM
L := lua.NewState()
defer L.Close()

// Preload LuaSocket modules
gluasocket.Preload(L)
```

### Read lines from a socket

```go
script := `
  local client = require 'socket'.connect('127.0.0.1', 8000)
  local line1 = client:receive('*l')
  local line2 = client:receive('*l')
  client:close()
  return line1, line2`
L.DoString(script)
line1 := L.ToString(-2)
line2 := L.ToString(-1)
```

### Get system time

```go
L.DoString("return require 'socket'.gettime()")
gettimeValue := float64(L.ToNumber(-1))
```

## Testing

```bash
$ go test github.com/nubix-io/gluasocket...
ok  	github.com/nubix-io/gluasocket	0.045s
?   	github.com/nubix-io/gluasocket/ltn12	[no test files]
?   	github.com/nubix-io/gluasocket/mime	[no test files]
ok  	github.com/nubix-io/gluasocket/mimecore	0.040s
?   	github.com/nubix-io/gluasocket/socket	[no test files]
ok  	github.com/nubix-io/gluasocket/socketcore	0.269s
?   	github.com/nubix-io/gluasocket/socketexcept	[no test files]
?   	github.com/nubix-io/gluasocket/socketftp	[no test files]
?   	github.com/nubix-io/gluasocket/socketheaders	[no test files]
?   	github.com/nubix-io/gluasocket/sockethttp	[no test files]
?   	github.com/nubix-io/gluasocket/socketsmtp	[no test files]
?   	github.com/nubix-io/gluasocket/sockettp	[no test files]
?   	github.com/nubix-io/gluasocket/socketurl	[no test files]
```

Some original Lua-based LuaSocket unit tests are used and wrapped in Go unit test functions. Tests that perform `os.exit()` are modified to perform `error()` instead so that errors are made detectable.

## Design

### Divergence from LuaSocket source codes

#### 1. Unit test calls to `os.exit()` replaced with `error()`

This was necessary in order to detect and report errors from a test runner. Filed [LuaSocket Issue #227](https://github.com/diegonehab/luasocket/issues/227).

#### 2. Finalized Exceptions moved to pure-Lua

LuaSocket's exception handling is based on Diego's [Finalized Exceptions](http://lua-users.org/wiki/FinalizedExceptions) whitepaper.

After struggling to port the C-based `socket.newtry()` and `socket.protect()` functions to GopherLua, an easier path emerged when I discovered the pure-Lua implementations found in the unmerged [LuaSocket Pull Request #161](https://github.com/diegonehab/luasocket/pull/161), which introduces a new `socket.except` module, and makes tiny modifications to the `socket` module in order to use it. This served the purposes of this project perfectly.
