{{ template "tmpl.banner" }}

version: {{ if .Values.tmpl -}}
{{ .Values.tmpl.version }}
{{- else -}}
{{ template "tmpl.version" }}
{{- end }}
