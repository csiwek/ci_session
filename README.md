# Codeigniter's CI Session

## Purpose

This library provides a way access and manipulate Sessions created by Codeigniter PHP framework.
It can help if you are planning to move to golang so part of the website can still operate on Codeigniter and you can gradually migrate your code to Golang.
The library also provides Middleware for Gin framework

## Supported Functions

- Has been tested with Codeigniter 3.0
- Provides Flash functions
- Provides Get/Set Userdata
- Support only sessions stored in files. If you are currently storing sessions in database you either need to switch to files or add functionality to store sessions in required storage.
- Provides Middleware for Gin
- Create session compatible with Codeigniter
- No unit tests for now. Feel free to contribute here

## Usage

- Make sure that the Go app has access to the directory where your sessions Codeigniter sessions are stored


## Examples

### Gin Middleware

See full example in the /example folder
