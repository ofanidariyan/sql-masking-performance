#!/bin/bash
make docker-down
make docker-up

# Ensure the wait-for-it script has execute permissions
chmod +x script/wait-for-it.sh

# Wait for the MySQL service to be available
./script/wait-for-it.sh localhost:3308 -- echo "MYSQL for masking data testing is up..."

sleep 10

# SQL Layer Testing
echo "Testing masking data with the SQL layer with a sample of 10000 records"
for i in 1 2 3 ; do
    echo "Testing $i"
    make masking-sql NUM_RECORDS=10000
done

echo "Testing masking data with the SQL layer with a sample of 100000 records"
for i in 1 2 3 ; do
    echo "Testing $i"
    make masking-sql NUM_RECORDS=100000
done

echo "Testing masking data with the SQL layer with a sample of 1000000 records"
for i in 1 2 3 ; do
    echo "Testing $i"
    make masking-sql NUM_RECORDS=1000000
done
########################################

make docker-down
make docker-up

# Ensure the wait-for-it script has execute permissions
chmod +x script/wait-for-it.sh

# Wait for the MySQL service to be available
./script/wait-for-it.sh localhost:3308 -- echo "MYSQL for masking data testing is up..."

sleep 10

# Golang Layer Testing
echo "Testing masking data with the Golang layer with a sample of 10000 records"
for i in 1 2 3 ; do
    echo "Testing $i"
    make masking-golang NUM_RECORDS=10000
done

echo "Testing masking data with the Golang layer with a sample of 100000 records"
for i in 1 2 3 ; do
    echo "Testing $i"
    make masking-golang NUM_RECORDS=100000
done

echo "Testing masking data with the Golang layer with a sample of 1000000 records"
for i in 1 2 3 ; do
    echo "Testing $i"
    make masking-golang NUM_RECORDS=1000000
done

#########################################

make docker-down
make docker-up

# Ensure the wait-for-it script has execute permissions
chmod +x script/wait-for-it.sh

# Wait for the MySQL service to be available
./script/wait-for-it.sh localhost:3308 -- echo "MYSQL for masking data testing is up..."

sleep 10

# Golang Layer Testing
echo "Testing masking data with the SQL & Golang layer with a sample of 10000 records"
for i in 1 2 3 ; do
    echo "Testing $i"
    make masking-sql-golang NUM_RECORDS=10000
done

echo "Testing masking data with the SQL & Golang layer with a sample of 100000 records"
for i in 1 2 3 ; do
    echo "Testing $i"
    make masking-sql-golang NUM_RECORDS=100000
done

echo "Testing masking data with the SQL & Golang layer with a sample of 1000000 records"
for i in 1 2 3 ; do
    echo "Testing $i"
    make masking-sql-golang NUM_RECORDS=1000000
done

#make docker-down