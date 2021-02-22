package main

import (
	"html/template"
	"os"
)

var gallerytemplate = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<link rel="icon" href="data:;base64,iVBORw0KGgo="/>
<title>Gallery</title>
<style>
html, body {
	margin: 0;
	padding: 0;
}
.folder {
	display: block;
	max-width: 30rem;
	border: 1px solid black;
	border-radius: 0.5rem;
	margin: 0.5rem;
	padding: 0.5rem;
}
.view {
	display: block;
	float: left;
}
img, video {
	display: block;
	margin: 0;
	padding: 0;
}
.binary {
	padding: 0.5rem;
	clear: both;
}
</style>
</head>
<body>

<div class="folders">
{{- range .Folders -}}
<a class="folder" href="{{ .Path }}">{{ .Name }}</a>
{{- end -}}
</div>

{{- range .Files -}}

{{- if .Image -}}
<a class="view" href="{{ .ViewPath }}"><img src="{{ .ThumbPath }}" width="{{ .ThumbWidth }}" height="{{ .ThumbHeight }}" /></a>
{{- else -}}
{{- if .Video -}}
<a class="view" href="{{ .ViewPath }}"><video autoplay=true loop=true muted=true src="{{ .ThumbPath }}" width="{{ .ThumbWidth }}" height="{{ .ThumbHeight }}" /></a>
{{- else -}}
<div class="binary"><a href="{{ .Path }}">{{ .Path }}</a> ({{ .SizeNice }})</div>
{{- end -}}
{{- end -}}

{{- end -}}

</body>
</html>
`))

var filetemplate = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<link rel="icon" href="data:;base64,iVBORw0KGgo="/>
<title>File</title>
<style>
html, body {
	margin: 0;
	padding: 0;
}
img, video {
	display: block;
	margin: 0;
	padding: 0;
	width: 100%;
}
</style>
</head>
<body>

<a href="../">Back</a> |
{{ .Path }} <a href="{{ .Path }}">Original</a> ({{ .SizeNice }})

{{ if .Image }} | <a href="{{ .BigPath }}">Large</a> ({{ .BigSizeNice }}){{ end }}
{{ if .Video }} | <a href="{{ .BigPath }}">Large</a> ({{ .BigSizeNice }}){{ end }}

<br />

{{- if .Image -}}
<img src="{{ .BigPath }}" />
{{- end -}}
{{- if .Video -}}
<video controls=true autoplay=true src="{{ .BigPath }}" />
{{- end -}}

</body>
</html>
`))

func renderGallery(gallery *Gallery) error {
	fh, err := os.Create(gallery.Dest)
	if err != nil {
		return err
	}
	defer fh.Close()

	if err := gallerytemplate.Execute(fh, gallery); err != nil {
		return err
	}

	if err := fh.Close(); err != nil {
		return err
	}
	return nil
}

func renderFile(file *File) error {
	dir := file.Dest + "_view/"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	fh, err := os.Create(dir + "index.html")
	if err != nil {
		return err
	}
	defer fh.Close()

	if err := filetemplate.Execute(fh, file); err != nil {
		return err
	}

	if err := fh.Close(); err != nil {
		return err
	}
	return nil
}
