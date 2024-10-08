---
title: zarf tools helm dependency list
description: Zarf CLI command reference for <code>zarf tools helm dependency list</code>.
tableOfContents: false
---

<!-- Page generated by Zarf; DO NOT EDIT -->

## zarf tools helm dependency list

list the dependencies for the given chart

### Synopsis


List all of the dependencies declared in a chart.

This can take chart archives and chart directories as input. It will not alter
the contents of a chart.

This will produce an error if the chart cannot be loaded.


```
zarf tools helm dependency list CHART [flags]
```

### Options

```
  -h, --help                 help for list
      --max-col-width uint   maximum column width for output table (default 80)
```

### Options inherited from parent commands

```
      --burst-limit int                 client-side default throttling limit (default 100)
      --debug                           enable verbose output
      --insecure-skip-tls-verify        Skip checking server's certificate for validity. This flag should only be used if you have a specific reason and accept the reduced security posture.
      --kube-apiserver string           the address and the port for the Kubernetes API server
      --kube-as-group stringArray       group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --kube-as-user string             username to impersonate for the operation
      --kube-ca-file string             the certificate authority file for the Kubernetes API server connection
      --kube-context string             name of the kubeconfig context to use
      --kube-insecure-skip-tls-verify   if true, the Kubernetes API server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kube-tls-server-name string     server name to use for Kubernetes API server certificate validation. If it is not provided, the hostname used to contact the server is used
      --kube-token string               bearer token used for authentication
      --kubeconfig string               path to the kubeconfig file
  -n, --namespace string                namespace scope for this request
      --plain-http                      Force the connections over HTTP instead of HTTPS. This flag should only be used if you have a specific reason and accept the reduced security posture.
      --qps float32                     queries per second used when communicating with the Kubernetes API, not including bursting
      --registry-config string          path to the registry config file
      --repository-cache string         path to the directory containing cached repository indexes
      --repository-config string        path to the file containing repository names and URLs
```

### SEE ALSO

* [zarf tools helm dependency](/commands/zarf_tools_helm_dependency/)	 - manage a chart's dependencies

