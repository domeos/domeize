domeize ![version v0.2.0](https://img.shields.io/badge/version-v0.2.0-brightgreen.svg) ![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)
=============
domeize if forked from domeize and added some more function to meet needs for DomeOS.

Utility to simplify running applications in docker containers.

domeize is a utility to simplify running applications in docker containers.  It allows you
to generate application configuration files at container startup time from templates and
container environment variables.  It also allows log files to be tailed to stdout and/or
stderr.

The typical use case for domeize is when you have an application that has one or more
configuration files and you would like to control some of the values using environment variables.

For example, a Python application using Sqlalchemy might not be able to use environment variables directly.
It may require that the database URL be read from a python settings file with a variable named
`SQLALCHEMY_DATABASE_URI`.  domeize allows you to set an environment variable such as
`DATABASE_URL` and update the python file when the container starts.

Another use case is when the application logs to specific files on the filesystem and not stdout
or stderr. This makes it difficult to troubleshoot the container using the `docker logs` command.
For example, nginx will log to `/var/log/nginx/access.log` and
`/var/log/nginx/error.log` by default. While you can sometimes work around this, it's tedious to find
the a solution for every application. domeize allows you to specify which logs files should
be tailed and where they should be sent.

See [A Simple Way To Dockerize Applications](http://jasonwilder.com/blog/2014/10/13/a-simple-way-to-dockerize-applications/)

## Installation

Download the latest version in your container:

For Ubuntu Images:

```
RUN apt-get update && apt-get install -y wget
RUN wget http://domeos-binpack.bjcnc.scs.sohucs.com/domeize -O /usr/local/bin/domeize && chmod +x /usr/local/bin/domeize
```

## Usage

domeize works by wrapping the call to your application using the `ENTRYPOINT` or `CMD` directives.

This would generate `/etc/nginx/nginx.conf` from the template located at `/etc/nginx/nginx.tmpl` and
send `/var/log/nginx/access.log` to `STDOUT` and `/var/log/nginx/error.log` to `STDERR` after running
`nginx`, only after waiting for the `web` host to respond on `tcp 8000`:

```
CMD domeize -template /etc/nginx/nginx.tmpl:/etc/nginx/nginx.conf -stdout /var/log/nginx/access.log -stderr /var/log/nginx/error.log -wait tcp://web:8000 nginx
```

modification:
1. read from env `NEED_PORTS` to get a number `n` to allocate unused ports, and export `AUTO_PORT0`...`AUTO_PORT(n-1)` to env, then domeize can replace the var to config file. also need env var `DOMEOS_SERVER_ADDR`, `CLUSTER_NAME`, `MY_POD_NAMESPACE`, `MY_POD_NAME` to report the port info (`AUTO_PORT0`...) through DomeOS to kubenetes, then put to annotation of meta for pod info.
2. env `DOMEIZE_TEMPLATES` to pass template file list, you can use `/etc/nginx/nginx.tmpl:/etc/nginx/nginx.conf` or `/path/to/1.tmpl:/path/to/1.conf,/path/to/2.tmpl:/path/to/2.conf` to replace too much `-template` parameter

### Command-line Options

You can specify multiple templates by passing using `-template` multiple times:

```
$ domeize -template template1.tmpl:file1.cfg -template template2.tmpl:file3

```

Templates can be generated to `STDOUT` by not specifying a dest:

```
$ domeize -template template1.tmpl

```


You can tail multiple files to `STDOUT` and `STDERR` by passing the options multiple times.

```
$ domeize -stdout info.log -stdout perf.log

```

If `inotify` does not work in you container, you use `-poll` to poll for file changes instead.

```
$ domeize -stdout info.log -stdout perf.log -poll

```


If your file uses `{{` and `}}` as part of it's syntax, you can change the template escape characters using the `-delims`.

```
$ domeize -delims "<%:%>"
```

## Waiting for other dependencies

It is common when using tools like [Docker Compose](https://docs.docker.com/compose/) to depend on services in other linked containers, however oftentimes relying on [links](https://docs.docker.com/compose/compose-file/#links) is not enough - whilst the container itself may have _started_, the _service(s)_ within it may not yet be ready - resulting in shell script hacks to work around race conditions.

Dockerize gives you the ability to wait for services on a specified protocol (`tcp`, `tcp4`, `tcp6`, `http`, and `https`) before starting your application:

```
$ domeize -wait tcp://db:5432 -wait http://web:80
```

See [this issue](https://github.com/docker/compose/issues/374#issuecomment-126312313) for a deeper discussion, and why support isn't and won't be available in the Docker ecosystem itself.

## Using Templates

Templates use Golang [text/template](http://golang.org/pkg/text/template/). You can access environment
variables within a template with `.Env`.

```
{{ .Env.PATH }} is my path
```

There are a few built in functions as well:

  * `default $var $default` - Returns a default value for one that does not exist. `{{ default .Env.VERSION "0.1.2" }}`
  * `contains $map $key` - Returns true if a string is within another string
  * `exists $path` - Determines if a file path exists or not. `{{ exists "/etc/default/myapp" }}`
  * `split $string $sep` - Splits a string into an array using a separator string. Alias for [`strings.Split`][go.string.Split]. `{{ split .Env.PATH ":" }}`
  * `replace $string $old $new $count` - Replaces all occurrences of a string within another string. Alias for [`strings.Replace`][go.string.Replace]. `{{ replace .Env.PATH ":" }}`
  * `parseUrl $url` - Parses a URL into it's [protocol, scheme, host, etc. parts][go.url.URL]. Alias for [`url.Parse`][go.url.Parse]
  * `atoi $value` - Parses a string $value into an int. `{{ if (gt (atoi .Env.NUM_THREADS) 1) }}`
  * `add $arg1 $arg` - Performs integer addition. `{{ add (atoi .Env.SHARD_NUM) -1 }}`

## License

MIT


[go.string.Split]: https://golang.org/pkg/strings/#Split
[go.string.Replace]: https://golang.org/pkg/strings/#Replace
[go.url.Parse]: https://golang.org/pkg/net/url/#Parse
[go.url.URL]: https://golang.org/pkg/net/url/#URL


