package run

import (
	"embed"
	"html/template"
)

//go:embed template/linux_html.html
var templateFile embed.FS

//go:embed template/oracle_html.html
var templateFileOracle embed.FS

//go:embed template/styles.css
var css template.CSS
