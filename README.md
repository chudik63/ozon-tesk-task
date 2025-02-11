# Tesk Task for OZON Bank
## Installation
```
git clone https://github.com/chudik63/ozon-tesk-task.git
cd ozon-tesk-task
```
## Building
### In-memory storage
```
make run-in-memory
```
or
```
docker compose -f docker-compose.memory.yaml up
```
### PostgreSQL
```
make run-postgres
```
or
```
docker compose -f docker-compose.yaml up --build
```

## Usage
Go to localhost:8080

### Queries
```graphql
    query ListPosts {
        posts(page: 1, limit: 10) {
            id
            title
            content
            author
            createdAt
            comments {
                id
                content
                author
            }
        }
    }

    query DeletePost {
        deletePost(postId: 1)
    }

    query DeleteComment {
        deleteComment(commentId: 1)
    }

    query GetComments {
        comments(postId: 1, page: 1, limit: 10) {
            id
            replies {
                id
                replies {
                    id
                }
            }
        }
    }

    query GetPost {
        post(id: 1) {
            id
            title
            content
            createdAt
            author
            comments {
                id
                createdAt
                content
                replies {
                    id
                }
            }
        }
    }

```
### Mutations
```graphql
mutation CreatePost($CreatePostInput: CreatePostInput!) {
  createPost(input: $CreatePostInput) {
    id
    title
    content
    allowComments
  }
}

mutation CreateComment($CreateCommentInput: CreateCommentInput!) {
  createComment(input: $CreateCommentInput) {
    id
    postId
    parentId
    author
    content
  }
}

mutation ReplyToComment($ReplyToCommentInput: CreateCommentInput!) {
  createComment(input: $ReplyToCommentInput) {
    id
    postId
    parentId
    author
    content
  }
}
```
Variables:
```graphql
{
  "CreatePostInput": {
    "title": "first post",
    "content": "test test test content",
    "allowComments": true
  },
  "CreateCommentInput": {
    "postId": 2,
    "content": "NEW COMMENT"
  },
  "ReplyToCommentInput": {
    "postId": 2,
    "parentId": 8,
    "content": "reply"
  }
}
```
### Subscription
```graphql
subscription NewCommentAdded {
  commentAdded(postId: 1) {
    id
    postId
    content
    parentId
    createdAt
  }
}
```
