// Command jtpl executes a template on JSON data.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s template-text < data > output\n",
			flag.CommandLine.Name())
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	err := run(os.Stdout, flag.Arg(0), os.Stdin)
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
