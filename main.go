package main

import (
	_ "embed"
	"encoding/xml"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"golang.org/x/net/html"
)

//go:embed post.txt.tmpl
var tmpl string

type feed struct {
	XMLName xml.Name `xml:"feed"`
	ID      string   `xml:"id"`
	Link    string   `xml:"link"`
	Title   string   `xml:"title"`
	Updated string   `xml:"updated"`
	Entries []entry  `xml:"entry"`
}

type category struct {
	XMLName xml.Name `xml:"category"`
	Term    string   `xml:"term,attr"`
}

type content struct {
	XMLName xml.Name `xml:"content"`
	Body    string   `xml:",innerxml"`
}

type entry struct {
	XMLName  xml.Name `xml:"entry"`
	Title    string   `xml:"title"`
	Category category `xml:"category"`
	Content  content  `xml:"content"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	res, err := http.Get("https://www.thegazette.co.uk/all-notices/data.feed?results-page=1")

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	doc, err := io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, doc)
	}

	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node)
	var b strings.Builder

	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			s := strings.TrimSpace(n.Data)
			t := strings.ReplaceAll(s, "…", "")

			b.WriteString(t)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	var out feed

	if err := xml.Unmarshal(doc, &out); err != nil {
		log.Fatal(err)
	}

	for _, entry := range out.Entries {
		if n, _ := html.Parse(strings.NewReader(entry.Content.Body)); n != nil {
			f(n)
		}
	}

	u := b.String()
	v := strings.Fields(u)

	rand.Shuffle(len(v), func(i, j int) {
		v[i], v[j] = v[j], v[i]
	})

	t := template.Must(template.New("default").Parse(tmpl))
	d := strings.Join(v, " ")
	p := struct {
		Date time.Time
		Content string
	}{
		Date: time.Now(),
		Content: d,
	}

	if err := t.Execute(os.Stdout, p); err != nil {
		log.Fatal(err)
	}
}
