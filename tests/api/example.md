## POST /quizzes/submissions

```json
{
  "user_id": "user-123",
  "quiz_id": "quiz-456",
  "quiz_type": "choice",
  "note": "ทำครั้งแรก",
  "status": "in_progress",
  "choice_answer": [
    {
      "answer": "A",
      "result": true
    },
    {
      "answer": "B",
      "result": false
    }
  ],
  "ordering_answer": [
    {
      "answer": "ขั้นตอนที่ 1",
      "order": 1,
      "result": true
    },
    {
      "answer": "ขั้นตอนที่ 2",
      "order": 2,
      "result": true
    }
  ],
  "missing_answer": {
    "answer": "คำตอบที่เติม",
    "result": true
  }
}
```

### Response (Success):

```json
{
  "id": "history-789",
  "user_id": "user-123",
  "quiz_id": "quiz-456",
  "quiz_type": "choice",
  "note": "ทำครั้งแรก",
  "status": "in_progress"
}
```

## GetQuizHistoryByID (GET /quiz-histories/:history_id)

```json
{
  "quiz_history": {
    "id": "history-789",
    "user_id": "user-123",
    "quiz_id": "quiz-456",
    "quiz_type": "choice",
    "note": "ทำครั้งแรก",
    "status": "in_progress",
    "choice_answer": [
      {
        "answer": "A",
        "result": true
      },
      {
        "answer": "B",
        "result": false
      }
    ],
    "ordering_answer": [
      {
        "answer": "ขั้นตอนที่ 1",
        "order": 1,
        "result": true
      },
      {
        "answer": "ขั้นตอนที่ 2",
        "order": 2,
        "result": true
      }
    ],
    "missing_answer": {
      "answer": "คำตอบที่เติม",
      "result": true
    }
  }
}
```

## GetAllQuizHistoryByQuizID (GET /quizzes/:quiz_id/histories)

```json

{
  "quiz_histories": [
    {
      "id": "history-789",
      "user_id": "user-123",
      "quiz_id": "quiz-456",
      "quiz_type": "choice",
      "note": "ทำครั้งแรก",
      "status": "in_progress",
      "choice_answer": [
        {
          "answer": "A",
          "result": true
        },
        {
          "answer": "B",
          "result": false
        }
      ]
    },
    {
      "id": "history-790",
      "user_id": "user-123",
      "quiz_id": "quiz-456",
      "quiz_type": "ordering",
      "note": "ทำครั้งที่สอง",
      "status": "in_progress",
      "ordering_answer": [
        {
          "answer": "ขั้นตอนที่ 1",
          "order": 1,
          "result": true
        },
        {
          "answer": "ขั้นตอนที่ 2",
          "order": 2,
          "result": false
        }
      ]
    }
  ]
}
```

### GET /users/:user_id/courses/:course_id/progress

```json
{
  "course_id": "course-123",
  "user_id": "user-456",
  "total_quizzes": 5,
  "answered_quizzes": 3,
  "progress_percent": 60,
  "is_complete": false
}
```
