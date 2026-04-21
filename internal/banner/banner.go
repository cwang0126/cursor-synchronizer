package banner

import (
	"fmt"
	"io"
)

const cyan = "\033[36m"
const bold = "\033[1m"
const reset = "\033[0m"

const art = `
   ____ _   _ ____  ____   ___  ____        ______   ___   _  ____
  / ___| | | |  _ \/ ___| / _ \|  _ \      / ___\ \ / / \ | |/ ___|
 | |   | | | | |_) \___ \| | | | |_) |____ \___ \\ V /|  \| | |
 | |___| |_| |  _ < ___) | |_| |  _ <_____|___) || | | |\  | |___
  \____|\___/|_| \_\____/ \___/|_| \_\    |____/ |_| |_| \_|\____|
`

// Print writes the banner to w using ANSI colors.
func Print(w io.Writer) {
	fmt.Fprint(w, cyan+bold+art+reset+"\n")
}
