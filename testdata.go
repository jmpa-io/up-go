package up

import "os"

type testdata struct {
	path    string
	content []byte
}

func (t *testdata) getContent() {
	b, _ := os.ReadFile(t.path)
	t.content = b
}

func NewTestdata(fileName string) *testdata {
	td := &testdata{
		path: "./testdata/" + fileName + ".json",
	}
	td.getContent()
	return td
}
