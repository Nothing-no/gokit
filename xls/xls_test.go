package xls

import "testing"

func TestIsExisted(t *testing.T) {
	ts := []struct {
		arr     []string
		element string
		result  bool
	}{
		{[]string{"hello", "world", "good", "day"}, "hello", true},
		{[]string{"hello", "world", "good", "day"}, "day", true},
	}
	var r bool
	for _, te := range ts {
		r = IsExisted(te.arr, te.element)
		if r != te.result {
			t.Errorf("Search %v in %v,expected %t, ouput %t\n",
				te.element, te.arr, te.result, r)
		}
	}
}
func TestGetIncAxis(t *testing.T) {
	ts := []struct {
		base   string
		step   int
		col    bool
		result string
	}{
		{"A5", 5, false, "A10"},
		{"A5", 5, true, "F5"},
		{"A10", -1, false, "A9"},
		{"B5", -1, true, "A5"},
		{"D5", 1, true, "E5"},
		{"D5", 0, true, "D5"},
	}
	for _, te := range ts {
		r := GetIncAxis(te.base, te.step, te.col)
		if r != te.result {
			t.Errorf("input:%s,%d,%t,expected %s, ouput %s\n",
				te.base, te.step, te.col, te.result, r)
		}
	}
}
