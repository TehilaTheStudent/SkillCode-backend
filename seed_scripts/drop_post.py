import os
import json
from pymongo import MongoClient
# Fetch environment variables
MONGO_URI = os.getenv("MONGO_URI", "mongodb://localhost:27017")
DATABASE_NAME = os.getenv("DATABASE_NAME", "skillcode_db")
COLLECTION_NAME = "questions"

# Read questions from questions.json
with open('questions.json', 'r') as file:
    questions = json.load(file)

# Select 10 questions
data = questions[:10]

# Connect to MongoDB
client = MongoClient(MONGO_URI)
db = client[DATABASE_NAME]
collection = db[COLLECTION_NAME]

# Drop the collection if it exists
collection.drop()

# Insert seed data
collection.insert_many(data)

print(f"Seeded {len(data)} documents into the '{COLLECTION_NAME}' collection.")
