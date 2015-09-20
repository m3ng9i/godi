// Godi: go package dependency information
package main

import "fmt"
import "os"
import "flag"
import "encoding/json"
import "github.com/m3ng9i/godi/di"


var usage = `Godi: go package dependency information

Usage:

    godi [-a] [-b] [-s] [-v | -j] [package]
    godi -e
    godi -h

Options:

    -a          Get infomation of all (recursively) imported packages.
                If not set this option, only get info of packages which are
                directly imported.
    -b=bool     If include builtin packages in the result. Default is true.
    -s=bool     If include sub-packages in the result. Default is true.
    -v          Show verbose result.
    -j          Show JSON format result.

    -e          Show examples.
    -h, -help   Show help message.

Notice:

    1) "go" command must be installed on your computer.
    2) The "package" option must be placed in the end of the command line.
    3) The package which you want to check must be installed on the computer.
    4) If not provide the "package" option, godi will try to gather information from the current directory.

Author:
    m3ng9i <https://github.com/m3ng9i>
`

var examples = `Examples:

1) See what the builtin package "log" is directly depends on:

    godi log

2) See what "log" and all it's dependent packages are depend on:

    godi -a log

3) See what "net/http" is directly depends on, remove sub-package:

    godi -s=false net/http

   You will find "net/http/internal" not shows in the result.

4) See what "github.com/m3ng9i/feedreader" is directly depends on and remove go's builtin package in the result:

    godi -b=false github.com/m3ng9i/feedreader

5) See all the dependency information of "github.com/m3ng9i/go-utils/cmd" in table format:

    godi -v -a github.com/m3ng9i/go-utils/cmd | column -t

6) See what "bufio" is directly depends on and output a JSON format result:

    godi.exe -j bufio
`


type cmdOption struct {
    all         bool
    builtin     bool
    subPkg      bool
    verbose     bool
    json        bool
    example     bool
    help        bool
    pkg         string
}


func outputNormal(info []di.PkgInfo) {
    for _, i := range info {
        fmt.Println(i.ImportPath)
    }
}


func outputVerbose(info []di.PkgInfo) {
    fmt.Printf("Name\tDirectly\tBuiltin\tSubPkg\tImportPath\n")
    for _, i := range info {
        fmt.Printf("%s\t%v\t%v\t%v\t%s\n", i.Name, i.Directly, i.Builtin, i.SubPkg, i.ImportPath)
    }
}


func outputJSON(info []di.PkgInfo) {
    b, err := json.MarshalIndent(info, "", " ")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    fmt.Println(string(b))
}


func main() {

    var o cmdOption

    flag.BoolVar(&o.all,       "a",       false,  "Get info of all imported packages")
    flag.BoolVar(&o.builtin,   "b",       true,   "Set if get information of builtin packages")
    flag.BoolVar(&o.subPkg,    "s",       true,   "Set if get information of sub-packages")
    flag.BoolVar(&o.verbose,   "v",       false,  "Show more information of packages")
    flag.BoolVar(&o.json,      "j",       false,  "Show JSON format result")
    flag.BoolVar(&o.example,   "e",       false,  "Show examples")
    flag.BoolVar(&o.help,      "h",       false,  "Show help message")
    flag.BoolVar(&o.help,      "help",    false,  "Show help message")

    flag.Usage = func() {
        fmt.Print(usage)
    }

    flag.Parse()

    if o.help {
        flag.Usage()
        os.Exit(0)
    }

    if o.example {
        fmt.Print(examples)
        os.Exit(0)
    }

    if o.verbose && o.json {
        fmt.Fprintln(os.Stderr, "Option -v/-verbose and -j/-json cannot be used together")
        os.Exit(1)
    }

    if !di.GoExist() {
        fmt.Fprintln(os.Stderr, `Command "go" not found`)
        os.Exit(1)
    }

    o.pkg  = flag.Arg(0)
    if o.pkg == "" {
        var err error
        o.pkg, err = di.ListCD()
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
    }

    pkgInfo, err := di.List(o.pkg, o.all, o.builtin, o.subPkg)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    if len(pkgInfo) == 0 {
        fmt.Fprintln(os.Stderr, "no information")
        os.Exit(1)
    }

    if o.verbose {
        outputVerbose(pkgInfo)
    } else if o.json {
        outputJSON(pkgInfo)
    } else {
        outputNormal(pkgInfo)
    }

}

