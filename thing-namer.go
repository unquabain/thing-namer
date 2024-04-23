package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Masterminds/sprig"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v2"
)

//go:embed data/words.yaml
var words []byte

//go:embed templates/index.html
var indexRaw string
var index = template.Must(template.New(`index`).Funcs(sprig.TxtFuncMap()).Parse(indexRaw))

func projectName(wf WordFile) (string, error) {
	var adjective,
		substantive string
	adjective, err := wf.Choose(`common`, `adjective`)
	if err != nil {
		return ``, fmt.Errorf("could not pick an adjective: %w", err)
	}
	for {
		substantive, err = wf.Choose(`common`, `substantive`)
		if err != nil {
			return ``, fmt.Errorf("could not pick an substantive: %w", err)
		}
		if substantive != adjective {
			break
		}
	}
	title := cases.Title(language.English)
	return title.String(fmt.Sprintf(`%s %s`, adjective, substantive)), nil
}

func (wf WordFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pn, err := projectName(wf)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err.Error())
		return
	}
	w.Header().Add(`Content-Type`, `text/html`)
	context := struct{ ProjectName string }{pn}
	if err := index.Execute(w, context); err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err.Error())
		return
	}
}

func main() {
	wf := make(WordFile)

	yaml.Unmarshal(words, wf)
	http.ListenAndServe(`:9099`, wf)

	/*
		var num int
		flag.IntVar(&num, `n`, 1, `Number of project titles to output`)
		flag.Parse()

		if num == 1 {
			pName, err := projectName(wf)
			if err != nil {
				fmt.Printf("Your project cannot be named: %v\n", err)
				return
			}
			fmt.Printf("Your project is now called \"%s\"\n", pName)
			return
		}

		for ; num > 0; num-- {
			pName, err := projectName(wf)
			if err != nil {
				fmt.Printf("Your project cannot be named: %v\n", err)
				return
			}
			fmt.Println(pName)
		}
	*/
}
