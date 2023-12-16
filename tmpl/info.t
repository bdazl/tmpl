{{ template "tmpl.banner" }}

version: {{ if .tmpl -}}
{{ .tmpl.version }}
{{- else -}}
{{ template "tmpl.version" }}
{{- end }}
