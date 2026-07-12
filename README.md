# bitc-go

Image compression tool with a TUI. Compresses images using oxipng, pngquant, cwebp, and avifenc.

## Requirements

```bash
brew install imagemagick oxipng pngquant webp libavif
```

## Install

```bash
go install github.com/yazdanctx/bitc-go/cmd/compressor@latest
```

Or build from source:

```bash
git clone https://github.com/yazdanctx/bitc-go.git
cd bitc-go
make build
make install
```

## Usage

```bash
compressor ./my-images
compressor ./my-images --output ~/my-folder
compressor --version
```

## License

MIT
