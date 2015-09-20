Godi: go package dependency information
=======================================

Godi is a command line tool for gathering a go package's dependency information.

## Requirement

Godi calls "go list" to gather information, so you need "go" command installed on your computer.

## Install godi

```bash
go get github.com/m3ng9i/godi
```

This will install godi binary in $GOPATH/bin.

## Command line options

```
-a          Get infomation of all (recursively) imported packages.
            If not set this option, only get info of packages which are
            directly imported.
-b=bool     If include builtin packages in the result. Default is true.
-s=bool     If include sub-packages in the result. Default is true.
-v          Show verbose result.
-j          Show JSON format result.

-e          Show examples.
-h, -help   Show help message.
```

## How to use

Example 1: see what the builtin package "log" is directly depends on:

```bash
godi log
```

The result is:

```
fmt
io
os
runtime
sync
time
```

Example 2: see what "log" and all it's dependent packages are depend on:

```bash
godi -a log
```

The result is:

```
errors
fmt
internal/syscall/windows
internal/syscall/windows/registry
io
math
os
reflect
runtime
strconv
sync
sync/atomic
syscall
time
unicode/utf16
unicode/utf8
unsafe
```

Example 3: see what "net/http" is directly depends on, remove sub-package:

```bash
godi -s=false net/http
```

The result is:

```
bufio
bytes
compress/gzip
crypto/tls
encoding/base64
encoding/binary
errors
fmt
io
io/ioutil
log
mime
mime/multipart
net
net/textproto
net/url
os
path
path/filepath
runtime
sort
strconv
strings
sync
sync/atomic
time
unicode/utf8
```

Because the use of `-s=false` option, "net/http/internal" not shows in the result.

Example 4: see what "github.com/m3ng9i/feedreader" is directly depends on and remove go's builtin package in the result:

```bash
godi -b=false github.com/m3ng9i/feedreader
```

Because the use of `-b=false` option, the result will not including any builting packages:

```
github.com/m3ng9i/go-utils/encoding
github.com/m3ng9i/go-utils/html
github.com/m3ng9i/go-utils/http
github.com/m3ng9i/go-utils/set
github.com/m3ng9i/go-utils/xml
```

Example 5: see all the dependency information of "github.com/m3ng9i/go-utils/cmd" in table format:

```bash
godi -v -a github.com/m3ng9i/go-utils/cmd | column -t
```

The result is:

```
Name      Directly  Builtin  SubPkg  ImportPath
bytes     true      true     false   bytes
errors    false     true     false   errors
io        true      true     false   io
math      false     true     false   math
os        true      true     false   os
exec      true      true     false   os/exec
filepath  false     true     false   path/filepath
runtime   false     true     false   runtime
sort      false     true     false   sort
strconv   false     true     false   strconv
strings   false     true     false   strings
sync      false     true     false   sync
atomic    false     true     false   sync/atomic
syscall   false     true     false   syscall
time      false     true     false   time
unicode   false     true     false   unicode
utf8      false     true     false   unicode/utf8
unsafe    false     true     false   unsafe
```

Example 6: see what "bufio" is directly depends on and output a JSON format result:

```bash
godi -j bufio
```

The result is:

```
[
 {
  "Name": "bytes",
  "ImportPath": "bytes",
  "Directly": true,
  "Builtin": true,
  "SubPkg": false
 },
 {
  "Name": "errors",
  "ImportPath": "errors",
  "Directly": true,
  "Builtin": true,
  "SubPkg": false
 },
 {
  "Name": "io",
  "ImportPath": "io",
  "Directly": true,
  "Builtin": true,
  "SubPkg": false
 },
 {
  "Name": "utf8",
  "ImportPath": "unicode/utf8",
  "Directly": true,
  "Builtin": true,
  "SubPkg": false
 }
]
```

Besides above, if you want to check the dependency information of the package in current working directory, just type "godi" and press return.

## Author

mengqi (aka m3ng9i) <https://github.com/m3ng9i>

