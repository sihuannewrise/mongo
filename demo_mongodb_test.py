import pymongo

myclient = pymongo.MongoClient("mongodb://mongodbuser:aw3se4dr5@almeling.ru:27017/")

mydb = myclient["mydatabase"]
mycol = mydb["customers"]

mydict = { "name": "Helen", "address": "Lenina 8A" }

x = mycol.insert_one(mydict)
print(x.inserted_id)

print(mycol)

# docker exec -it mongodb mongosh -u mongodbuser -p aw3se4dr5
# mongosh --host mogdodb.almeling.ru --username mongodbuser --password aw3se4dr5 --authenticationDatabase knhdb
# mongosh --host mogdodb.almeling.ru -u mongodbuser -p aw3se4dr5
# mongosh mongodb://mongodbuser:aw3se4dr5@mongodb.almeling.ru:27017/mydatabase -u mongodbuser -p aw3se4dr5

# start mongo shell
# docker exec -it mongodb mongosh -u <mongodbuser> -p <password>
# db.auth({user: "mongodbuser", passwordPrompt()})
# db.auth({user: "mongodbuser", pwd: "aw3se4dr5"})
