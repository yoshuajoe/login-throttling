# Login Throttling

## Definition
Loging throtling is consisted of two different words, "login" and "throttling". Throttling is something happens related to pace.
It's all the matter of User Experience (UX). For example when a user is about signing in, we will provide a disabled button after a request has been made,
this will prevent user to click the button repetitively as "throttle" in horce race, they whip the horse repetitively in small delay time to make the horse runs faster.
Thus, we will implement login throttle to prevent user from abuse the request.

## Request API
Here we have an endpoint named `/login` to make login request. We will use maximum 10 similar request in a minute, if the number goes more that 10, server will reject the request and user will have to wait 30 seconds to cool down.

## Similar request definition
We will capture the request's IP Address as key.
```golang
IPAddress := c.Request().Header.Get("X-Real-Ip")
if IPAddress == "" {
    IPAddress = c.Request().Header.Get("X-Forwarded-For")
}
if IPAddress == "" {
    IPAddress = c.Request().RemoteAddr
}
```
So in other words, if too many requests come from the same IP address, we will consider it as abusive and will reject.

## Implementation
We will use Redis to memorize total request made, and using IP Address + Minute of first request as key

## Deployment
Docker compose is using in this project, so if you had already installed it on your local machine, just type
```
docker compose up
```
then voilla, your app's running now.