package main

import "os"
import "fmt"
import "time"
import "strconv"
import "math/rand"

const NUM_IN_ROW = 4
const DIVSPACE = 10

func main() {
	var err error
	var numlines int64 = 0
	
	if len(os.Args) < 2 {
		println("Usage solve <number of lines>")
	} else {
		if numlines, err = strconv.ParseInt(os.Args[1], 10, 32); err != nil {
			println("Cannot parse " + os.Args[1])
			return
		}
	}

	out, err := os.Create("solve.html")
	if err != nil {
		println("Cannot write solve.html")
		return
	}
	fmt.Fprint(out, `<html><head>
<style>
table, th, td {
    border: 1px solid black;
    border-collapse: collapse;
}
th, td {
    padding: 15px;
}
</style>
</head><body><table style=\"width:100%\">`)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	count := 0
	answers := make([]int, NUM_IN_ROW * numlines)
	for i := int64(0); i < numlines; i++ {
		fmt.Fprint(out, "<tr>")
		for j := 0; j < NUM_IN_ROW; j++ {
			first := r.Intn(89) + 11
			second := r.Intn(99)
			var pad string = ""
			if count + 1 < 10 {
				pad = "&nbsp;"
			}
			if second < 10 {
				fmt.Fprintf(out, `<td><code>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;%d<br>
%s%d)&nbsp;&nbsp;&nbsp;&nbsp;<u>&nbsp;&nbsp;%d</u><br>
<br>
<br>
<br>
<br>
</code></td>`, first, pad, count + 1, second)
			} else {
				fmt.Fprintf(out, `<td><code>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;%d<br>
%s%d)&nbsp;&nbsp;&nbsp;&nbsp;<u>&nbsp;%d</u><br>
<br>
<br>
<br>
<br>
</code></td>`, first, pad, count + 1, second)
			}
			answers[count] = first * second
			count++
		}
		fmt.Fprint(out, "</tr>")
	}
	fmt.Fprint(out, "</table></body></html>")
	out.Close()

	out, err = os.Create("answers.txt")
	for i := 0; i < count; i++ {
		fmt.Fprintf(out, "%2d) %d\n", i + 1, answers[i])
	}
	out.Close()
}
