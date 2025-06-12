# Gale
Gale is a lightweight HTTP server load tester with beautiful report generation. It works through the command line for instant testing of your servers.


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

