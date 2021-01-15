# fwd

Forwards local port(s) to remote host(s). Like socat, but with better syntax.

    $ fwd google.com:80
    2021/01/15 18:48:13 127.0.0.1:80 -> google.com:80

    $ fwd 8080:google.com:80
    2021/01/15 18:48:13 127.0.0.1:8080 -> google.com:80

    $ fwd 0.0.0.0:8080:google.com:80
    2021/01/15 18:48:13 0.0.0.0:8080 -> google.com:80

    $ fwd 1234:google.com:80 5678:microsoft.com:80
    2021/01/15 18:49:57 127.0.0.1:1234 -> google.com:80
    2021/01/15 18:49:57 127.0.0.1:5678 -> microsoft.com:80
