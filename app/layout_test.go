package app

import (
	"fmt"
	"testing"
)

var (
	oneSub       = "0x0"
	oneMain      = "66x66"
	halfVMain    = "66x33"
	halfVSub     = "66x33"
	halfHMain    = "33x66"
	halfHSub     = "33x66"
	thirdVMain   = "66x44"
	thirdVSub    = "66x22"
	thirdHMain   = "44x66"
	thirdHSub    = "22x66"
	quarterVMain = "66x48"
	quarterVSub  = "66x16"
	quarterHMain = "48x66"
	quarterHSub  = "16x66"
)

func TestOneHorizontalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(1).Split(Horizontal)
	main := printSize(l.main())
	sub := printSize(l.sub())
	//fmt.Printf("h: main %s sub %s\n", main, sub)
	if main != oneMain {
		t.Errorf("got %s expect %s\n", main, oneMain)
	}
	if sub != oneSub {
		t.Errorf("got %s expect %s\n", sub, oneSub)
	}
	l.Position(Left)
	left := printSize(l.Left())
	right := printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}

	l.Position(Right)
	left = printSize(l.Left())
	right = printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}
}

func TestHalfHorizontalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(Half).Split(Horizontal)
	main := printSize(l.main())
	sub := printSize(l.sub())
	//fmt.Printf("h: main %s sub %s\n", main, sub)
	if main != halfHMain {
		t.Errorf("got %s expect %s\n", main, halfHMain)
	}
	if sub != halfHSub {
		t.Errorf("got %s expect %s\n", sub, halfHSub)
	}
	l.Position(Left)
	left := printSize(l.Left())
	right := printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}

	l.Position(Right)
	left = printSize(l.Left())
	right = printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}
}

func TestThirdHorizontalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(Third).Split(Horizontal)
	main := printSize(l.main())
	sub := printSize(l.sub())
	//fmt.Printf("h: main %s sub %s\n", main, sub)
	if main != thirdHMain {
		t.Errorf("got %s expect %s\n", main, thirdHMain)
	}
	if sub != thirdHSub {
		t.Errorf("got %s expect %s\n", sub, thirdHSub)
	}
	l.Position(Left)
	left := printSize(l.Left())
	right := printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}

	l.Position(Right)
	left = printSize(l.Left())
	right = printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}
}

func TestQuarterHorizontalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(Quarter).Split(Horizontal)
	main := printSize(l.main())
	sub := printSize(l.sub())
	//fmt.Printf("h: main %s sub %s\n", main, sub)
	if main != quarterHMain {
		t.Errorf("got %s expect %s\n", main, quarterHMain)
	}
	if sub != quarterHSub {
		t.Errorf("got %s expect %s\n", sub, quarterHSub)
	}
	l.Position(Left)
	left := printSize(l.Left())
	right := printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}

	l.Position(Right)
	left = printSize(l.Left())
	right = printSize(l.Right())
	switch l.mainPos {
	case Left:
		if left != main {
			t.Errorf("got %s expect %s\n", left, main)
		}
		if right != sub {
			t.Errorf("got %s expect %s\n", right, sub)
		}
	case Right:
		if right != main {
			t.Errorf("got %s expect %s\n", right, main)
		}
		if left != sub {
			t.Errorf("got %s expect %s\n", left, sub)
		}
	}
}

func TestOneVerticalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(1)
	main := printSize(l.main())
	sub := printSize(l.sub())

	if main != oneMain {
		t.Errorf("got %s expect %s\n", main, oneMain)
	}
	if sub != oneSub {
		t.Errorf("got %s expect %s\n", sub, oneSub)
	}

	l.Position(Top)
	top := printSize(l.Top())
	bottom := printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

	l.Position(Bottom)
	top = printSize(l.Top())
	bottom = printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

}

func TestHalfVerticalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(Half)
	main := printSize(l.main())
	sub := printSize(l.sub())

	if main != halfVMain {
		t.Errorf("got %s expect %s\n", main, halfVMain)
	}
	if sub != halfVSub {
		t.Errorf("got %s expect %s\n", sub, halfVSub)
	}

	l.Position(Top)
	top := printSize(l.Top())
	bottom := printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

	l.Position(Bottom)
	top = printSize(l.Top())
	bottom = printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

}

func TestThirdVerticalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(Third)
	main := printSize(l.main())
	sub := printSize(l.sub())
	if main != thirdVMain {
		t.Errorf("got %s expect %s\n", main, thirdVMain)
	}
	if sub != thirdVSub {
		t.Errorf("got %s expect %s\n", sub, thirdVSub)
	}
	l.Position(Top)
	top := printSize(l.Top())
	bottom := printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

	l.Position(Bottom)
	top = printSize(l.Top())
	bottom = printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

}

func TestQuarterVerticalLayout(t *testing.T) {
	l := NewLayout().SetSize(66, 66).Sections(Quarter)
	main := printSize(l.main())
	sub := printSize(l.sub())
	if main != quarterVMain {
		t.Errorf("got %s expect %s\n", main, quarterVMain)
	}
	if sub != quarterVSub {
		t.Errorf("got %s expect %s\n", sub, quarterVSub)
	}
	l.Position(Top)
	top := printSize(l.Top())
	bottom := printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

	l.Position(Bottom)
	top = printSize(l.Top())
	bottom = printSize(l.Bottom())
	switch l.mainPos {
	case Top:
		if top != main {
			t.Errorf("got %s expect %s\n", top, main)
		}
		if bottom != sub {
			t.Errorf("got %s expect %s\n", bottom, sub)
		}
	case Bottom:
		if bottom != main {
			t.Errorf("got %s expect %s\n", bottom, main)
		}
		if top != sub {
			t.Errorf("got %s expect %s\n", top, sub)
		}
	}

}

func printSize(w, h int) string {
	size := fmt.Sprintf("%dx%d", w, h)
	return size
}
