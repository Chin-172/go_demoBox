# go_demoBox


Hi👋 DemoBox is a simple work of gRPC functions run between Server and Client side. You may through this project to see how the SQL function works on gRPC functions and how the gRPC communicate between server and client side.

### System Env
DB: PostgreSQL\
Host: localhost\
Port: 5432\
DB Tools: [pgAdmin4](https://www.pgadmin.org/download/)


### How to run

Run Server:
```bash
cd server
go run .\serverSide.go
```

Run Client:
```bash
cd server
go run .\serverSide.go
```

### How gRPC functions work
<img width="500" height="550" alt="image" src="https://github.com/user-attachments/assets/2f72ee73-181b-4b12-86b8-c3e4da5e5201" />
<img width="500" height="550" alt="image" src="https://github.com/user-attachments/assets/32a9a01f-5c53-4551-ba43-65276c9ccb58" />

### Oneof Field Type
In order to handle multiple data entity into one rpc function, here using `oneof` field to define the data entity type
<img width="702" height="562" alt="image" src="https://github.com/user-attachments/assets/68ab81df-ee05-4f37-bec6-87cf3e76c5a3" />

### Progress
| Date | Progress |
| --- | --- |
| 20/12/2025 | Create DB table (user_list & group_list) and define rpc service in proto file |
| 21/12/2025 | Create one streaming and one entire entity gRPC function on server side |
| 22/12/2025 | Coding the client side function to connect and call server side gRPC functions |
| 27/12/2025 | Create more relational tables and Added Login & Data Inspection Authority Check Functions  |
