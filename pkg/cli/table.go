package cli

import (
	"os"
	"strings"

	"github.com/liggitt/tabwriter"
)

type table struct {
	w   *tabwriter.Writer
	err error
}

func newTable(cols ...string) *table {
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', tabwriter.RememberWidths)
	_, err := w.Write([]byte(strings.Join(cols, "\t") + "\n"))
	return &table{w: w, err: err}
}

func (t *table) Flush() {
	_ = t.w.Flush()
}

func (t *table) Err() error {
	if t.err != nil {
		return t.err
	}
	return t.w.Flush()
}

func (t *table) WriteRow(cols ...string) {
	if t.err != nil {
		return
	}
	_, t.err = t.w.Write([]byte(strings.Join(cols, "\t") + "\n"))
}
