package up

import (
	"os"
	"time"
)

type testdata struct {
	path    string
	content []byte
}

func (t *testdata) getContent() {
	b, err := os.ReadFile(t.path)
	if err != nil {
		panic(err)
	}
	t.content = b
}

func newTestdata(fileName string) *testdata {
	td := &testdata{
		path: "./testdata/" + fileName + ".json",
	}
	td.getContent()
	return td
}

// a location used for tests.
var location = time.FixedZone("AEST", 11*60*60)
