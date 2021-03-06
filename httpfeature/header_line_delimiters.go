package httpfeature

import (
	"fmt"
)

type HeaderLineDelimiters struct {
	BaseFeature
	Symbols []string
}

func (f *HeaderLineDelimiters) Name() string {
	return "Header line delimitiers"
}

func (f *HeaderLineDelimiters) Export() interface{} {
	return f.Symbols
}

func (f *HeaderLineDelimiters) String() string {
	return PrintableStrings(f.Symbols)
}

func (f *HeaderLineDelimiters) Collect() error {
	f.Symbols, _ = f.collectSymbols()
	return nil
}

func (f *HeaderLineDelimiters) collectSymbols() ([]string, error) {
	var toCheck = [...]string{
		"\r",
		"\n",
		"\r\n",
		"\x00",
	}

	symbols := make([]string, 0)
	for _, ch := range toCheck {
		if f.check(ch) {
			symbols = append(symbols, ch)
		}
	}
	return symbols, nil
}

func (f *HeaderLineDelimiters) check(delim string) bool {
	req := f.BaseRequest.Clone()
	req.AddHeader("X-Foo", fmt.Sprintf("foo%sX-Bar:bar", delim))
	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		return false
	}

	for _, h := range resp.HeadersSlice("X-Bar") {
		if h.Value == "bar" {
			return true
		}
	}
	return false
}