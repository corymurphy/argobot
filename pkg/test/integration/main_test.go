package integration

import "flag"

var kubeContext = flag.String("kubecontext", "kind-kind", "kubernetes context to run tests against")
var kubeNamespace = flag.String("kubenamespace", "bienhoa", "kubernetes namespace to run tests against")
var appName = flag.String("appname", "argobot", "application to deploy")
var appVersion = flag.String("appversion", "60a4727c6c", "application version to deploy")
