package server

import "log"

func main() {
	//1.получаем env
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

}
