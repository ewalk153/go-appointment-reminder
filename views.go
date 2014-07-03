package main

import (
	"fmt"
)

func DiretionsXml(redirectTo string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Say>Your appointment is located in a building near a road.</Say>
    <Redirect method="POST">%s</Redirect>
</Response>`, redirectTo)
}

func GoodbyeXml() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Say>Good bye.</Say>
</Response>`
}

func ReminderXml(postTo string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
	<Gather action="%s" numDigits="1">
    	<Say>Hello this is a call from Twilio.  You have an appointment tomorrow at 9 AM.</Say>
    	<Say>Please press 1 to repeat this menu. Press 2 for directions. Or press 3 if you are done.</Say>
	</Gather>
</Response>`, postTo)
}
