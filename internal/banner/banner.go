package banner

import (
	"fmt"
	"io"
)

// Cyan is the bright cyan ANSI color used by the banner; exported so other
// packages can color related text (e.g. tagline / version) to match.
const Cyan = "\033[1;96m"
const Reset = "\033[0m"

// art is figlet's "slant" rendering of "Cursor Sync".
const art = `
   ______                              _____                 
  / ____/_  ________________  _____   / ___/__  ______  _____
 / /   / / / / ___/ ___/ __ \/ ___/   \__ \/ / / / __ \/ ___/
/ /___/ /_/ / /  (__  ) /_/ / /      ___/ / /_/ / / / / /__  
\____/\__,_/_/  /____/\____/_/      /____/\__, /_/ /_/\___/  
                                         /____/              
`

// Print writes the banner to w using ANSI colors.
func Print(w io.Writer) {
	fmt.Fprint(w, Cyan, art, Reset, "\n")
}
