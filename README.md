# Daterr (sticking with this till i come up with a doper name)
A web API for a dating site


---

# **Dating Site API - Specifications**
### **Tech Stack:**
- **Backend:** Go (Gin or Fiber)
- **Database:** PostgreSQL (`sqlc` for queries)
- **Authentication:** PASETO Tokens
- **File Storage:** GCP Cloud Storage (optional)
- **Real-Time:** WebSockets (for chat)

---

# ** Prospective Features #

## **Authentication**
### **1. Register User**  
**Rate Limit:** `5 requests per minute`  
**Endpoint:** `POST /api/auth/register`  
**Request Body:**  
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```
**Possible Responses:**
✅ **201 Created:**  
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "token": "paseto_token"
}
```
❌ **400 Bad Request:**  
```json
{
  "error": "Invalid email format"
}
```
❌ **409 Conflict:**  
```json
{
  "error": "Email already registered"
}
```

---

### **2. Login**
**Rate Limit:** `10 requests per minute`  
**Endpoint:** `POST /api/auth/login`  
**Request Body:**  
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```
**Possible Responses:**  
✅ **200 OK:**  
```json
{
  "token": "paseto_token"
}
```
❌ **401 Unauthorized:**  
```json
{
  "error": "Invalid credentials"
}
```

---

## **User Profile**
### **3. Get Profile**
**Rate Limit:** `60 requests per minute`  
**Endpoint:** `GET /api/users/{user_id}`  
**Response:**  
```json
{
  "id": "uuid",
  "name": "John Doe",
  "bio": "I love coffee!",
  "gender": "Male",
  "age": 25,
  "image_url": "https://storage.googleapis.com/bucket/image.jpg"
}
```

---

### **4. Update Profile**  
**Rate Limit:** `30 requests per minute`  
**Endpoint:** `PUT /api/users/{user_id}`  
**Request Body:**  
```json
{
  "name": "John Doe",
  "bio": "I love coffee!",
  "gender": "Male",
  "age": 25,
  "image_url": "https://storage.googleapis.com/bucket/image.jpg"
}
```
**Response:**  
✅ **200 OK:**  
```json
{
  "message": "Profile updated successfully"
}
```

---

## **Swiping & Matching**
### **5. Swipe on a User**  
**Rate Limit:** `100 swipes per day per user`  
**Endpoint:** `POST /api/swipes`  
**Request Body:**  
```json
{
  "swiper_id": "uuid",
  "swipee_id": "uuid",
  "liked": true
}
```
**Possible Responses:**  
✅ **201 Created:**  
```json
{
  "match": true
}
```
✅ **201 Created (No Match):**  
```json
{
  "match": false
}
```
❌ **429 Too Many Requests:**  
```json
{
  "error": "Daily swipe limit reached"
}
```

---

### **6. Get Matches**  
**Rate Limit:** `60 requests per minute`  
**Endpoint:** `GET /api/matches/{user_id}`  
**Response:**  
```json
[
  {
    "match_id": "uuid",
    "user_id": "uuid",
    "name": "Jane Doe",
    "image_url": "https://storage.googleapis.com/bucket/image.jpg",
    "matched_at": "2025-03-24T12:00:00Z"
  }
]
```

---

## **Chats & Messages**
### **7. Start Chat**  
**Rate Limit:** `10 requests per minute`  
**Endpoint:** `POST /api/chats`  
**Request Body:**  
```json
{
  "user1_id": "uuid",
  "user2_id": "uuid"
}
```
**Response:**  
```json
{
  "chat_id": "uuid"
}
```

---

### **8. Send Message**  
**Rate Limit:** `120 messages per user per hour`  
**Endpoint:** `POST /api/messages`  
**Request Body:**  
```json
{
  "chat_id": "uuid",
  "sender_id": "uuid",
  "content": "Hey there!"
}
```
**Response:**  
```json
{
  "message_id": "uuid",
  "sent_at": "2025-03-24T12:05:00Z"
}
```

---

### **9. Get Messages**  
**Rate Limit:** `100 requests per minute`  
**Endpoint:** `GET /api/chats/{chat_id}/messages`  
**Response:**  
```json
[
  {
    "sender_id": "uuid",
    "content": "Hey there!",
    "sent_at": "2025-03-24T12:05:00Z"
  }
]
```

---

## **Feed (Posts & Comments)**
### **10. Create Post**  
**Rate Limit:** `20 posts per day per user`  
**Endpoint:** `POST /api/posts`  
**Request Body:**  
```json
{
  "user_id": "uuid",
  "content": "Just had an amazing day!",
  "image_url": "https://storage.googleapis.com/bucket/post.jpg"
}
```
**Response:**  
```json
{
  "post_id": "uuid",
  "created_at": "2025-03-24T14:00:00Z"
}
```

---

### **11. Get Posts**  
**Rate Limit:** `60 requests per minute`  
**Endpoint:** `GET /api/posts`  
**Response:**  
```json
[
  {
    "post_id": "uuid",
    "user_id": "uuid",
    "content": "Just had an amazing day!",
    "image_url": "https://storage.googleapis.com/bucket/post.jpg",
    "created_at": "2025-03-24T14:00:00Z"
  }
]
```

---

### **12. Comment on Post**  
**Rate Limit:** `60 comments per hour per user`  
**Endpoint:** `POST /api/posts/{post_id}/comments`  
**Request Body:**  
```json
{
  "user_id": "uuid",
  "content": "That's awesome!"
}
```
**Response:**  
```json
{
  "comment_id": "uuid",
  "created_at": "2025-03-24T14:05:00Z"
}
```

---

## **Security & Rate Limiting**
- **PASETO-based JWT authentication**
- **Rate limit swipes to prevent abuse** (e.g., max 100 swipes/day)
- **WebSockets for real-time chat**
- **Cloud Storage for image uploads**  
- **IP-based rate limiting per endpoint**  
- **Account lockout after repeated failed login attempts**  
---
