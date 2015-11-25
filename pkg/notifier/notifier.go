package notifier

import "log"

type Notifier interface {
	Notify(source string, status string, payload string)
}

// LogNotified notifies to logger
type LogNotifier struct {
}

func (l *LogNotifier) Notify(source string, status string, payload string) {
	log.Printf("Notifier: source:%s status:%s payload:%s", source, status, payload)
}
