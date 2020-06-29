# Request

## Usage

    Send an HTTP request.

    Flags:
    -a, --add-headers strings   Add a list of comma-separated request headers. (format : key:value,key:value,etc...)
    -f, --file string           Path to JSON or XML file to use for request body. (content-type headers for each file-type are set automatically)
    -h, --help                  help for request
    -m, --method string         Specify method. (supported methods include GET, POST, PUT, PATCH, and DELETE) (default "GET")
    -t, --time-out int          Specify number of seconds for time-out. (default 3)
    -u, --url string            URL to send request.

## Examples

When sending a `GET` request, the result of not explicitly specifying a request method defaults to a `GET`.

    networker request --url https://api.thecatapi.com/v1/breeds


When sending a `POST` request, the request body can be provisioned with JSON from a local file. The request headers can be specified as an unterminated list of comma-separated key/value pairs.

    networker request --url https://api.thecatapi.com/v1/votes --method POST --file scrap.json --add-headers <key>:<value>,<key>:<value>


