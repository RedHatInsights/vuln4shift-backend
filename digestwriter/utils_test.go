package digestwriter_test

// Utility functions for the package's unit tests
import (
	"time"

	"bou.ke/monkey"
)

func patchCurrentTime() {
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	})
}
