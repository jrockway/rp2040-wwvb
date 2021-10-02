package screen

import (
	"image/color"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type fakeDisplay struct {
	text string
}

func (d *fakeDisplay) SetText(x string) {
	if len(x) != 64 {
		d.text = "error"
		return
	}
	d.text = strings.Join([]string{
		strings.ReplaceAll(x[0:16], "\x00", " "),
		strings.ReplaceAll(x[16:32], "\x00", " "),
		strings.ReplaceAll(x[32:48], "\x00", " "),
		strings.ReplaceAll(x[48:64], "\x00", " "),
	}, "\n")
}

func (d *fakeDisplay) SetPixel(_, _ int16, _ color.RGBA) {}
func (d *fakeDisplay) Display() error                    { return nil }

func TestScreen(t *testing.T) {
	d := new(fakeDisplay)
	s := Screen{Display: d}

	td := []struct {
		input string
		args  []interface{}
		want  [4]string
	}{
		{
			input: "1 ",
			want: [4]string{
				"1               ",
				"                ",
				"                ",
				"                ",
			},
		},
		{
			input: "2 ",
			want: [4]string{
				"1 2             ",
				"                ",
				"                ",
				"                ",
			},
		},
		{
			input: "this is a very long string that probably wraps around if it goes on long enough",
			want: [4]string{
				"goes on long eno",
				"ughlong string t",
				"hat probably wra",
				"ps around if it ",
			},
		},
		{
			input: "",
			want: [4]string{
				"                ",
				"                ",
				"                ",
				"                ",
			},
		},
		{
			input: "1\n2\n3\n4\n",
			want: [4]string{
				"1               ",
				"2               ",
				"3               ",
				"4               ",
			},
		},
		{
			input: "X",
			want: [4]string{
				"X               ",
				"2               ",
				"3               ",
				"4               ",
			},
		},
		{
			input: "",
			want: [4]string{
				"                ",
				"                ",
				"                ",
				"                ",
			},
		},
		{
			input: "0123456789ABCDEF\n2\n3\n4\n",
			want: [4]string{
				"0123456789ABCDEF",
				"2               ",
				"3               ",
				"4               ",
			},
		},
	}
	for _, test := range td {
		if test.input == "" {
			s.Clear()
		}
		s.Printf(test.input, test.args...)
		want := strings.Join(test.want[:], "\n")
		if diff := cmp.Diff(d.text, want); diff != "" {
			t.Errorf("after printing %q (%v) -got +want:\n%s", test.input, test.args, diff)
		}
	}
}
