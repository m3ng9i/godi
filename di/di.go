// Package di calls 'go list' command to gather information about go package dependency.
package di

import "strings"
import "strconv"
import "encoding/json"
import "os/exec"
import "errors"
import "github.com/m3ng9i/go-utils/cmd"


// max length of "go list" command line parameters
var maxChars = 8000


// PkgInfo stores information of a package
type PkgInfo struct {
    Name        string  // name of the package
    ImportPath  string  // import path of the package
    Directly    bool    // if the package is directly imported
    Builtin     bool    // if the package is builtin package
    SubPkg      bool    // if the package is a sub package
}


// used in channel
type pkgInfoC struct {
    PkgInfo
    err error
}


type pkg struct {
    ImportPath string
    Directly bool
}


type deps struct {
    Imports []string
    Deps    []string
}


// call "go list -json 'pkg name'", get all dependency package's importpath
func di(importPath string) (p []pkg, err error) {
    sout, serr, e := cmd.Run("go", "list", "-json", importPath)
    if len(serr) > 0 {
        err = errors.New(string(serr))
        return
    }
    if e != nil {
        err = e
        return
    }

    d := new(deps)
    err = json.Unmarshal(sout, d)
    if err != nil {
        return
    }

    p = make([]pkg, len(d.Deps))

    Loop:
    for i, j := range d.Deps {
        p[i].ImportPath = j
        for _, k := range d.Imports {
            if k == j {
                p[i].Directly = true
                continue Loop
            }
        }
    }

    return
}


func pkgInfo(importPath string, p []pkg) (info []PkgInfo, err error) {

    // use "go list" command to list one or more packages' information
    // command is like: go list -f "{{.Standard}}: {{.Name}}" package1 package2

    args := []string{"list", "-f", "{{.Name}}:{{.Standard}}"}

    for _, i := range p {
        args = append(args, i.ImportPath)
    }

    sout, serr, e := cmd.Run("go", args...)
    if len(serr) > 0 {
        err = errors.New(string(serr))
        return
    }
    if e != nil {
        err = e
        return
    }

    for n, line := range strings.Split(string(sout), "\n") {
        if n >= len(p) {
            break
        }

        result := strings.SplitN(line, ":", 2)

        i := PkgInfo {
            Name:       result[0],
            ImportPath: p[n].ImportPath,
            Directly:   p[n].Directly,
            SubPkg:     strings.HasPrefix(p[n].ImportPath, importPath),
        }

        i.Builtin, err = strconv.ParseBool(result[1])
        if err != nil {
            return
        }

        info = append(info, i)
    }

    return
}


// SetMaxChars set value of maxChars, if n is < 200, the function return false
func SetMaxChars(n int) (success bool) {
    if n < 200 {
        return false
    }
    maxChars = n
    return true
}


/*List call "go list" command to gather package dependency information

Because the command line has a string length limitation, if parameters of
"go list" command is too long, this function will split the parameters
and call "go list" more times.

You can use SetMaxChars() to set the max length of command line parameters.

Parameters:
    all:        true:   get info of all (recursively) imported packages
                false:  get info of packages which are directly imported
    builtin:    whether or not including builtin packages
    subpkg:     whether or not including sub-packages
*/
func List(importPath string, all, builtin, subpkg bool) (info []PkgInfo, err error) {
    var p []pkg
    p, err = di(importPath)
    if err != nil {
        return
    }

    if len(p) == 0 {
        return
    }

    var pinfos []PkgInfo
    count := 0
    pos := 0

    for i, j := range p {
        count += len(j.ImportPath)

        // if length of the command string is too long,
        // or if the loop is iterating for the last time,
        // call "go list" command to gather the package information.
        if count > maxChars || i == len(p) - 1 {
            pinfo, e := pkgInfo(importPath, p[pos : i + 1])
            if e != nil {
                err = e
                return
            }
            pinfos = append(pinfos, pinfo...)
            pos = i + 1
            count = 0
        }
    }

    for _, i := range pinfos {
        if (!all && !i.Directly) || (!builtin && i.Builtin) || (!subpkg && i.SubPkg) {
            continue
        }
        info = append(info, i)
    }

    return
}


// ListCD get package's import path from the current directory
func ListCD() (p string, err error) {
    sout, serr, e := cmd.Run("go", "list")
    if len(serr) > 0 {
        err = errors.New(string(serr))
        return
    }
    if e != nil {
        err = e
        return
    }
    p = strings.TrimSpace(string(sout))
    return
}


// GoExist check if go command is exist in PATH environment variable.
func GoExist() bool {
    _, err := exec.LookPath("go")
    return err == nil
}

