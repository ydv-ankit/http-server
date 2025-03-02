package main

import (
	"net"
	"strconv"
)

var statusText = map[int]string{
	200: "OK",
	201: "Created",
	202: "Accepted",
	204: "No Content",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	307: "Temporary Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	410: "Gone",
	500: "Internal Server Error",
	501: "Not Implemented",
	503: "Service Unavailable",
}

func WriteResponse(c *net.Conn, status int, body string) {
	response := "HTTP/1.1 " + strconv.Itoa(status) + " " + statusText[status] + "\r\n"
	response += "Content-Type: text/html\r\n"
	response += "Content-Length: " + strconv.Itoa(len(body)) + "\r\n"
	response += "\r\n"
	response += body
	(*c).Write([]byte(response))
}
