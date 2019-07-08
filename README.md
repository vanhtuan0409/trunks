# Trunks

<p align="center">
  <img src="./resources/trunks.png"
</p>

Trunks is a lightweight, template-based load testing tool built upon [vegeta](https://github.com/tsenart/vegeta). Trunks combined the power of Vegeta and Golang Template to generate randomized data, mimics the pattern of real-live traffic.

If you like my work, consider buy me a coffee :D

<a href="https://www.buymeacoffee.com/sHZbgvYh0" target="_blank"><img src="https://bmc-cdn.nyc3.digitaloceanspaces.com/BMC-button-images/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;" ></a>

### Installation

#### Pre-built binary

Please checkout [Release page](https://github.com/vanhtuan0409/trunks/releases)

#### Source

You need `go` installed and `GOBIN` in your `PATH`. Once that is done, run the command:

```shell
$ go get -u github.com/vanhtuan0409/trunks
```

### Example

```shell
$ trunks -rate 10 -duration 20 -target targets.sample.yml | vegeta report
```

### Usage manual

```console
Usage of trunks:
  -debug string
        Write debug log to file (default discard)
  -duration int
        Duration to run the request (in seconds)
  -output string
        Output file (default "stdout")
  -rate int
        Request per second to send (default 5)
  -target string
        Targets config file path (default "targets.yml")
```

#### `-debug`

Specifies the path of debug log file. It defauls to discarded `/dev/null`

#### `-duration`

Specifies the amount of time to issue request to the targets. The internal concurrency structure's setup has this value as a variable. The actual run time of the test can be longer than specified due to the responses delay. Use 0 for an infinite attack.

#### `-output`

Specifies the output file to which the binary results will be written to. Made to be piped to the report command input. Defaults to stdout.

For more details about load test analytics, please refer to [vegeta attack](https://github.com/tsenart/vegeta#report-command)

#### `-rate`

Specifies the request rate per time unit to issue against the targets. The actual request rate can vary slightly due to things like garbage collection, but overall it should stay very close to the specified. If no time unit is provided, 1s is used.

#### `-target`

Specifies the file config for targets. Refer to [Targets format](#targets-format)

### Targets format

```
meta:
  headers:
    Accept: application/json
targets:
  - url: "http://localhost:8080/api1?lat={{ randNumeric 3 }}&long={{ randNumeric 3 }}"
    method: GET
    repeat: 2
    headers:
      Authorization: "Bearer xxx"
  - url: "http://localhost:8080/api2?token={{ randAlphaNum 12 }}"
    method: POST
    repeat: 3
    body: |
      {
        "timestamp": {{ now | unixEpoch }},
      }
```

Go template functions are powered by [Sprig](http://masterminds.github.io/sprig/)

#### `meta.headers`

Specifies common Header for all targets. Not allow templating

#### `targets[].url`

Specifies target URL for load test request. Allow templating

#### `targets[].method`

Specifies HTTP Method for load test request.

#### `targets[].body`

Specifies HTTP Body for load test request. Allow templating and only taken into affect when Method is not `GET`

#### `targets[].headers`

Request-specific headers, will override [`meta.headers`](#metaheaders) if duplicated.

#### `targets[].repeat`

Number of repeated time for a request (default to 1). Used for balancing request distribution among apis
