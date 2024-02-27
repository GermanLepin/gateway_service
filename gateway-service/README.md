# Information about the gateway-service

The gateway service is an entry point for the application. All requests from a client go through the gateway service, and the gateway service decides where it should send a request.

What I already did:
1. user creation handler `v1/user/create`
2. user logging handler `v1/user/login`
3. user deletion handler `v1/user/delete/{uuid}`
4. user fetching handler `v1/user/fetch/{uuid}`
5. migrations
