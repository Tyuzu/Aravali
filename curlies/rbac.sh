curl -i -X POST http://localhost:8080/products \
  -H "Authorization: Bearer YOUR_JWT_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":"p101","name":"Mechanical Keyboard","price":129.99}'
