package spinner

import (
	"time"

	"github.com/briandowns/spinner"
)

var s = spinner.New(spinner.CharSets[36], 50*time.Millisecond)

func Start() { s.Start() }
func Stop()  { s.Stop() }
