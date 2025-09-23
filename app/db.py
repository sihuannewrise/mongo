from pymongo import MongoClient

client = MongoClient("mongodb://mongodbuser:aw3se4dr5@almeling.ru:27017/")
db = client["knhdb"]
students_collection = db["students"]
