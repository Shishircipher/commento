curl -X POST http://localhost:8000/comments -H "Content-Type: application/json" -d '{
      "author": {
        "name": "John Doe",
        "email": "johndoe@example.com",
        "picture": "https://example.com/avatar.jpg"
      },
      "post_id": 1,
      "parent_id": null,
      "content": "This is a top-level comment"
    }'
