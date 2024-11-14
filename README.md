app=> main.go
 Flow Overview:
==============
Configuration Loading: Load configuration based on the environment (Prod/UAT).
Logging Initialization: Set up logging based on the configuration.
TCP Connection: Connect to the primary gateway, and fallback to secondary if needed.
Session and User Login: Authenticate the session and user after connecting to the gateway.
Subscribe and Capture Data: Subscribe to data feeds and start capturing the data.

Key Functions:
==============
messages.ConnectToGateway: Sends a connection request and handles responses from the gateway.
messages.SessionLogin: Logs in to the session after establishing the gateway connection.
messages.UserLogin: Authenticates the user.
messages.Subscribe: Subscribes to specific feeds or data streams.
messages.StartFeedCapturing: Begins capturing the data from the feed.

messages=> connections.go
=========================
Overall Flow:
The getConnectionRequest() creates the connection request.
ConnectToGateway() sends the request and waits for a response.
The response is parsed depending on whether it’s accepted or rejected.
The parsed gateway IPs are returned as primary and secondary connections.

messages=> rejection.go
=======================
What it does: This function reads a rejected packet from a network connection, extracts the header and body, and retrieves the rejection reason (a text message) from the body.
Key Components:
Logging: Logs details about the rejected message.
utils.Recv(): Reads the body from the network connection.
Text extraction: Extracts the rejection reason based on the length provided in the packet body.

messages=> retranmission.go
===========================
getRetransmissionRequest constructs and returns a retransmission request packet.
Retransmit handles the entire retransmission flow, sending the request, receiving the response, and handling errors and rejections.
parseRetransmissionResponsePkt parses the retransmission response packet and returns the parsed data.

messages=> session.go
=====================
Key Concepts:
Session Login Request: This is a request that includes the session ID, password, and other metadata necessary for a session to log in.
Networking: The connection to the server is represented by net.Conn, and data is sent/received using utils.Send() and utils.Recv().
Error Handling: Logs errors at various steps and returns appropriate errors if something goes wrong (e.g., packet rejection, parsing errors).
Logging: All events, errors, and data are logged using the zap logging system for debugging and monitoring.
Error Cases:
FailedToSend: Logs and returns an error if sending the request fails.
FailedToReadHeader: Logs and returns an error if reading the response header fails.
Rejected Packet: Logs details if the request was rejected and parses the rejection reason.
Unknown TemplateId: Logs and returns an error if the response has an unexpected template ID.

messages=> subscription.go
==========================
Key Concepts:
Message Headers: The request and response packets are structured with headers (MHeader and RHeader), defining the message length, sequence numbers, and template IDs.
Template IDs: These IDs uniquely identify the type of message being sent or received (e.g., subscription request or response, rejection).
Networking: The server communication is done via net.Conn, a network connection object. Data is sent/received using utils.Send() and utils.Recv().
Error Handling: The code ensures that errors during sending, receiving, or parsing are logged and returned for further action.
Logging: All key operations, including request creation, sending, receiving, and errors, are logged using the zap logging system.
Error Handling:
FailedToSend: If the subscription request fails to be sent, an error is logged, and the function returns with an error.
FailedToRecvHeader: If the header of the response cannot be read, it logs the error and returns.
RejectedResponse: If the response is rejected, it parses the rejection reason, logs it, and returns an error.
Unknown Packet: If an unrecognized template ID is received, it logs the unknown packet and returns an error.
Conclusion:
This code is designed to handle the full lifecycle of a subscription request from creating the message to sending it, receiving a response, and handling either a successful or rejected result. Logging and error handling ensure that issues during the process are captured for debugging purposes.

messages=> user.go
==================
Flow Overview:
This code defines functions related to user login within a system, likely for some network communication using the Go language. The code involves constructing and sending a user login request, handling the response, and logging information.
Create Login Request: The getUserLoginRequest function creates a structured request packet based on the expected message format.
Send Request: The login request is sent over the established TCP connection.
Receive Response Header: The system waits for the response header and processes it to identify whether the login was successful or rejected.
Handle Response:
If successful, it parses the response body.
If rejected, it handles the rejection packet.
If the response is unknown, it raises an error.
Logging: Throughout the process, detailed logging helps track each step for debugging and monitoring.

==============================================================================================================================================================================================================================
































