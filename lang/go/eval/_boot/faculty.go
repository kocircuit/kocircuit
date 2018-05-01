package boot

type BootMacroMacro struct {
	BootPlaceholderMacro `ko:"name=placeholder"`
	Macro                string `ko:"name=macro"`
}

func BootFaculty(ideals []string) Faculty {
	XXX // pass faculty macros in user program
}
