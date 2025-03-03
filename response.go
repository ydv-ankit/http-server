package main

import (
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

func createHeader(key string, value string) string {
	return key + ": " + value + "\r\n"
}

func (c *ClientConnection) WriteTextResponse() {
	response := c.Response.PROTOCOL + " " + strconv.Itoa(c.Response.STATUS) + " " + statusText[c.Response.STATUS] + "\r\n"
	for key, value := range c.Response.HEADERS {
		response += createHeader(key, value)
	}
	response += "\r\n" + string(c.Response.BODY)
	c.Conn.Write([]byte(response))
}
