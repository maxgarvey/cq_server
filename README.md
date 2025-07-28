# cq_server
CQ Server is short for Concurrent Queue Server

This application is an HTTP server to accept requests which can be enqueued to handlers to do work. Once a request is issued, a record is created for it in Redis to track progress, which is updated by the worker asynchronously. Then the message is enqueued to RabbitMQ, where a queue worker will process it. Progress can be queried by making requests to the service which will check the record in Redis and return a status update.

## endpoints

### ask
`/ask/{requestType}`
	- Creates a new identifier token for this request.
	- Creates a record to track the status of this work and uploads to Redis.
	- Enqueues the work to RabbitMQ based on the `requestType` for process by one of the queue workers.

### get
`/get/{requestType}/{id}`
	- Retrieves record from redis for message with matching `id` and `requestType`
	- Returns response

### update
`/update/{requestType}/{id}`
	- Updates the status based on the "status" parameter in the request body in Redis
	- Called from worker context

## workers
Workers are defined in the `worker` directory. For implementation, copy the `worker.go` file, and fill in the `TODO` section to perform the work desired.