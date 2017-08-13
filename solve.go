package main

import "html/template"
import "io"
import "math/rand"
import "os"
import "strconv"
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
{{range $index, $element := .}}{{if rowstart $index}}<tr>{{end}}<td><code>{{if lt (indexinc $index) 10}}&nbsp;{{end}}{{indexinc $index}})&nbsp;<span class="interline">&times;</span>&nbsp;&nbsp;{{if lt .First 100}}&nbsp;{{end}}{{if lt .First 10}}&nbsp;{{end}}{{.First}}<br>
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
{{range $index, $element := .}}{{if lt (indexinc $index) 10}}&nbsp;{{end}}{{indexinc $index}})&nbsp;{{if lt .First 100}}&nbsp;{{end}}{{if lt .First 10}}&nbsp;{{end}}{{.First}}&nbsp;&times;&nbsp;{{if lt .Second 100}}&nbsp;{{end}}{{if lt .Second 10}}&nbsp;{{end}}{{.Second}} = {{.Result}}</br>
{{end}}</tt>`

	footer = `{{define "footer"}}</body></html>{{end}}`
)

type SolveData struct {
	First, Second, Result int
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
	var numlines int64 = 0

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

	if len(os.Args) < 2 {
		println("Usage solve <number of lines>")
	} else {
		if numlines, err = strconv.ParseInt(os.Args[1], 10, 32); err != nil {
			println("Cannot parse " + os.Args[1])
			return
		}
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
	for i := int64(0); i < numlines * NUM_IN_ROW; i++ {
		first := r.Intn(999)
		second := r.Intn(999)
		data := SolveData{
			First: first,
			Second: second,
			Result: first * second,
		}
		toSolve <- data
		toAnswers <- data
	}
	close(toSolve)
	close(toAnswers)

	wg.Wait()
	solve.Close()
	answers.Close()
}
