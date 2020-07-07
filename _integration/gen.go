// +build none

// Copyright Project Contour Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"html/template"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("generator", "Generate mock integration testing data.")
	count := app.Flag("count", "Count of entries to create.").Short('c').Default("500").Int()

	services := app.Command("services", "Generate mock service data.")
	selector := services.Flag("selector", "Name of the service's selector.").Default("").String()

	rand.Seed(time.Now().UnixNano())
	args := os.Args[1:]
	switch kingpin.MustParse(app.Parse(args)) {
	case services.FullCommand():
		genServices(*count, *selector)
	default:
		app.Usage(args)
		os.Exit(2)
	}
}

const serviceTmpl = `
# autogenerated: do not edit!
# source:{{ range .args }} {{ . }}{{end}}
{{ $selector := .selector -}}
{{ range .names -}}
apiVersion: v1
kind: Service
metadata:
  name: {{ . }}
spec:
  ports:
  - port: 80
    protocol: TCP
{{ if $selector }}  selector:
    app: {{ $selector }}{{ end }}
---
{{ end}}
`

func genServices(n int, selector string) {
	var names []string
	for i := 0; i < n; i++ {
		name := namesgenerator.GetRandomName(0)
		name = strings.Replace(name, "_", "-", -1) // must be a valid rfc 1035 value
		names = append(names, name)
	}

	t, err := template.New("services").Parse(serviceTmpl[1:])
	check(err)
	check(t.Execute(os.Stdout, map[string]interface{}{
		"names":    names,
		"args":     os.Args,
		"selector": selector,
	}))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
