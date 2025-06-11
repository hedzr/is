package color

// c1 control code
type c1ccS struct {
	to EraseTo
	*Cursor
}

// func (s c1ccS) Printf(format string, args ...any) *Cursor {
// 	// _, _ = s.cS.sb.WriteString(csi)
// 	// if s.n > 0 {
// 	// 	_, _ = s.cS.sb.WriteString(strconv.Itoa(s.n))
// 	// }
// 	// if s.m > 0 {
// 	// 	_ = s.cS.sb.WriteByte(';')
// 	// 	_, _ = s.cS.sb.WriteString(strconv.Itoa(s.m))
// 	// }
// 	// _ = s.cS.sb.WriteByte(s.ch)

// 	return s.Cursor.Printf(format, args...)
// }

// func (s c1ccS) Echo(args ...string) *Cursor {
// 	s.Cursor.Echo(args...)
// 	return s.Cursor
// }

// func (s c1ccS) Println(args ...any) *Cursor {
// 	s.Cursor.Println(args...)
// 	return s.Cursor
// }

// func (s c1ccS) Print(args ...any) *Cursor {
// 	s.Cursor.Print(args...)
// 	return s.Cursor
// }

func (s *Cursor) EraseLine() *Cursor {
	eraseLine(s.w, CursorEraseAll)
	return s
}

type EraseTo int

const (
	CursorEraseToEnd EraseTo = iota
	CursorEraseToBegin
	CursorEraseAll
)
