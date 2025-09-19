import pymongo

myclient = pymongo.MongoClient("mongodb://mongodbuser:aw3se4dr5@mongodb.almeling.ru:433/")

mydb = myclient["mydatabase"]
mycol = mydb["customers"]

# mydict = { "name": "Matwey", "address": "Highlands 23" }

# x = mycol.insert_one(mydict)
# print(x.inserted_id)

print(mycol)
