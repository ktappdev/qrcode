package urlhandler

import "os"

func GetEnvItems() (port string, server string) {
	port = os.Getenv("PORT")
	server = os.Getenv("SERVER")
	return port, server
}
