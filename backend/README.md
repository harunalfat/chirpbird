## Chirpbird Backend

To run this project, copy the `.env.example` and rename it as `.env`. Change necessary variables as needed
```
PORT=8000
REDIS_ADDRESS=localhost:6379 
REDIS_PASSWORD=
REDIS_SENTINEL_ADDRESSES= <-- needed if you plan to use replication
REDIS_MASTER_NAME= <-- needed if you plan to use replication
MONGODB_CONN_URI=mongodb://localhost:27017 <--- comma separated for replication
```
- Run `./bin/air`
- Server should run on `http://localhost:${PORT}`