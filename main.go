// Command jtpl executes a template on JSON data.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	textTemplate "text/template"
)

func main() {
	file := flag.String("f", "", "read template from file")
	text := flag.Bool("text", false, "use text/template instead of html/template")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [-text] { -f template-file | template-text } < data > output\n",
			flag.CommandLine.Name())
		flag.PrintDefaults()
	}
	flag.Parse()

	var tpl string
	if *file != "" {
		data, err := ioutil.ReadFile(*file)
		if err != nil {
			fmt.Fprint(os.Stderr, flag.CommandLine.Name(), ": ", err, "\n")
			os.Exit(1)
		}
		tpl = string(data)
	} else if flag.NArg() == 1 {
		tpl = flag.Arg(0)
	} else {
		flag.Usage()
		os.Exit(1)
	}

	err := run(os.Stdout, tpl, os.Stdin, *text)
	if err != nil {
		fmt.Fprint(os.Stderr, flag.CommandLine.Name(), ": ", err, "\n")
		os.Exit(1)
	}
}

type anyTemplate interface {
	Execute(io.Writer, interface{}) error
}

func run(out io.Writer, tpl string, data io.Reader, useText bool) error {
	var (
		t   anyTemplate
		err error
	)
	if useText {
		t, err = textTemplate.New("").Funcs(tplFuncs).Parse(tpl)
	} else {
		t, err = template.New("").Funcs(tplFuncs).Parse(tpl)
	}
	if err != nil {
		return err
	}
	dec := json.NewDecoder(data)
	dec.UseNumber()
	var d interface{}
	err = dec.Decode(&d)
	// EOF is acceptable, since we may receive no input.
	if err != nil && err != io.EOF {
		return err
	}
	return t.Execute(out, d)
}
