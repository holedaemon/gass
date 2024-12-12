# gass

A dead-simple, poorly named [Sass](https://sass-lang.com) transpiler created with the intention of making the simultaneous transpilation of multiple Sass files easy.

`gass` accepts a list of sources in form of a "Gassfile." Which is really just a newline delimited file of input/output paths. See [#usage](#usage) for more info.

# Installing

This is honestly just something I wrote for myself, but if you find yourself needing similar functionality, you can install `gass` with the Go toolchain.

```sh
$ go install github.com/holedaemon/gass@latest
```

# Usage

After installing, one may invoke the tool by calling `gass` in a shell

# How it works

`gass` is very simple. Realistically, a script could be made to do the same exact thing in a few minutes.

There aren't any maintained Go packages for Sass. That being the case, we actually turn to the [reference implementation](https://github.com/sass/dart-sass) of Sass for help. `dart-sass` exposes an IPC interface for consumers to use via the `--embedded` flag. `gass` leverages this by using [bep's godartsass package](https://github.com/bep/godartsass) (the one made for [Hugo's](https://gohugo.io) Sass support) to interface with the process. As inputs are read from a Gassfile, they are sent to `dart-sass` via stdin, transpiled to CSS, and returned along with any source maps (if enabled). From there, we just write the CSS and maps to their respective output paths and call it good. Quick n' dirty.

# License

See [LICENSE](LICENSE).
