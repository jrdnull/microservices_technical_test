# Article Service

Provides access to articles in the news feed system.

## Public API

The public API is exposed via HTTP.

Authentication is using two hardcoded test keys in a query paramter *auth_key* 
for ease of testing via a browser.

Two users are provided in the fixture data, use `auth_key=user1` or `auth_key=user2`.

### Public Endpoint

**HTTP GET /v1/articles/feed**

**Predicate**

> An article is included in the response if it has at least 1 tag matching those of the user.

**Behaviour:**

Return news articles, filtered by the tags of the current user.

### Public Endpoint

**HTTP GET /v1/articles**

**Parameters:**

#### all_tags
Comma separated list of tag ids, the articles must have all the tag ids to be returned.
This is limited to exactly 2 tags if used to meet the required predicate.

> Request must include 2 tags

The requirement for 2 can be configured via the env var `ARTICLE_TAG_FILTER_INPUTS`.

**Predicate**

> An article is included in the response if its tags contain _both_ of the tags specified in the request.

> The product team mentioned that they might need to change these predicate values. Modifying both the number
> of tags in the request and that used in the matching criteria. So, bonus points if you make these configurable! ;)

**Bonus Points:**

Provided to this endpoint is also the `any_tags` query parameter which will return any articles
having at least one of the ids provided. Aside from being configurable via the env var the `all_tags`
filter is implemented in such a way that you can just remove the validation from the HTTP handler to
allow any amount of tags to be passed to it.

The option of requesting with no filters to return all articles paginated is also avaialble.

**Behaviour:**

Return news articles, filtered by the tags included in the request.
