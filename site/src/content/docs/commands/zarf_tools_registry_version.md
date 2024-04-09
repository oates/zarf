---
title: zarf tools registry version
description: Zarf CLI command reference for <code>zarf tools registry version</code>.
tableOfContents: false
---

<!-- Page generated by Zarf; DO NOT EDIT -->

## zarf tools registry version

Print the version

### Synopsis

The version string is completely dependent on how the binary was built, so you should not depend on the version format. It may change without notice.

This could be an arbitrary string, if specified via -ldflags.
This could also be the go module version, if built with go modules (often "(devel)").

```
zarf tools registry version [flags]
```

### Options

```
  -h, --help   help for version
```

### Options inherited from parent commands

```
      --allow-nondistributable-artifacts   Allow pushing non-distributable (foreign) layers
      --insecure                           Allow image references to be fetched without TLS
      --platform string                    Specifies the platform in the form os/arch[/variant][:osversion] (e.g. linux/amd64). (default "all")
  -v, --verbose                            Enable debug logs
```

### SEE ALSO

* [zarf tools registry](/commands/zarf_tools_registry/)	 - Tools for working with container registries using go-containertools
