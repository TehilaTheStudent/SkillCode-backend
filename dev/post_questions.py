import requests

# Configuration
backend_url = "http://localhost:8080/skillcode/questions"  # Replace with your actual backend URL
question_template = {
    "title": "Binary Search",
    "description": "Implement binary search to find the target value in a sorted array.",
    "difficulty": "Easy",
    "category": "Array",
    "stats": 0,
    "examples": [
        {
            "parameters": [
                "[1, 2, 3, 4, 5]",
                "3"
            ],
            "expected_output": "2"
        }
    ],
    "test_cases": [
        {
            "parameters": [
                "[1, 2, 3, 4, 5]",
                "3"
            ],
            "expected_output": "2"
        },
        {
            "parameters": [
                "[1, 2, 3, 4, 5]",
                "6"
            ],
            "expected_output": "-1"
        }
    ],
    "function_config": {
        "name": "binarySearch",
        "parameters": [
            {
                "name": "arr",
                "param_type": {
                    "type": "Array",
                    "type_children": {
                        "type": "Integer"
                    }
                }
            },
            {
                "name": "target",
                "param_type": {
                    "type": "Integer"
                }
            }
        ],
        "return_type": {
            "type": "Integer"
        }
    },
    "languages": ["Python", "JavaScript", "Go"]
}

# Function to post questions
def post_questions():
    for i in range(5):
        question = question_template.copy()
        response = requests.post(backend_url, json=question)
        if response.status_code == 201:
            print(f"Question {i+1} posted successfully.")
        else:
            print(f"Failed to post question {i+1}. Response: {response.text}")

if __name__ == "__main__":
    post_questions()
