## Quiz microservices with Golang
I have created five microservices: auth, users, quiz, answers and solution. 
I have chosen nosql database: mongodb to store data. I apply synchronous communication between microservices, 
they communicate with each other through rest interface.

## Usage
Command to run application
```bash
docker-compose up --build -d
docker-compose logs -f name_service
```

## License

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)  
