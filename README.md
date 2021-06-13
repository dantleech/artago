Artago
======

[![Build](https://github.com/dantleech/artago/actions/workflows/go.yml/badge.svg)](https://github.com/dantleech/artago/actions/workflows/go.yml)

**Art**ifact**ag**gogeregator aggregates artifacts - you `POST` build artifacts
to it and you configure what should be done with the artifacts. Artag provides
a web server to serve the artifacts.

It's an artifact server. It serves artifacts.

This is my first Go project, so the code is even more awful than normal.

Installation
------------

Download this and build it:

```
$ git clone git@github.com/dantleech/artago
```

```
$ go build
```

Usage
-----

Create configuration file `artago.yml` in a new directory:

```yaml
address: 127.0.0.1:8080
workspacePath: workspace
publicDir: docs
rules:
  -
    rule: artifact.Name == "index.html"
    actions:
      -
         type: copy
         params:
           destination: docs/index.html
      -
         type: copy
         params:
           destination: docs/archive/%artifact.BuildId%/index.html
```

Run the server:

```
$ artago
```

Use CuRL to send an artifact (or artifacts):

```
$ curl -XPOST -F "data=@/home/daniel/www/phpbench/phpbench/.phpbench/html/index.html" http://localhost:8080/artifact/upload
```

You will see a response:

```
{
  "buildId": "20210613-135921",
  "results": {}
}
```

You can now view the processed artifacts on `https://localhost:8080`.

Custom Headers
--------------

You can specify special headers to 

### BuildId

By default the BuildId will default to the current date (YYY-MM-DD H:i:s). If
you may want to override this to correspond with your CI servers build ID:

```
$ curl  -H"BuildId:123456" // ...
```

Configuration
-------------

```yaml
address: 127.0.0.1:8080
workspacePath: workspace
publicDir: docs
rules:
  -
    rule: artifact.Name == "index.html"
    actions:
      -
         type: copy
         params:
           destination: docs/index.html
```

- `address`: The address to serve on
- `workspacePath`: Relative to the CWD, the directory to use for temporary
  files.
- `publicDir`: The directory which should be exposed by the HTTP file server
- `rules`: See next section

Rules
-----

The rules in the configuration decide what should be done with the uploaded artifacts. Each rule is an [expression](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md) followed by a list of actions to take if that expression evaluates to true for the current artifact:

```
artifact.Name == "index.html"
```

Actions
-------

### copy

Copy the artifact to the given `destination`.

```
rules:
  -
    rule: artifact.Name == "index.html"
    actions:
      -
         type: copy
         params:
           destination: docs/index.html
```

The `destination` can have an expression embedded in it - the following will
use the artifact's build ID:

```
destination: docs/%artifact.BuildId%/index.html
```

Note that paths are always relative to the current working directory of
`artago`.

### publishLink

Add links to the HTTP JSON response:

```
rules:
  -
    rule: artifact.Name == "contents.rst"
    actions:
      -
         type: publishLink
         params:
             name: documentation
             template: "http://127.0.0.1:8080/documentation.rst"
```

This will add a `links` _result_ in the HTTP response, in this case:

```
{
  "BuildId": "20210613-135921",
  "Results": {
    "links": {
      "documentation": "http://127.0.0.1:8080/documentation.rst"
    }
  }
}
```

As with copy, you can embed an expressoin in the `template` with
`%myExpression%`.
