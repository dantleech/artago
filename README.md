Artag
=====

**Art**ifact**ag**gregator aggregates artifacts.

This is my first Go project, so the code is even more awful than normal.

Installation
------------

Download this and build it:

```
$ git clone git@github.com/dantleech/artag
```

```
$ go build
```

Usage
-----

Create configuration file `artag.yml` in a new directory:

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
$ artag
```

Use CuRL to send an artifact:

```
$ curl -XPOST -F "data=@/home/daniel/www/phpbench/phpbench/.phpbench/html/index.html" http://localhost:8080/artifact/upload
```

You can now view the processed artifacts on `https://localhost:8080`.
