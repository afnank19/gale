# Gale
Gale is a lightweight HTTP server load tester with beautiful report generation. It works through the command line for instant testing of your servers.

## Installation

### Linux

1. Download the binary and the install.sh script from the releases tab.
2. Run `./install.sh` in the dir where both the binary and script are.

### Windows

1. Download the windows release
2. Move to a place of your choice
3. Edit environment variables
4. Add a new PATH to the binary you just placed

## Usage

**Flag Structure:** `-[flagletter]=[value]` OR `--[flagname]=[value]`

| Flag (Long / Short)       | Description                                                                 | Example                   |
|---------------------------|-----------------------------------------------------------------------------|---------------------------|
| `--threads` / `-t`        | Number of maximum threads to use. Defaults to number of physical cores.    |                           |
| `--connections` / `-c`    | Number of concurrent connections.                                           | `-c=10`                   |
| `--duration` / `-d`       | Time to run the test. Units: `s`, `m`, `h`.                                 | `-d=10s`                  |
| `--url` / `-u`            | The URL of the server. **Required**                                         | `http://localhost:3000`   |

### Notes:
- Gale is quite minimalist right now, only supporting GET requests with minimal headers. Support for more headers might get worked on, but not confirmed yet. You can contribute/fork if you require that feature.

