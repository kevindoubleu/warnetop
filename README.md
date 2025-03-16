# warnetop

Software to manage internet cafe(s)

## features

### v1.0

all warnet operations are manpower (operator/cashier) based: pc availability, time logging

- [ ] CRUD available `Device`s in the cafe
- [ ] CRUD ongoing `PlaySession`s



## Setup

### DB

Postgresql 16 running on `localhost:5432`

```sql
CREATE DATABASE warnetop;
CREATE USER warnetop_admin WITH ENCRYPTED PASSWORD 'secret';
GRANT CONNECT ON DATABASE warnetop TO warnetop_admin;
GRANT ALL PRIVILEGES ON DATABASE warnetop TO warnetop_admin;
GRANT ALL PRIVILEGES ON devices IN SCHEMA public TO warnetop_admin;
```
