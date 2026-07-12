# bitc-go

Image compression tool. Compresses images to AVIF using avifenc, with optional black & white conversion via ImageMagick.

## Requirements

```bash
brew install imagemagick libavif
```

## Install

```bash
go install github.com/yazdanctx/bitc-go/cmd/bitc@latest
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
bitc <directory> [flags]
```

### Examples

```bash
# Compress all images in a directory to AVIF
bitc ./my-images

# Convert to black & white (2-color) then compress
bitc --bw ./my-images

# Specify output directory
bitc ./my-images --output ~/my-folder
```

### Flags

| Flag | Description |
|------|-------------|
| `--bw` | Convert images to 2-color black & white before compressing |
| `--output` | Override output directory (defaults to `~/Downloads/compressed-<timestamp>`) |
| `--version` | Print version |

## License

MIT
