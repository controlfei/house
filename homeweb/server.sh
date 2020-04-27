#启动redis服务
#redis-server ./conf/redis.conf

#启动trackerd
fdfs_trackerd  /home/linpengfei/go/src/house/homeweb/conf/tracker.conf restart

#启动storaged
fdfs_storaged  /home/linpengfei/go/src/house/homeweb/conf/storage.conf restart