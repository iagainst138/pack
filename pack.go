package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "regexp"
    "strings"
    "text/template"
    )

const map_file_template = `package {{.PackageName}}

var {{.VarName}} = map[string]string {
{{range .Vars}}
"{{.Key}}":` +  "`{{.Data}}`," + `
{{end}}
}
`

type Var struct {
    Key     string
    Data    string
}

type FileDetails struct {
    PackageName string
    VarName     string
    Vars        []Var
}

func gen_file(include_dir, out_file, package_name, var_name string, regex *regexp.Regexp) error {
    vars := []Var{}

    if file_info,err := ioutil.ReadDir(include_dir); err == nil {
        for _,f := range file_info {
            if f.IsDir() == false {
                if regex.MatchString(f.Name()){
                    p := path.Join(include_dir, f.Name())
                    if b,err := ioutil.ReadFile(p); err == nil {
                        fmt.Println("[" + out_file + "] adding", f.Name())
                        vars = append(vars, Var{f.Name(), string(b)})
                    } else { return err }
                }
            } else { fmt.Println("Directories are not supported -", f.Name()) }
        }
    } else {
        return err
    }

    t := template.Must(template.New("source").Parse(map_file_template))
    out,err := os.Create(out_file)
    if err != nil {
        return err
    }
    return t.Execute(out, FileDetails{package_name, var_name, vars})
}

func main() {
    var out_file string
    var package_name string
    var regex string
    var resource_dir string
    var var_name string

    flag.StringVar(&resource_dir, "i", "", "directory to load from")
    flag.StringVar(&out_file, "o", "", "go file to create")
    flag.StringVar(&package_name, "p", "resources", "package name to use")
    flag.StringVar(&regex, "m", ".*", "regexp to match")
    flag.StringVar(&var_name, "var", "R", "variable name")
    flag.Parse()

    if out_file == "" { panic("no output file specified.") }
    if resource_dir == "" { panic("no resource directory specified.") }

    d := strings.Split(out_file, "/")[0]
    if _,e := os.Stat(d); os.IsNotExist(e) {
        os.Mkdir(d, 0755)
    }

    if err := gen_file(resource_dir, out_file, package_name, var_name, regexp.MustCompile(regex)); err != nil {
        panic(err)
    }
}
