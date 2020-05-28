package cmd

import "fmt"

var (
	subExampleFormat = "\n%s:\n\n\tlong form:\n\n\t\t%s\n\n\tshort form:\n\n\t\t%s\n"

	backDoorExample = newExample("backdoor", []subExample{
		subExample{
			description: "Create a new backdoor",
			longForm:    "networker backdoor --create --port <port>",
			shortForm:   "networker bd --create -p <port>",
		},
		subExample{
			description: "Connect to an existing backdoor",
			longForm:    "networker backdoor --connect --address <host>:<port>",
			shortForm:   "networker bd --connect -a <host>:<port>",
		},
	})

	captureExample = newExample("capture", []subExample{
		subExample{
			description: "Capture packets on en1 for 10 seconds or until 100 packets have been captured and log the capture status during capture",
			longForm:    "networker capture --devices en1 --seconds 10 --out myCaptureSession --limit --num 100 --verbose",
			shortForm:   "networker c -d en1 -s 10 -out myCaptureSession -l -n 100 -v",
		},
		subExample{
			description: "Don't specify an outfile and instead print captured packets from the en0 interface to stdout",
			longForm:    "networker capture --device en0 --seconds 10 --limit --num 100 --verbose",
			shortForm:   "networker c -d en0 -s 10 -l -n 100 -v",
		},
	})

	listExample = newExample("list", []subExample{
		subExample{
			description: "List the IP of the current network gateway, local IP of this machine, and remote IP of this machine",
			longForm:    "networker list --me",
			shortForm:   "networker ls -m",
		},
		subExample{
			description: "List the hostname, IP address, and connection status of all devices on the current network",
			longForm:    "sudo networker list --all",
			shortForm:   "sudo networker ls -a",
		},
	})

	lookUpExample = newExample("lookup", []subExample{
		subExample{
			description: "Look up the network for a given hostname or IP",
			longForm:    "networker lookup --network 31.13.65.36",
			shortForm:   "networker lu -n 31.13.65.36",
		},
		subExample{
			description: "Look up the hostname for a given IP",
			longForm:    "networker lookup --hostnames 157.240.195.35",
			shortForm:   "no short form as -h is reserved for help",
		},
		subExample{
			description: "Look up nameservers for a given hostname",
			longForm:    "networker lookup --nameservers youtube.com",
			shortForm:   "networker lu -s youtube.com",
		},
		subExample{
			description: "Look up the addresses for a given hostname",
			longForm:    "networker lookup --addresses youtube.com",
			shortForm:   "networker lu -a youtube.com",
		},
	})

	proxyExample = newExample("proxy", []subExample{
		subExample{
			description: "Start a new proxy server that listens on a given port and forwards traffic to a given address",
			longForm:    "networker proxy --listen-on <port> --upstream <host>:<port>",
			shortForm:   "networker p -l <port> -u <host>:<port>",
		},
	})

	scanExample = newExample("scan", []subExample{
		subExample{
			description: "Scan a comma-separated list of TCP ports of an address and only log out the ones that are open",
			longForm:    "networker scan --ip <address> --ports 22,80,3389 --tcp-only --open-only",
			shortForm:   "networker s --ip <address> -p 22,80,3389 -t -o",
		},
		subExample{
			description: "Scan all TCP ports up to port 1024 and only log out the ones that are open",
			longForm:    "networker scan --ip <someIPaddress> --up-to 1024 --tcp-only --open-only",
			shortForm:   "networker s --ip <address> -u 1024 -t -o",
		},
	})

	requestExample = newExample("request", []subExample{
		subExample{
			description: "Explicitly passing the file path of a JSON or XML file to add the file contents to the body of a POST request",
			longForm:    "networker request --url https://api.thecatapi.com/v1/votes --method POST --file scrap.json --add-headers <key>:<value>,<key>:<value>",
			shortForm:   "networker r -u https://api.thecatapi.com/v1/votes -m POST -f scrap.json -a <key>:<value>,<key>:<value>",
		},
		subExample{
			description: "Send a Delete request. Supported methods include GET, POST, PATCH, PUT, and DELETE",
			longForm:    "networker request --url https://api.thecatapi.com/v1/votes/<voteID> --method DELETE --add-headers x-api-key:<api-key>",
			shortForm:   "networker r -u https://api.thecatapi.com/v1/votes/<voteID> -m DELETE -a x-api-key:<api-key>",
		},
		subExample{
			description: "Networker will set the protocol scheme (defaults to https://) and method (defaults to GET) if not set",
			longForm:    "networker request --url api.thecatapi.com/v1/votes --add-headers x-api-key:<api-key>",
			shortForm:   "networker r -u api.thecatapi.com/v1/votes -a x-api-key:<api-key>",
		},
	})
)

type (
	example struct {
		cmdName     string
		subExamples []subExample
	}
	subExample struct {
		description string
		longForm    string
		shortForm   string
	}
)

func newExample(name string, subExamples []subExample) string {
	e := example{
		cmdName:     name,
		subExamples: subExamples,
	}
	return e.format()
}

func (e *example) format() string {
	var formattedExample string
	for _, se := range e.subExamples {
		formattedExample += fmt.Sprintf(subExampleFormat,
			se.description,
			se.longForm,
			se.shortForm,
		)
	}
	return formattedExample
}
