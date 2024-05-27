package csv

import (
	"encoding/csv"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

type MyData struct {
	Name   string `csv:"name"`
	HasPet bool   `csv:"has_pet"`
	Age    int    `csv:"age"`
}

func TestCSV(t *testing.T) {
	data := `name,age,has_pet
Jon,"100",true
"Fred ""The Hammer"" Smith",42,false
Martha,37,"true"
`
	r := csv.NewReader(strings.NewReader(data))
	allData, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	var entries []MyData
	Unmarshal(allData, &entries)

	//now to turn entries into output
	out, err := Marshal(entries)
	if err != nil {
		panic(err)
	}
	sb := &strings.Builder{}
	w := csv.NewWriter(sb)
	w.WriteAll(out)

	expected := `name,has_pet,age
Jon,true,100
"Fred ""The Hammer"" Smith",false,42
Martha,true,37
`
	if diff := cmp.Diff(expected, sb.String()); diff != "" {
		t.Error(diff)
	}
}
