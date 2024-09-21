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

