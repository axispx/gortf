package gortf

func RTFToHTML(r *RtfDocument) (string, error) {
	body := r.Body

	var htmlBody string
	for _, styleBlock := range body {
		closingTagStack := []string{}
		if styleBlock.Painter.Bold {
			htmlBody += "<bold>"
			closingTagStack = append(closingTagStack, "</bold>")
		}

		if styleBlock.Painter.Italic {
			htmlBody += "<italic>"
			closingTagStack = append(closingTagStack, "</italic>")
		}

		if styleBlock.Painter.Underline {
			htmlBody += "<u>"
			closingTagStack = append(closingTagStack, "</u>")
		}

		htmlBody += styleBlock.Text

		for i := len(closingTagStack) - 1; i >= 0; i-- {
			htmlBody += closingTagStack[i]
		}
	}

	return htmlBody, nil
}
