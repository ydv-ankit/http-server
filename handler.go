package main

func handleRequest(uri string) string {
	switch uri {
	case "/":
		return "Welcome to the homepage"
	case "/hello":
		return "Hello there"
	case "/bye":
		return "Goodbye"
	default:
		return "Not Found"
	}
}
