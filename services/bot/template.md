# Open Positions

*Numbers may be incorrect*

| Company | Open Positions | Employees | Revenue |
|---|---|---|---|{{range .}}
| [{{.Attributes.Name}}]({{.Attributes.WebsiteUrl}}) | [{{.Attributes.OpenPositionsCount}} Open Positions]({{.Attributes.OpenPositionsUrl}}) | ~{{.Attributes.EmployeesCount}} | ~{{.Attributes.Revenue}} |{{end}}
