package main

import "os"
import "io"
import "time"
import "strconv"
import "sync"
import "math/rand"
import "html/template"

const NUM_IN_ROW = 4

const (
	solveTemplate = `<html><head>
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
</head><body><table style="width:100%">
{{range $index, $element := .}}
{{if rowstart $index}}<tr>{{end}}
<td><code>
{{if lt $index 10}}&nbsp;{{end}}{{$index}})&nbsp;&nbsp;&nbsp;<span class="interline">&times;</span>&nbsp;{{.First}}<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<u>{{if lt .Second 10}}&nbsp;{{end}}&nbsp;{{.Second}}</u><br>
<br>
<br>
<br>
<br>
</code></td>
{{if rowend $index}}</tr>{{end}}
{{end}}
</table></body></html>`

	answersTemplate = `<html><head></head><body><tt>
{{range $index, $element := .}}
{{if lt $index 10}}&nbsp;{{end}}{{$index}})&nbsp;{{.First}}&nbsp;&times;&nbsp;{{if lt .Second 10}}&nbsp;{{end}}{{.Second}} = {{.Result}}</br>
{{end}}
</tt></body></html>`
)

type SolveData struct {
	First, Second, Result int
}

func genTemplate(output io.Writer, t *template.Template, dataChan <-chan SolveData, wg *sync.WaitGroup) {
	err := t.Execute(output, dataChan)
	if err != nil {
		panic(err)
	}

	wg.Done()
}

func main() {
	var err error
	var numlines int64 = 0

	funcMap := template.FuncMap{
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
		println("Cannot write solve.html")
		return
	}

	st := template.Must(template.New("solve").Funcs(funcMap).Parse(solveTemplate))
	at := template.Must(template.New("answers").Parse(answersTemplate))
	toSolve := make(chan SolveData)
	toAnswers := make(chan SolveData)

	var wg sync.WaitGroup
	wg.Add(2)
	go genTemplate(solve, st, toSolve, &wg)
	go genTemplate(answers, at, toAnswers, &wg)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := int64(0); i < numlines * NUM_IN_ROW; i++ {
		first := r.Intn(89) + 11
		second := r.Intn(99)
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
