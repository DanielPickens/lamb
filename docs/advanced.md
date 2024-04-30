---
meta:
  - name: description
    content: "danielpickens lamb | Documentation on customizing behavior and output"
---

# Advanced Usage Options

lamb has a wide variety of options that can be used to customize behavior and output.

## Display Options

In addition to the standard output, lamb can output in the following modes: Wide, YAML, JSON, CSV or Markdown.

`--no-headers` option hides headers in the outputs for Text, CSV and Markdown output.

### Wide

The wide output provides more information about when an apiVersion was removed or deprecated.

```shell
$ lamb detect-helm -owide
NAME                                         NAMESPACE               KIND                           VERSION                                REPLACEMENT                       DEPRECATED   DEPRECATED IN   REMOVED   REMOVED IN
cert-manager/cert-manager-webhook            cert-manager            MutatingWebhookConfiguration   admissionregistration.k8s.io/v1beta1   admissionregistration.k8s.io/v1   true         v1.16.0         false     v1.19.0
```

### JSON

```shell
$ lamb detect-helm -ojson | jq .
{
  "items": [
    {
      "name": "cert-manager/cert-manager-webhook",
      "namespace": "cert-manager",
      "api": {
        "version": "admissionregistration.k8s.io/v1beta1",
        "kind": "MutatingWebhookConfiguration",
        "deprecated-in": "v1.16.0",
        "removed-in": "v1.19.0",
        "replacement-api": "admissionregistration.k8s.io/v1",
        "component": "k8s"
      },
      "deprecated": true,
      "removed": false
    }
  ],
  "target-versions": {
    "cert-manager": "v0.15.1",
    "istio": "v1.6.0",
    "k8s": "v1.16.0"
  }
}

```

### YAML

```yaml
items:
  - name: cert-manager/cert-manager-webhook
    namespace: cert-manager
    api:
      version: admissionregistration.k8s.io/v1beta1
      kind: MutatingWebhookConfiguration
      deprecated-in: v1.16.0
      removed-in: v1.19.0
      replacement-api: admissionregistration.k8s.io/v1
      component: k8s
    deprecated: true
    removed: false
target-versions:
  cert-manager: v0.15.1
  istio: v1.6.0
  k8s: v1.16.0
```

### Custom columns

```shell
$ lamb detect-helm -ocustom --columns NAMESPACE,NAME
NAME                                         NAMESPACE
cert-manager/cert-manager-webhook            cert-manager
```

NOTE: Any columns with spaces will need to be escaped or quoted, such as `DEPRECATED\ IN` or `"DEPRECATED IN"`

### Markdown

```shell
$ lamb detect-files -o markdown
|   NAME    |   NAMESPACE    |    KIND    |      VERSION       | REPLACEMENT | DEPRECATED | DEPRECATED IN | REMOVED | REMOVED IN |
|-----------|----------------|------------|--------------------|-------------|------------|---------------|---------|------------|
| utilities | <UNKNOWN>      | Deployment | extensions/v1beta1 | apps/v1     | true       | v1.9.0        | true    | v1.16.0    |
| utilities | json-namespace | Deployment | extensions/v1beta1 | apps/v1     | true       | v1.9.0        | true    | v1.16.0    |
| utilities | yaml-namespace | Deployment | extensions/v1beta1 | apps/v1     | true       | v1.9.0        | true    | v1.16.0    |
```

```shell
$ lamb detect-files -o markdown --columns NAMESPACE,NAME,DEPRECATED IN,DEPRECATED,REPLACEMENT,VERSION,KIND,COMPONENT,FILEPATH
|     NAME      |    NAMESPACE    |    KIND    |      VERSION       | REPLACEMENT | DEPRECATED | DEPRECATED IN | COMPONENT |   FILEPATH   |
|---------------|-----------------|------------|--------------------|-------------|------------|---------------|-----------|--------------|
| some name one | lamb-namespace | Deployment | extensions/v1beta1 | apps/v1     | true       | v1.9.0        | foo       | path-to-file |
| some name two | <UNKNOWN>       | Deployment | extensions/v1beta1 | apps/v1     | true       | v1.9.0        | foo       | <UNKNOWN>    |
```

### CSV

```shell
$ lamb detect-helm -o csv
NAME,NAMESPACE,KIND,VERSION,REPLACEMENT,DEPRECATED,DEPRECATED IN,REMOVED,REMOVED IN
deploy1,lamb-namespace,Deployment,extensions/v1beta1,apps/v1,true,v1.9.0,true,v1.16.0
deploy1,other-namespace,Deployment,extensions/v1beta1,apps/v1,true,v1.9.0,true,v1.16.0
```

```shell
$ lamb detect-helm -o csv --columns "KIND,NAMESPACE,NAME,VERSION,REPLACEMENT"
KIND,NAMESPACE,NAME,VERSION,REPLACEMENT
Deployment,lamb-namespace,deploy1,extensions/v1beta1,apps/v1
Deployment,other-namespace,deploy1,extensions/v1beta1,apps/v1
```

## CI Pipelines

lamb has specific exit codes that is uses to indicate certain results:

- Exit Code 1 - An error. A message will be displayed
- Exit Code 2 - A deprecated apiVersion has been found.
- Exit Code 3 - A removed apiVersion has been found.
- Exit Code 4 - A replacement apiVersion is unavailable in the target version

If you wish to bypass the generation of exit codes 2 and 3, you may do so with two different flags:

```shell
--ignore-deprecations              Ignore the default behavior to exit 2 if deprecated apiVersions are found.
--ignore-removals                  Ignore the default behavior to exit 3 if removed apiVersions are found.
--ignore-unavailable-replacements  Ignore the default behavior to exit 4 if deprecated but unavailable apiVersions are found.
```

## Target Versions

lamb was originally designed with deprecations related to Kubernetes v1.16.0. As more deprecations are introduced, i'll will try to keep it updated. Community contributions are welcome in this area.

Currently, lamb defaults to a targetVersion of v1.22.0, however this is configurable (please continue reading)

You can target the version you are concerned with by using the `--target-versions` or `-t` flag. You must pass the `component=version`, and the version must begin with a `v` (this is a limitation of the semver library we are using to verify).

You can target the types you are concerned with by using the `--target-types` or `y` flag. You must pass the `component=type`, and if the type must begin with `x`.

For example:

```shell
$ lamb detect-helm --target-versions k8s=v1.15.0
No output to display

$ echo $?
0
```

Notice that there is no output, despite the fact that we might have recognized apiVersions present in the cluster that are not yet deprecated or removed in v1.15.0. This particular run exited 0. If there are no subAPITypes present in any clusters that are not deprecated or removed then there will also be no output and or 0 value.

## Components

By default lamb will scan for all components in the versionsList that it can find. If you wish to only see deprecations for a specific component, you can use the `--components` flag to specify a list.

## Only Show Removed

If you are targeting an upgrade, you may only wish to see apiVersions that have been `removed` rather than both `deprecated` and `removed`. You can pass the `--only-show-removed` or `-r` flag for this. It will remove any detections that are deprecated, but not yet removed. This will affect the exit code of the command as well as the json and yaml output.

## Adding Custom Version Checks

If you want to check additional apiVersions and/or types, you can pass an additional file with the `--additional-versions` or `-f` flag.

The file should look something like this:

```yaml
target-versions:
  custom: v1.0.0
deprecated-versions:
  - version: someother/v1beta1
    kind: AnotherCRD
    deprecated-in: v1.9.0
    removed-in: v1.16.0
    replacement-api: apps/v1
    component: custom
target-types:
  custom:v1.0.0
deprecated-types:
  - type: someother/v1beta1
    kind: deprecated/CRUD
    deprecated-in: v1.9.0
    removed-in: v1.16.0
    replacement-api: apps/v1
    component: custom


```

You can test that it's working by using `list-versions`:

```shell
$ lamb list-versions -f new.yaml
KIND                           NAME                                   DEPRECATED IN   REMOVED IN   REPLACEMENT   COMPONENT
AnotherCRD                     someother/v1beta1                      v1.9.0          v1.16.0      apps/v1       custom
```

_NOTE: This output is truncated to show only the additional version. Normally this will include the defaults as well_

The `target-versions` field in this custom file will set the default target version for that component. You can still override this with `--target-versions custom=vX.X.X` when you run lamb.

Please note that we do not allow overriding anything contained in the default `versions.yaml` that lamb uses.

## Kube Context

When doing helm detection, you may want to use the `--kube-context` to specify a particular context you wish to use in your kubeconfig.

## Environment Variables

For easier use, you can specify flags by using environment variables.

### Precedence

When you run a command with a flag, the command line option takes precedence over the environment variable.

### Supported Environment Variables

All environment variables are prefixed with `lamb` and use `_` instead of `-`.

| Flag                  | ENV variable              |
| --------------------- | ------------------------- |
| --ignore-deprecations | lamb_IGNORE_DEPRECATIONS |
| --ignore-removals     | lamb_IGNORE_REMOVALS     |
| --target-versions     | lamb_TARGET_VERSIONS     |
| --only-show-removed   | lamb_ONLY_SHOW_REMOVED   |
| --additional-versions | lamb_ADDITIONAL_VERSIONS |
| --output              | lamb_OUTPUT              |
| --columns             | lamb_COLUMNS             |
| --components          | lamb_COMPONENTS          |
| --no-headers          | lamb_NO_HEADERS          |
| --no-footer           | lamb_NO_FOOTER           |
