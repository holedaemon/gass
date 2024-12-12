# gass

A dead-simple, poorly named [Sass](https://sass-lang.com) transpiler created with the intention of making the simultaneous transpilation of multiple Sass files easy.

`gass` accepts a list of sources in form of a "Gassfile." Which is really just a newline delimited file of input/output paths. See [#usage](#usage) for more info.

# Installing

This is honestly just something I wrote for myself, but if you find yourself needing similar functionality, you can install `gass` with the Go toolchain.

```sh
$ go install github.com/holedaemon/gass@latest
```

# Usage

Once installed, `gass` will be available in your shell under the same name (assuming you have $GOBIN in your $PATH). All one needs to get started is a Gassfile:

## Gassfile

As mentioned, a Gassfile is just a newline delimited plaintext file. Each line should contain a pair of sources in form of an input file and output. e.g.

```txt
/home/max/git/github.com/holedaemon/gass/testdata/input.scss /home/max/git/github.com/holedaemon/gass/testdata/output
```

The input MUST lead to a file, while the output can lead to either a file or directory. When using the latter, the output CSS file will have the same name as the input.

> [!WARNING]
> Both the input and output MUST be an absolute path.

The Gassfile also "supports" comments. If a line starts with '#', `gass` will skip it.

## Flags

If you want to fine-tune the transpiler or its output, `gass` provides a number of configuration options via flags.

```
-v               Prints the current version and exits      | Default: false

-d               Runs gass in debug mode                   | Default: false
-f string        Set the path of your Gassfile             | Default: ".gassfile"
-b string        Set the path to your dart-sass binary     | Default: ""
-t duration      Set the transpilation timeout             | Default: 30s

-c               Tells gass to minify outputs              | Default: false
-m               Tells gass to generate source maps        | Default: true
```

# How it works

`gass` is very simple. Realistically, a script could be made to do the same exact thing in a few minutes.

There aren't any maintained Go packages for Sass. That being the case, we actually turn to the [reference implementation](https://github.com/sass/dart-sass) of Sass for help. `dart-sass` exposes an IPC interface for consumers to use via the `--embedded` flag. `gass` leverages this by using [bep's godartsass package](https://github.com/bep/godartsass) (the one made for [Hugo's](https://gohugo.io) Sass support) to interface with the process. As inputs are read from a Gassfile, they are sent to `dart-sass` via stdin, transpiled to CSS, and returned along with any source maps (if enabled). From there, we just write the CSS and maps to their respective output paths and call it good. Quick n' dirty.

# License

See [LICENSE](LICENSE).
