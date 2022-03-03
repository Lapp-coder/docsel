# Docsel 
[![Go](https://img.shields.io/badge/go-1.17-blue)](https://golang.org/doc/go1.17) [![Release](https://img.shields.io/badge/release-1.0.0-success)](https://github.com/Lapp-coder/docsel/releases)

Docsel is a utility that allows you to run the services you choose based on the docker-compose file.

## Example
![docsel example](https://github.com/Lapp-coder/docsel/blob/resourses/docsel-example.gif)

## Installation 
1. Clone the repository
```shell
$ git clone https://github.com/Lapp-coder/docsel
```
2. Go to the directory of the utility 
```shell
$ cd docsel
```
3. Follow the installation steps for your OS
* MacOS on M1
```shell
$ chmod +x build/bin/docsel-darwin_arm64 \
&& make setup-macos_m1
```
* MacOS on Intel
```shell
$ chmod +x build/bin/docsel-darwin_amd64 \
&& make setup-macos_intel
```

* Linux
```shell
$ chmod +x build/bin/docsel-linux_amd64 \
&& make setup-linux
```

* Windows
```shell
$ setup-windows.bat
```
