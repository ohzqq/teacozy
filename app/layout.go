package app

const (
	Top = iota + 1
	Right
	Bottom
	Left
)

const (
	Half = iota + 2
	Third
	Quarter
)

type Layout struct {
	split    int
	mainPos  int
	sections int
	width    int
	height   int
}

func NewLayout() *Layout {
	w, h := TermSize()
	l := &Layout{
		split:    Vertical,
		sections: 1,
		width:    w,
		height:   h,
	}
	switch l.split {
	case Vertical:
		l.mainPos = Top
	case Horizontal:
		l.mainPos = Left
	}
	return l
}

func (l *Layout) SetSize(w, h int) *Layout {
	l.width = w
	l.height = h
	return l
}

func (l *Layout) Top() (int, int) {
	switch l.mainPos {
	case Top:
		return l.Main()
	case Bottom:
		return l.Sub()
	default:
		return 0, 0
	}
}

func (l *Layout) Left() (int, int) {
	switch l.mainPos {
	case Left:
		return l.Main()
	case Right:
		return l.Sub()
	default:
		return 0, 0
	}
}

func (l *Layout) Bottom() (int, int) {
	switch l.mainPos {
	case Top:
		return l.Sub()
	case Bottom:
		return l.Main()
	default:
		return 0, 0
	}
}

func (l *Layout) Right() (int, int) {
	switch l.mainPos {
	case Left:
		return l.Sub()
	case Right:
		return l.Main()
	default:
		return 0, 0
	}
}

func (l *Layout) Main() (int, int) {
	w, h := l.width, l.height
	switch l.split {
	case Vertical:
		if l.sections > 1 {
			h = h / l.sections * (l.sections - 1)
		}
	case Horizontal:
		if l.sections > 1 {
			w = w / l.sections * (l.sections - 1)
		}
	}
	return w, h
}

func (l *Layout) Sub() (int, int) {
	w, h := 0, 0
	switch l.split {
	case Vertical:
		if l.sections > 1 {
			h = l.height / l.sections
			w = l.width
		}
	case Horizontal:
		if l.sections > 1 {
			w = l.width / l.sections
			h = l.height
		}
	}
	return w, h
}

func (l *Layout) GetSize(w, h int) (int, int) {
	switch l.split {
	case Vertical:
		return w, h / l.sections
	case Horizontal:
		return w / l.sections, h
	default:
		return w, h
	}
}

func (l *Layout) Position(p int) *Layout {
	l.mainPos = p
	return l
}

func (l *Layout) Split(s int) *Layout {
	l.split = s
	return l
}

func (l *Layout) Sections(s int) *Layout {
	l.sections = s
	return l
}

//func (l Layout) Thirds(w, h int) (int, int) {
//}
