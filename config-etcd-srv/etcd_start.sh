etcd --config-file=/home/secondtonone/workspace/goProject/src/lbaas/config-etcd-srv/config1.yml &
sleep 5
etcd --config-file=/home/secondtonone/workspace/goProject/src/lbaas/config-etcd-srv/config2.yml &
sleep 5
etcd --config-file=/home/secondtonone/workspace/goProject/src/lbaas/config-etcd-srv/config3.yml &
sleep 6