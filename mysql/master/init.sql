create user 'replica_user'@'%' identified by 'password';
grant replication slave on *.* to 'replica_user'@'%' with grant option;
flush privileges;
