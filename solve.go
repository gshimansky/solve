package main

import "flag"
import "html/template"
import "io"
import "math/rand"
import "os"
import "sync"
import "time"

const NUM_IN_ROW = 4

const (
	header = `{{define "header"}}<html><head>
<style>
table, th, td {
    border: 1px solid black;
    border-collapse: collapse;
}
th, td {
    padding: 15px;
}
.interline {
    position: relative;
    top: 0.5em;
}
</style>
</head><body><h1>Generated on {{.}}</h1>{{end}}`

	solveTemplate = `
<table style="width:100%">
{{range $index, $element := .}}{{if rowstart $index}}<tr>{{end}}<td><code>{{if lt (indexinc $index) 10}}&nbsp;{{end}}{{indexinc $index}})&nbsp;<span class="interline">{{if .Add}}+{{else}}-{{end}}</span>&nbsp;&nbsp;{{if lt .First 100}}&nbsp;{{end}}{{if lt .First 10}}&nbsp;{{end}}{{.First}}<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<u>{{if lt .Second 100}}&nbsp;{{end}}{{if lt .Second 10}}&nbsp;{{end}}&nbsp;{{.Second}}</u><br>
<br>
<br>
<br>
<br>
<br>
<br>
</code></td>
{{if rowend $index}}</tr>
{{end}}{{end}}</table>`

	answersTemplate = `<tt>
{{range $index, $element := .}}{{if lt (indexinc $index) 10}}&nbsp;{{end}}{{indexinc $index}})&nbsp;{{if lt .First 100}}&nbsp;{{end}}{{if lt .First 10}}&nbsp;{{end}}{{.First}}&nbsp;{{if .Add}}&#xff0b;{{else}}&#xFF0D;{{end}}&nbsp;{{if lt .Second 100}}&nbsp;{{end}}{{if lt .Second 10}}&nbsp;{{end}}{{.Second}} = {{if lt .Result 100}}&nbsp;{{end}}{{if lt .Result 10}}&nbsp;{{end}}{{.Result}}</br>
{{end}}</tt>`

	footer = `{{define "footer"}}</body></html>{{end}}`
)

type SolveData struct {
	First, Second, Result int
	Add bool
}

func genTemplate(output io.Writer, t *template.Template, name string, timestr string, dataChan <-chan SolveData, wg *sync.WaitGroup) {
	err := t.ExecuteTemplate(output, "header", timestr)
	if err != nil {
		panic(err)
	}

	err = t.ExecuteTemplate(output, name, dataChan)
	if err != nil {
		panic(err)
	}

	err = t.ExecuteTemplate(output, "footer", nil)
	if err != nil {
		panic(err)
	}

	wg.Done()
}

func main() {
	var err error
	numlines := flag.Int("n", 5, "number of lines")
	flag.Parse()

	fm := template.FuncMap{
		"indexinc": func(i int) int {
			return i + 1
		},
		"rowstart": func(i int) bool {
			return i % NUM_IN_ROW == 0
		},
		"rowend": func(i int) bool {
			return (i + 1) % NUM_IN_ROW == 0
		},
	}

	solve, err := os.Create("solve.html")
	if err != nil {
		println("Cannot write solve.html")
		return
	}
	answers, err := os.Create("answers.html")
	if err != nil {
		println("Cannot write answers.html")
		return
	}

	at := template.New("answers")
	at = at.Funcs(fm)
	at = template.Must(at.Parse(answersTemplate))
	at = template.Must(at.Parse(header))
	at = template.Must(at.Parse(footer))
	st := template.New("solve")
	st = st.Funcs(fm)
	st = template.Must(st.Parse(solveTemplate))
	st = template.Must(st.Parse(header))
	st = template.Must(st.Parse(footer))
	toSolve := make(chan SolveData)
	toAnswers := make(chan SolveData)

	var wg sync.WaitGroup
	wg.Add(2)
	timestr := time.Now().Format(time.RFC3339)
	go genTemplate(solve, st, "solve", timestr, toSolve, &wg)
	go genTemplate(answers, at, "answers", timestr, toAnswers, &wg)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < *numlines * NUM_IN_ROW; i++ {
		first := r.Intn(100)
		second := r.Intn(100)
		add := r.Intn(2) == 0
		result := 0
		if add {
			result = first + second
		} else {
			if first < second {
				first, second = second, first
			}
			result = first - second
		}
		data := SolveData{
			First: first,
			Second: second,
			Result: result,
			Add: add,
		}
		toSolve <- data
		toAnswers <- data
		println("Doing cell", i)
	}
	close(toSolve)
	close(toAnswers)

	wg.Wait()
	solve.Close()
	answers.Close()
}
