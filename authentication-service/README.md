# Information about the authentication-service

The authentication service is a service that accepts all requests from the gateway service. All requests from a client go through the gateway service, and the gateway service sends them to the authentication service to authorize and authenticate users.

What I already did:
1. user logging handler `user/login`
2. user validation token handler `user/validate-token`
3. migrations
