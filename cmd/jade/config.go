package main

import (
	"bytes"
	"io"
	"log"
	"strings"
	"text/template"

	"github.com/Joker/jade"
)

const (
	file_bgn = `// Code generated by "jade.go"; DO NOT EDIT.

package {{.Package}}

import (
{{- range .Import}}
	{{.}}
{{- end}}
)

{{- range .Def}}
	{{.}}
{{- end}}

{{.Func}} {
	{{.Before}}
`
	file_end = `
	{{.After}}
}
`
)

var golang = jade.ReplaseTokens{
	GolangMode: true,
	TagBgn:     "\nbuffer.WriteString(`<%s%s>`)",
	TagEnd:     "\nbuffer.WriteString(`</%s>`)",
	TagVoid:    "\nbuffer.WriteString(`<%s%s/>`)",
	TagArgEsc:  " buffer.WriteString(` %s=\"`)\n var esc%d = %s\n buffer.WriteString(`\"`);",
	TagArgUne:  " buffer.WriteString(` %s=\"`)\n var unesc%d = %s\n buffer.WriteString(`\"`);",
	TagArgStr:  " buffer.WriteString(` %s=\"%s\"`);",
	TagArgAdd:  `%s + " " + %s`,
	TagArgBgn:  "`);",
	TagArgEnd:  "buffer.WriteString(`",

	CondIf:     "\nif %s {",
	CondUnless: "\nif !%s {",
	CondCase:   "\nswitch %s {",
	CondWhile:  "\nfor %s {",
	CondFor:    "\nfor %s, %s := range %s {",
	CondEnd:    "\n}",
	CondForIf:  "\nif len(%s) > 0 { for %s, %s := range %s {",

	CodeForElse:   "\n}\n} else {",
	CodeLongcode:  "\n%s",
	CodeBuffered:  "\n var esc%d = %s",
	CodeUnescaped: "\n var unesc%d = %s",
	CodeElse:      "\n} else {",
	CodeElseIf:    "\n} else if %s {",
	CodeCaseWhen:  "\ncase %s:",
	CodeCaseDef:   "\ndefault:",
	CodeMixBlock:  "\nbuffer.Write(block)",

	TextStr:     "\nbuffer.WriteString(`%s`)",
	TextComment: "\nbuffer.WriteString(`<!-- %s -->`)",

	MixinBgn:         "\n{ %s",
	MixinEnd:         "}\n",
	MixinVarBgn:      "\nvar (",
	MixinVar:         "\n%s = %s",
	MixinVarRest:     "\n%s = %#v",
	MixinVarEnd:      "\n)\n",
	MixinVarBlockBgn: "var block []byte\n{\nbuffer := new(bytes.Buffer)",
	MixinVarBlock:    "var block []byte",
	MixinVarBlockEnd: "\nblock = buffer.Bytes()\n}\n",
}

type layout struct {
	Package string
	Import  []string
	Def     []string
	Bbuf    string
	Func    string
	Before  string
	After   string
}

func (data *layout) writeBefore(wr io.Writer) {
	t := template.Must(template.New("file_bgn").Parse(file_bgn))
	err := t.Execute(wr, data)
	if err != nil {
		log.Fatalln("executing template: ", err)
	}
}
func (data *layout) writeAfter(wr *bytes.Buffer) {
	t := template.Must(template.New("file_end").Parse(file_end))
	err := t.Execute(wr, struct{ After string }{data.After})
	if err != nil {
		log.Fatalln("executing template: ", err)
	}
}

func newLayout(constName string) layout {
	var tpl layout
	tpl.Package = pkg_name

	tpl.Import = []string{
		`"bytes"`,
		`"io"`,
		`"fmt"`,
		`"html"`,
		`"strconv"`,
		`"github.com/Joker/hpp"`,
		`pool "github.com/valyala/bytebufferpool"`,
	}

	if !inline {
		tpl.Def = []string{"const ()"}
	}

	if writer {
		tpl.Bbuf = "wr io.Writer"
		tpl.Before = "buffer := &WriterAsBuffer{wr}"
	} else if stdbuf {
		tpl.Bbuf = "buffer *bytes.Buffer"
	} else {
		tpl.Bbuf = "buffer *pool.ByteBuffer"
	}

	if format {
		tpl.Before = `
			r, w := io.Pipe()
			go func() {
				buffer := &WriterAsBuffer{w}`

		if writer {
			tpl.After = `
				w.Close()
			}()
			hpp.Format(r,wr)`
		} else {
			tpl.After = `
				w.Close()
			}()
			hpp.Format(r,buffer)`
		}
	}

	//

	goFilter := jade.UseGoFilter()

	if goFilter.Name != "" {
		tpl.Func = "func " + goFilter.Name
		goFilter.Name = ""
	} else {
		tpl.Func = `func Jade_` + constName
	}

	if goFilter.Args != "" {
		args := strings.Split(goFilter.Args, ",")
		buffer := true
		for k, v := range args {
			args[k] = strings.Trim(v, " \t\n")
			if strings.HasPrefix(args[k], "buffer ") {
				args[k] = tpl.Bbuf
				buffer = false
			}
		}
		if buffer {
			args = append(args, tpl.Bbuf)
		}
		tpl.Func += "(" + strings.Join(args, ",") + ")"
		goFilter.Args = ""
	} else {
		tpl.Func += `(` + tpl.Bbuf + `) `
	}

	if goFilter.Import != "" {
		imp := strings.Split(goFilter.Import, "\n")
		for k, v := range imp {
			str := strings.Trim(v, " \t")
			if v[len(v)-1:] != `"` { // lastChar != `"`
				imp[k] = `"` + str + `"`
			} else {
				imp[k] = str
			}
		}
		tpl.Import = append(tpl.Import, imp...)
		goFilter.Import = ""
	}
	return tpl
}
