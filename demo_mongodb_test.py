import pymongo

myclient = pymongo.MongoClient("mongodb://mongodbuser:aw3se4dr5@mongodb.almeling.ru:27017/")

mydb = myclient["mydatabase"]
mycol = mydb["customers"]

mydict = { "name": "Kiwi", "address": "Lenina 17A" }

x = mycol.insert_one(mydict)
print(x.inserted_id)

print(mycol)
