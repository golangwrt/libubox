package libubox

import "testing"

var (
	buf *BlobBuf
)
func init() {
	buf = NewBlobBuf()
}

func TestBlobBuf(t *testing.T) {
	buf.Init(0)
	buf.AddString("name", "abc")
	t.Logf("%s\n", buf.Head().FormatJSON(true))
}

func TestBlobBuf_Printf(t *testing.T) {
	buf.Init(0)
	buf.AddString("name", "abc")
	buf.Printf("", `{"age":33}`)
	t.Logf("%s\n", buf.Head().FormatJSON(true))
}