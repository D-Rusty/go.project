
# 会用到一些命令
1. curl -i -X POST -H "Content-Type: application/json"  -d '{"content":"My first post","author":"Sau Sheong"}' http://127.0.0.1:8080/post/
2. curl -i -X DELETE http://127.0.0.1:8080/post/1
3. curl -i -X GET http://127.0.0.1:8080/post/1
4. curl -i -X PUT -H "Content-Type: application/json"  -d '{"content":"Updated post","author":"Sau Sheong"}' http://127.0.0.1:8080/post/1