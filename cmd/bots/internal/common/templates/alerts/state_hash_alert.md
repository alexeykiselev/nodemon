```yaml
Alert type: State Hash ❌
Details: Nodes on the same chain have diverging state hashes at {{ .SameHeight}}
{{ with .FirstGroup }}
BlockID (First group): {{ .BlockID}}
State Hash (First group): {{ .StateHash}}{{range .Nodes}}
{{.}}{{end}}{{end}}
{{ with .SecondGroup }}
BlockID (Second group): {{ .BlockID}}
State Hash (Second group): {{ .StateHash}}{{range .Nodes}}
{{.}}{{end}}{{end}}
{{ if .LastCommonStateHashExist }}
Fork occurred after block {{ .ForkHeight}}
BlockID: {{ .ForkBlockID}}
State Hash: {{ .ForkStateHash}}
{{ end }}
```
