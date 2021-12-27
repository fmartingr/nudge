# nudge

Simple service that will run periodic ping checks on provided IPs mostly to check if internet or an specified network is available.

## Configuration

Using environment variables:

- `NUDGE_LOG_LEVEL` (default: `warn`): Log level to use. See [logrus.Level](https://godocs.io/github.com/sirupsen/logrus#Level).
- `NUDGE_PORT` (default: `2000`): The port to serve the HTTP server
- `NUDGE_IPS` (default: `1.1.1.1 9.9.9.9`): **Space** [^1] separated list of IPs to check for connectivity.
- `NUDGE_INTERVAL` (default `60`, seconds): Interval to perform checks on the specified IPs.

## API

- `GET /health`: Simple endpoint to check the service is working
    - `200`: Service working
    - Anything else: problems!
- `GET /status`: Check the status of the connection
    - `200`: All ips answered to ICPM
    - `204`: All ips failed to answer to ICMP

[^1]: https://github.com/spf13/viper/issues/380
