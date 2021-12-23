# Networker Load Balancer Demo

## From project root

Start three TLS servers on ports 3001,3002,3003.

    make servers

Start a loadbalancer on port 80.

    make balance

Refreshing the page should render a changing message
showing what backend server was hit.
