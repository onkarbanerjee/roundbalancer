# RoundBalancer - HTTP(?) Round Robin API

## Overview

RoundBalancer is an HTTP-based Round Robin API that acts as a load balancer for routing HTTP POST requests to multiple instances of an Application API. It ensures that incoming requests are distributed evenly among available instances using a round-robin algorithm.

## Features

- **Application API:** A simple API that accepts HTTP POST requests with a JSON payload and responds with the same payload.
- **Round Robin API:** A routing API that distributes requests to multiple Application API instances in a round-robin fashion.
- **Configurable Instances:** Ability to configure the Round Robin API with multiple Application API instances.
- **Error Handling:** Handles failures and slow responses from Application API instances.
- **Scalable:** Supports running multiple instances of the Application API.

## Technologies Used

- **Programming Language:** [Specify your language, e.g., Python, Node.js, Go]
- **Framework:** [Specify your framework, e.g., Flask, Express, FastAPI]
- **HTTP Client:** [Specify the library used for making requests, e.g., requests, axios, fetch]

## Installation & Setup

### 1. Clone the Repository

```sh
 git clone https://github.com/your-repo/roundbalancer.git
 cd roundbalancer
```

### 2. Install Dependencies

```sh
 # If using Python:
 pip install -r requirements.txt

 # If using Node.js:
 npm install
```

### 3. Run Application API Instances

```sh
 # Example: Running three instances on different ports
 python app_api.py --port 5001
 python app_api.py --port 5002
 python app_api.py --port 5003
```

### 4. Run the Round Robin API

```sh
 python round_robin_api.py --instances http://localhost:5001 http://localhost:5002 http://localhost:5003
```

## Usage

### Sending a Request

Send a POST request to the Round Robin API:

```sh
 curl -X POST http://localhost:8000 -H "Content-Type: application/json" -d '{"game":"Mobile Legends", "gamerID":"GYUTDTE", "points":20}'
```

### Expected Response

```json
{
  "game": "Mobile Legends",
  "gamerID": "GYUTDTE",
  "points": 20
}
```

## Round Robin Logic

- The Round Robin API maintains a list of Application API instances.
- Each incoming request is forwarded to the next available instance in the list.
- If an instance is down, the request is skipped to the next available instance.
- If an instance is slow, timeouts are handled appropriately to avoid delays.

## Handling Failures

- If an instance becomes unresponsive, it is temporarily removed from the rotation.
- If an instance slows down, a timeout mechanism prevents blocking subsequent requests.

## Testing

To test the functionality, you can:

- Use **Postman** or **cURL** to send HTTP POST requests.
- Simulate failures by stopping one instance and observing request distribution.
- Use a load-testing tool like **Apache JMeter** or **locust**.

## Future Enhancements

- Implement a health-check mechanism for detecting failed instances.
- Add support for weighted round-robin for instance prioritization.
- Integrate logging and monitoring.

## License

This project is for recruitment purposes only and remains the property of **Coda Payments Pte. Ltd.**

---

Let me know if you'd like to modify anything or add more details! 🚀

give the contents that i can copy paste



