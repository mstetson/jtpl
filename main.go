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
)

func main() {
	file := flag.String("f", "", "read template from file")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s { -f template-file | template-text } < data > output\n",
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

	err := run(os.Stdout, tpl, os.Stdin)
	if err != nil {
		fmt.Fprint(os.Stderr, flag.CommandLine.Name(), ": ", err, "\n")
		os.Exit(1)
	}
}

func run(out io.Writer, tpl string, data io.Reader) error {
	t, err := template.New("").Parse(tpl)
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
