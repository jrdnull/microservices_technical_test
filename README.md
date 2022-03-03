# News Service

## Design

The two requested services are Go modules in the `user_service` and `article_service`
directories.

The services each have a separate database with only identifiers being shared.

### User Service

Provides a gRPC interface to user tag preferences persisted in PostgreSQL.

### Article Service

Provides access to articles persisted in PostgreSQL optionally filtering by tags. 

### Module Layout

Inspired by a style suggested by Ben Johnson I read years ago (a more upto date look at
it by him is available here: https://www.gobeyond.dev/wtf-dial/).

#### Services

Our business logic lives here, working with data from Repositories via interfaces.

#### Repositories

Our persistence layer.

#### Presentation (grpc/http)

Making use of Services, should only be concerned with the presentation layer they're 
implementing.

### Libraries Used

My library choice has been for ease of getting going or because wanting to try 
a library out (like uptrace/bun, although I regret this).

- DATA-DOG/go-txdb - easily isolate tests by running each in a transaction
- beme/abide - http snapshot testing, allows quickly writing cases and confirming/updating desired output
- caarlos0/env - easily marshal env vars into struct
- labstack/echo - http framework 
- uptrace/bun - SQL ORM, while the fixtures saved me a lot of time I found the documentation
  to be lacking and ran into difficult to resolve issues with it
- google/go-cmp - diff lib allowing for easy assertions in tests with readable output
- go.opentelemetry.io/otel - OpenTelemetry implementation

## TODO

- [ ] Proper authentication, keys `user1` and `user2` are hardcoded for now
- [ ] Secure + encrypt services with TLS (plus mutual TLS authentication for gRPC)

## Running Services

Download dependencies for both services:
```shell
make vendor
```

Everything required to run is provided by the docker-compose file, bring it up:
```shell
docker-compose up -d
```

Once running initialise the databases and insert some test data:
```shell
DB_HOST=localhost make db_reset
```

You should now be able to test the two requested endpoints with:

http://localhost:8080/v1/articles/feed?auth_key=user1

and

http://localhost:8080/v1/articles?auth_key=user1&all_tags=1,2

## Testing

All tests can be run as usual for go tests, or `make test` is provided to run for
both of the modules.

Tests requiring the database are behind the `integration` build tag and require the
docker-compose environment to be running.

## Tracing

Tracing is enabled and exporting to Jaeger, you can view traces here:
http://localhost:16686/search

-------------------------

# The Vision

A startup company has a vision for a service that provides a tailored news feed to its users.

News will be acquired from a multitude of online sources. Each news article will be analyzed (using AI) and tagged with one or more keywords, 
before being stored in a database within the company's infrastructure. Users of the service are able to specify tags for news areas that interest them. When
a user opens their dashboard, they should see a feed of the most recent and relevant news articles.

# The Mission

There are few architectural guidelines so it’s completely up to you; however, to support their growth, the company have opted for
the microservice route. Obviously, a number of services are going to be 
required, but don't panic, we're not expecting you to implement _everything_ (although that would be impressive)!

Your task is to implement two services as follows:

- A `user` service that will provide a user's chosen tags
- A `news article` service that will provide a feed of news articles

## 1. User Service

The `User` service is a microservice that stores each user's tag selection.

### Internal Endpoint

_Protocol, method, signature, etc. - (TBD, by candidate!)_

**Role:**

This endpoint is called by the `News Article` service.

**Behaviour:**

For a given user, return all the tags specified.


## 2. News Article Service

The `News Article` service is a microservice that stores news articles.

Each article must contain a `Title`, `Timestamp` and list of `Tags`.

This service also provides an external endpoint that allows users to retrieve articles, filtered and sorted by their timestamp.

### Public Endpoint

_Protocol, method, signature, etc. - (TBD, by candidate!)_

**Role:**

This endpoint is called from user's browser/phone/tablet.

**Predicate**

> An article is included in the response if it has at least 1 tag matching those of the user.

**Behaviour:**

Return news articles, filtered by the tags of the current user.

### Public Endpoint

_Protocol, method, signature, etc. - (TBD, by candidate!)_

Request must include 2 tags

**Role:**

This endpoint is called from user's browser/phone/tablet.

**Predicate**

> An article is included in the response if its tags contain _both_ of the tags specified in the request.

The product team mentioned that they might need to change these predicate values. Modifying both the number 
of tags in the request and that used in the matching criteria. So, bonus points if you make these configurable! ;)


**Behaviour:**

Return news articles, filtered by the tags included in the request.

## Prerequisites
- Handle all failure cases
- Your code should be tested
- Provide a `docker-compose.yaml` file for any third party services that you use 
- Provide a clear explanation of your approach and design choices (while submitting your pull request)
- Provide a proper `README.md`:
    - Explain how to setup and run your code
    - Include all information you consider useful for a seamless coworker on-boarding

## Workflow
- Create a new branch
- Commit and push to this branch
- Submit a pull request once you have finished

We will then write a review for your pull request!

## Bonus

- Add metrics / request tracing / authentication 📈
- Add whatever you think is necessary to make the app awesome ✨
