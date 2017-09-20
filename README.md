Requests Counter
================

A Golang implementation for request counter in the last 60 seconds (sliding window).
The server should continue to the return the correct numbers after restarting it, by persisting data to a file.

Installation
-------------

### Prerequisites

* [Golang][1] installation, having [$GOPATH][2] properly set.

To install [**requests-counter**](https://github.com/ahmedkamals/requests-counter)

```bash
$ go get github.com/ahmedkamals/requests-counter
```

Test Driver
-----------

You can use the following steps as a testing procedure

  * **Server**
    ```bash
    $ make build
    $ bin/requsts-counter-{OS}-{ARCH} -config "{CONFIG_PATH.json}"
    ```
    
    **`Environment`**
     * `OS` - the current operating system, e.g. (linux, darwin, ...etc.)
     * `ARCH` - the current system architecture, e.g. (386, amd64)
        
    **`Params`**:
     * `CONFIG_PATH` - a path to a json fil config e.g. `"config/main.json"`      

* **Multiple Requests**
    ```bash
    # Sending one request per second.
    $ for i in {0..1000..1}
      do 
         curl "http://localhost:8000/metrics"
         sleep 1;
      done
    ```

## Tests
    
Not all items covered, just make one example.
    
```bash
$ make unit
```

## Coding - __Structure & Design__
* `Server` - reports metrics/stats over HTTP using the path `/metrics` or `/stats`,
also used for health check using `/health` or `/status`.  
* `Dispatcher` - tracks requests count, and handles data recover in case of failure.
* `Saver` - handles saving data.

## Todo
   - More unit tests
   - Performance and memory optimization.
   - Refactoring

Enjoy!

[1]: https://golang.org/dl/
[2]: https://golang.org/doc/install
