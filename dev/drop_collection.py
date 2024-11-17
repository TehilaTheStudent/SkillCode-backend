from pymongo import MongoClient

# Configuration
mongo_host = "localhost"  # MongoDB host (e.g., localhost or container name)
mongo_port = 27017        # Default MongoDB port
db_name = "skillcode_db"  # Database name
collection_name = "questions"  # Collection name to drop

# Function to drop collection
def drop_questions_collection():
    try:
        client = MongoClient(mongo_host, mongo_port)
        db = client[db_name]
        if collection_name in db.list_collection_names():
            db.drop_collection(collection_name)
            print(f"Collection '{collection_name}' dropped successfully.")
        else:
            print(f"Collection '{collection_name}' does not exist.")
    except Exception as e:
        print(f"Error: {e}")
    finally:
        client.close()

if __name__ == "__main__":
    drop_questions_collection()
