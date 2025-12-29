# go_demoBox


Hi👋 DemoBox is a simple work of gRPC functions run between Server and Client side. You may through this project to see how the SQL function works on gRPC functions and how the gRPC communicate between server and client side.
<img width="1329" height="687" alt="image" src="https://github.com/user-attachments/assets/63da4673-b28a-46c2-a991-0a9048b9b76d" />

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

### Progress
| Date | Progress |
| --- | --- |
| 20/12/2025 | Create DB table (user_list & group_list) and define rpc service in proto file |
| 21/12/2025 | Create one streaming and one entire entity gRPC function on server side |
| 22/12/2025 | Coding the client side function to connect and call server side gRPC functions |
| 27/12/2025 | Create more relational tables and Added Login & Data Inspection Authority Check Functions  |
