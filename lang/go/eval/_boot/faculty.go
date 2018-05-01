package boot

type BootReservedMacro struct {
	BootPlaceholderMacro `ko:"name=placeholder"`
	Macro                string `ko:"name=macro"`
}

//XXX: add reserved ideals to booter

func BootReservedFaculty(ideals []Ideal) Faculty {
	XXX // pass faculty macros in user program
}
