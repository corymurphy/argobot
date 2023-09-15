{{/*
Expand the name of the chart.
*/}}
{{- define "argobot.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "argobot.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "argobot.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "argobot.labels" -}}
helm.sh/chart: {{ include "argobot.chart" . }}
{{ include "argobot.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "argobot.selectorLabels" -}}
app.kubernetes.io/name: {{ include "argobot.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "argobot.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "argobot.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Generates secret-webhook name
*/}}
{{- define "argobot.webhookSecretName" -}}
{{- if .Values.githubApp.webhookSecretName -}}
    {{ .Values.githubApp.webhookSecretName }}
{{- else -}}
    {{ template "argobot.fullname" . }}
{{- end -}}
{{- end -}}

{{/*
Generates secret-githubapp-private-key name
*/}}
{{- define "argobot.ghPrivateKeySecretName" -}}
{{- if .Values.githubApp.privateKeySecretName -}}
    {{ .Values.githubApp.privateKeySecretName }}
{{- else -}}
    {{ template "argobot.fullname" . }}
{{- end -}}
{{- end -}}
