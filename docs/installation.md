# Installation

AGEN supports macOS, Linux, and Windows. Choose the method that works best for you.

## macOS / Linux (Homebrew)

If you use Homebrew, you can install from our tap:

```bash
brew install eshanized/tap/agen
```

## Windows (Scoop)

For Windows users using Scoop:

```bash
scoop bucket add eshanized https://github.com/eshanized/scoop-bucket
scoop install agen
```

## Arch Linux (AUR)

Arch Linux users can install from the AUR:

```bash
yay -S agen-bin
```

## Debian / Ubuntu

Download the latest `.deb` package from [GitHub Releases](https://github.com/eshanized/agen/releases) and install:

```bash
sudo dpkg -i agen_*.deb
```

## Binary Download

You can download pre-compiled binaries for your platform from the [Releases Page](https://github.com/eshanized/agen/releases).

1. Download the archive for your OS/Arch.
2. Extract the `agen` binary.
3. Move it to a directory in your `$PATH` (e.g., `/usr/local/bin`).

## Building from Source

Requirements: Go 1.22+

```bash
git clone https://github.com/eshanized/agen.git
cd agen
go build -o agen ./cmd/agen
```
