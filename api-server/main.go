package main

func main() {
	addr := ":8080"
	server := NewServer(addr)
	server.Run()
}
