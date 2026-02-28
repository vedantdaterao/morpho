```
NAME
      torgo - minimal BitTorrent client

SYNOPSIS
      torgo -f <file.torrent> [-o <output_dir>] [-v]

DESCRIPTION
      A BitTorrent client written in Go with zero external dependencies.

OPTIONS
      -f  file.torrent
            Path to .torrent file (required)

      -o  output_dir
            Output directory for downloaded file (default: current directory)

      -v  Enable verbose logging
            Shows detailed peer connection info and piece-by-piece progress

```
## Demo
![Demo](/assets/demo.gif)

## REFERENCES
  - [BitTorrent Protocol Specification (BEP 0003)](http://www.bittorrent.org/beps/bep_0003.html)
