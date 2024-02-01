package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"
	"text/template"
)

const (
	FlagUsageOutput     = "(optional) output file path"
	FlagUsageName       = "the name of the proxy configuration type to be generated"
	FlagUsagePortNumber = "a port value for the port type to be generated"
)

func exit1(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}

type config struct {
	Output     string
	TypeName   string
	PortNumber int
}

func newConfig() (cfg config) {
	flag.StringVar(&cfg.Output, "output", ".", FlagUsageOutput)
	flag.StringVar(&cfg.TypeName, "type-name", "RegenerateMeWithName", FlagUsageName)
	flag.IntVar(&cfg.PortNumber, "port-number", -1, FlagUsagePortNumber)
	flag.Parse()
	return cfg
}

func main() {
	cfg := newConfig()

	tpl, err := template.New("proxy-conf-type-gen").Parse(text)
	if err != nil {
		exit1(err)
	}

	var b bytes.Buffer
	if err := tpl.Execute(&b, cfg); err != nil {
		exit1(err)
	}

	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, "", &b, parser.ParseComments)
	if err != nil {
		exit1(err)
	}

	fileName := fmt.Sprintf("type_%s.go", cfg.TypeName)
	fileName = strings.ToLower(fileName)
	fileName = path.Join(cfg.Output, fileName)

	file, err := os.Create(fileName)
	if err != nil {
		exit1(err)
	}

	written, err := io.Copy(file, &b)
	if err != nil {
		exit1(err)
	}

	slog.Info(fmt.Sprintf("%s: %d", fileName, written))
}

const text = `package v211

import "k8s.io/apimachinery/pkg/util/json"

type {{ .TypeName }} ProxyCfg

const DefaultPort{{ .TypeName }} = {{ .PortNumber }}

func (p *{{ .TypeName }}) UnmarshalJSON(text []byte) error {
	var result ProxyCfg
	if err := json.Unmarshal(text, &result); err != nil {
		return err
	}

	if result.Port == 0 {
		result.Port = DefaultPort{{ .TypeName }}
	}

	*p = {{ .TypeName }}(result)
	return nil
}

func (p *{{ .TypeName }}) MarshalJSON() ([]byte, error) {
	if p.Port == 0 {
		p.Port = DefaultPort{{ .TypeName }}
	}

	proxyCfg := ProxyCfg(*p)
	return json.Marshal(&proxyCfg)
}
`
