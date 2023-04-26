package keys

func DefaultKeyMap() KeyMap {
	k := []*Binding{
		PgUp(),
		PgDown(),
		Up(),
		Down(),
		HalfPgUp(),
		HalfPgDown(),
		Home(),
		End(),
	}
	return NewKeyMap(k...)
}

func VimKeyMap() KeyMap {
	km := []*Binding{
		Up().AddKeys("k"),
		Down().AddKeys("j"),
		HalfPgUp().AddKeys("K"),
		HalfPgDown().AddKeys("J"),
		Home().AddKeys("g"),
		End().AddKeys("G"),
		Quit().AddKeys("q"),
	}
	return NewKeyMap(km...)
}
