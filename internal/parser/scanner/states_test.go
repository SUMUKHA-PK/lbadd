package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_scanDoubleQuoteSymbol(t *testing.T) {
	assert := assert.New(t)

	st := &testStream{}

	sc := newTestScanner(t, `""`, st)
	sc.executeState(scanDoubleQuoteSymbol)

	sc.assertCurrentStateIsNil()

	errs := st.errors()
	assert.Empty(errs, "exptected no errors")
}
