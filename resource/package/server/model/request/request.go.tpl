package request

import (
{{- if not .OnlyTemplate }}
	"DYCLOUD/model/common/request"
	{{ if or .HasSearchTimer .GvaModel}}"time"{{ end }}
{{- end }}
)

type {{.StructName}}Search struct{
{{- if not .OnlyTemplate}}
{{- if .GvaModel }}
    StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
    EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
{{- end }}
{{- range .Fields}}
    {{- if ne .FieldSearchType ""}}
        {{- if eq .FieldSearchType "BETWEEN" "NOT BETWEEN"}}
    Start{{.FieldName}}  *{{.FieldType}}  `json:"start{{.FieldName}}" form:"start{{.FieldName}}"`
    End{{.FieldName}}  *{{.FieldType}}  `json:"end{{.FieldName}}" form:"end{{.FieldName}}"`
        {{- else }}
            {{- if or (eq .FieldType "enum") (eq .FieldType "picture") (eq .FieldType "pictures") (eq .FieldType "video") (eq .FieldType "richtext") (eq .FieldType "json") }}
            {{.FieldName}}  string `json:"{{.FieldJson}}" form:"{{.FieldJson}}" `
            {{- else if ne .FieldType "string" }}
    {{.FieldName}}  *{{.FieldType}} `json:"{{.FieldJson}}" form:"{{.FieldJson}}" `
            {{- else }}
    {{.FieldName}}  {{.FieldType}} `json:"{{.FieldJson}}" form:"{{.FieldJson}}" `
            {{- end }}
        {{- end }}
    {{- end}}
{{- end }}
    request.PageInfo
    {{- if .NeedSort}}
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
    {{- end}}
{{- end}}
}