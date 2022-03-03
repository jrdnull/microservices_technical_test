# User Service

Provides access to users, and their settings in the news feed system.

## Internal API

The internal API is provided via a gRPC server:

```protobuf
service UserService {
  rpc GetNewsFeedTags(GetNewsFeedTagsRequest) returns (GetNewsFeedTagsResponse) {}
}

message GetNewsFeedTagsRequest {
  int64 id = 1;
}

message GetNewsFeedTagsResponse {
  repeated int64 tag_ids = 1;
}
```

## Internal Endpoint

**rpc GetNewsFeedTags(GetNewsFeedTagsRequest) returns (GetNewsFeedTagsResponse)**

**Role:**

This endpoint is called by the `News Article` service.

**Behaviour:**

For a given user, return all the tags specified.