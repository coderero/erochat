CREATE KEYSPACE erochat WITH REPLICATION = { 
    'class' : 'SimpleStrategy', 
    'replication_factor' : 1 
    };