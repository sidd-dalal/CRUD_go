
# Example requests:

1. Add 'Owner' info: curl -X POST http://localhost:6459/owners -d '{"name":"Siddharth"}'
2. Add 'Pet' info: curl -X POST http://localhost:6459/owners/pets?ownerId=1 -d '{"name":"kutta","type":"Dog","age":3}'
3. Update 'Owner' infocurl -X PUT http://localhost:6459/owners/update?id=1 -d '{"name":"New Siddharth"}'
4. Delete 'Owner' info: curl -X DELETE http://localhost:6459/owners/delete?id=1
5. Show all 'Owner': curl http://localhost:6459/owners/all