package logging

import "fmt"

func FormatLog(worker int, job int, log string) string {
	return fmt.Sprintf("[%d-%d]: %s", worker, job, log)
}
