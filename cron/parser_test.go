package cron

import (
	"fmt"
	"testing"
)

func TestGetSequenceFromCommaEpr(t *testing.T) {
	cases := map[string]struct {
		in     string
		out    []uint
		bounds bounds
		err    error
	}{
		"TestMinuteSimpleCase": {
			in:     "1,15",
			out:    []uint{1, 15},
			bounds: minuteBounds,
			err:    nil,
		},
		"TestMinuteLoverOutOfBound": {
			in:     "-1,15",
			out:    nil,
			bounds: minuteBounds,
			err:    fmt.Errorf("Should fail"),
		},
		// "TestMinuteLoverIsHigherThanHither": {
		// 	in:     "17,15",
		// 	out:    "0 15 30 45",
		// 	bounds: minuteBounds,
		// 	err:    true,
		// },
		// "TestMinuteHigherOutOfBound": {
		// 	in:     "17,60",
		// 	out:    "0 15 30 45",
		// 	bounds: minuteBounds,
		// 	err:    true,
		// },
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			result, err := getSequenceFromCommaEpr(c.in, c.bounds)
			if c.err == nil {
				if err != nil {
					t.Errorf("Test faild but should not")
					return
				}
			} else {
				if err == nil {
					t.Errorf("Test should not fail but it failed")
					return
				}
			}
			if !fieldRangesEqual(result, c.out) {
				t.Errorf("Results do not math. Expected: %v, got: %v", c.out, result)
			}
		})
	}
}

func fieldRangesEqual(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func Test_getSequenceFromHifenEpr(t *testing.T) {
	t.Errorf("Tests have not been implemented yet")
}
func Test_getSequenceFromSlashExpr(t *testing.T) {
	t.Errorf("Tests have not been implemented yet")
}

func Test_getRange(t *testing.T) {
	t.Errorf("Tests have not been implemented yet")
}

// TODO: finish it
func TestParse(t *testing.T) {
	t.Errorf("Tests have not been implemented yet")
	// cases := []struct{
	// 	in string
	// 	out schedule
	// 	err error
	// } {
	// 	{
	// 		args: "*/15 0 1,15 * 1-5 /usr/bin/find",
	// 		exp: {
	// 			ranges: [5][]uint{
	// 				[]uint{0, 15, 30, 45},
	// 				[]uint{0},
	// 				[]uint{1, 15},
	// 				[]uint{1,2,3,4,5,6,7,8,9,10,11,12},
	// 				[]uint{0, 15, 30, 45},
	// 			},
	// 			command: "/usr/bin/find",
	// 		},
	// 		err: nil,
	// 	},
	// }
}
