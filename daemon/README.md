# Daemon

This project is built using [Wails](https://wails.io/), a framework for creating desktop applications using Go and modern web technologies. It provides a Go backend and a frontend built with standard web technologies.

## Prerequisites

Before you can run the application, ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.19 or higher)
- [Node.js](https://nodejs.org/) (version 14 or higher)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) (Follow the Wails installation instructions)
- [osquery](https://osquery.io/) (for system monitoring integration)


### Install Wails CLI

If you haven't installed Wails yet, you can install it by running:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### To run 

Be sure the osquery daemon is running, follow instructions to start osquery daemon from these links for your platform
https://osquery.readthedocs.io/en/stable/installation/install-windows/
https://osquery.readthedocs.io/en/stable/installation/install-macos/



```bash
In a separate termainal, run these commands first for both dev and prod builds

osqueryi --nodisable_extensions
osquery> select value from osquery_flags where name = 'extensions_socket';
+-----------------------------------+
| value                             |
+-----------------------------------+
| /Users/USERNAME/.osquery/shell.em |
+-----------------------------------+

Run this command in another termainal

wails dev
```
This should start dev build


### To build

```bash
wails build
```

### Profiling the application

This application includes Go's built-in profiling tool pprof to measure performance and identify bottlenecks.

## Running the Application with Profiling
The application is already set up with the necessary code to expose the pprof profiling interface.
Once you run the application, the profiling server will be available on localhost:6060.

## Collect CPU profile by running this command (profile duration is typically 30 seconds):

```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```
Once the profile is downloaded, you can analyze it with pprof:

```bash
go tool pprof cpu.prof
```

To visualize the profile, generate a graph (requires Graphviz to be installed):
```bash
go tool pprof -svg cpu.prof > cpu_profile.svg
```

## Further Profiling Options
You can also profile goroutines, threads, and blocking events:

Goroutine profile: http://localhost:6060/debug/pprof/goroutine
Thread profile: http://localhost:6060/debug/pprof/threadcreate
Blocking profile: http://localhost:6060/debug/pprof/block

Refer to Go pprof docs


