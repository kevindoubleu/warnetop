package device

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Device struct {
	ID    int    `json:"id"`
	Rate  int    `json:"rate"`
	Model string `json:"model"`
}

func (d Device) String() string {
	return fmt.Sprintf(
		"<ID: %d, hourly rate: Rp%s, model: %s",
		d.ID,
		message.NewPrinter(language.English).Sprint(d.Rate),
		d.Model,
		">",
	)
}
