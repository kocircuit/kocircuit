package inject

import (
	"bytes"
	"text/template"
)

type InjectionProgram struct {
	GatePkgPath   string   `ko:"name=gatePkgPath"`
	GateTypeName  string   `ko:"name=gateTypeName"`
	InjectPkgPath []string `ko:"name=injectPkgPath"`
}

func (prog *InjectionProgram) Render() string {
	var w bytes.Buffer
	if err := injectionTmpl.Execute(&w, prog); err != nil {
		panic(err)
	}
	return w.String()
}

var injectionTmpl = template.Must(template.New("").Parse(injectionSrc))

var injectionSrc = `
package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"io/ioutil"
	"log"

	"github.com/kocircuit/kocircuit/lang/go/runtime"

	gatePkg {{ printf "%q" .GatePkgPath }}

{{range .InjectPkgPath -}}
	_ {{ printf "%q" . }}
{{- end -}}
)

var flagArgs = flag.String("args", "", "Path to input file with encoded arguments structure")
var flagReturn = flag.String("return", "", "Path to output file with encoded return value")

func main() {
	flag.Parse()
	argBytes, err := ioutil.ReadFile(*flagArgs)
	if err != nil {
		log.Fatalf("reading arguments from file %q (%v)", *flagArgs, err)
	}
	dec := gob.NewDecoder(bytes.NewBuffer(argBytes))
	var args = new(gatePkg.{{.GateTypeName}})
	if err = dec.Decode(args); err != nil {
		log.Fatalf("decoding arguments (%v)", err)
	}
	result := args.Play(runtime.NewContext())
	var resultBuf bytes.Buffer
	enc := gob.NewEncoder(&resultBuf)
	if err = enc.Encode(result); err != nil {
		log.Fatalf("encoding returned value (%v)", err)
	}
	if err = ioutil.WriteFile(*flagReturn, resultBuf.Bytes(), 0644); err != nil {
		log.Fatalf("writing returned value to file %q (%v)", *flagReturn, err)
	}
}`
