# gass

A dead-simple, poorly named [Sass](https://sass-lang.com) transpiler created with the intention of making the simultaneous transpilation of multiple Sass files easy.

`gass` accepts a list of sources in form of a "Gassfile." Which is really just a newline delimited file of input/output paths. See [#usage](#usage) for more info.

# Installing

This is honestly just something I wrote for myself, but if you find yourself needing similar functionality, you can install `gass` with the Go toolchain.

```sh
$ go install github.com/holedaemon/gass@latest
```

# Usage

`gass` is configured in two ways:

1. with a "gassfile"
2. with command-line flags

## Gassfile

A Gassfile is really just a newline delimited, plaintext file of paths to Sass/CSS files. Consider the following:

```
/home/max/git/github.com/holedaemon/website/style/main.scss /home/max/git/github.com/holedaemon/website/static/css/main.css
```

The first path presented is the **input**: `gass` will transpile this into CSS.  
The second path is the **output**: this is where `gass` will output the transpiled input. You can either specify a file or a directory. When using the latter, the output file will have the same name as the input.

The Gassfile also "supports" comments; if a line starts with a '#', `gass` will skip that line.

## Flags

`gass` uses Go's `flag` package to configure both the transpiler and arguments passed to it upon execution. The following flags are provided:

### Transpiler

The following are used to configure the transpiler itself.

```
-h    Output help.
-f    Path to the Gassfile to use. Defaults to ".gassfile".
-d    Starts the transpiler in debug mode. Defaults to false.
-b    Path to the dart-sass binary to use as a backend. Default is blank as the transpiler will attempt to infer.
-t    The maximum alloted time to wait for transpilation. Defaults to 30s.
```

### Arguments

The following are used to configure transpilation.

```
-c    Configures the transpiler to output minified CSS. Defaults to false.
-s    Configures the transpiler what syntax to expect input files to use. Defaults to SCSS. Valid options are SCSS, Sass, and CSS.
-m    Configures the transpiler to emit source maps when transpiling. Defaults to true.
-e    Configures the transpiler to embed sources into source maps, if they are enabled. Defaults to false.
```

For now, these are the only way to configure transpilation settings, and they apply to all jobs. In the future, I'd like to also optionally support setting these per-job in a Gassfile.

# How it works

`gass` is very simple. Realistically, a script could be made to do the same exact thing in a few minutes.

There aren't any maintained Go packages for Sass. That being the case, we actually turn to the [reference implementation](https://github.com/sass/dart-sass) of Sass for help. `dart-sass` exposes an IPC interface for consumers to use via the `--embedded` flag. `gass` leverages this by using [bep's godartsass package](https://github.com/bep/godartsass) (the one made for [Hugo's](https://gohugo.io) Sass support) to interface with the process. As inputs are read from a Gassfile, they are sent to `dart-sass` via stdin, transpiled to CSS, and returned along with any source maps (if enabled). From there, we just write the CSS and maps to their respective output paths and call it good. Quick n' dirty.

# License

See [LICENSE](LICENSE).
