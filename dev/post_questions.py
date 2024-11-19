import requests

# Configuration
backend_url = "http://localhost:8080/skillcode/questions"  # Replace with your actual backend URL

# Base question template
base_question_template = {
    "title": "",
    "description": "",
    "difficulty": "",
    "category": "",
    "stats": 0,
    "examples": [],
    "test_cases": [],
    "function_config": {
        "name": "",
        "parameters": [],
        "return_type": {}
    },
    "languages": ["Python", "JavaScript", "Go"]
}

# Question variations
questions = [
    {
        "title": "Return 1",
        "description": "Write a function that always returns the number 1.",
        "difficulty": "Easy",
        "category": "Matrix",
        "examples": [
            {
                "parameters": [],
                "expected_output": "1"
            }
        ],
        "test_cases": [
            {
                "parameters": [],
                "expected_output": "1"
            },
            {
                "parameters": [],
                "expected_output": "1"
            }
        ],
        "function_config": {
            "name": "returnOne",
            "parameters": [],
            "return_type": {
                "type": "Integer"
            }
        }
    },
    {
        "title": "Sum of Two Numbers",
        "description": "Write a function that takes two numbers and returns their sum.",
        "difficulty": "Easy",
        "category": "LinkedList",
        "examples": [
            {
                "parameters": ["2", "3"],
                "expected_output": "5"
            }
        ],
        "test_cases": [
            {
                "parameters": ["5", "7"],
                "expected_output": "12"
            },
            {
                "parameters": ["-1", "4"],
                "expected_output": "3"
            }
        ],
        "function_config": {
            "name": "sumTwoNumbers",
            "parameters": [
                {
                    "name": "a",
                    "param_type": {"type": "Integer"}
                },
                {
                    "name": "b",
                    "param_type": {"type": "Integer"}
                }
            ],
            "return_type": {
                "type": "Integer"
            }
        }
    },
    {
        "title": "Binary Search",
        "description": "Implement binary search to find the target value in a sorted array.",
        "difficulty": "Medium",
        "category": "Array",
        "examples": [
            {
                "parameters": ["[1, 2, 3, 4, 5]", "3"],
                "expected_output": "2"
            }
        ],
        "test_cases": [
            {
                "parameters": ["[1, 2, 3, 4, 5]", "3"],
                "expected_output": "2"
            },
            {
                "parameters": ["[1, 2, 3, 4, 5]", "6"],
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
                        "type_children": {"type": "Integer"}
                    }
                },
                {
                    "name": "target",
                    "param_type": {"type": "Integer"}
                }
            ],
            "return_type": {
                "type": "Integer"
            }
        }
    },
    {
        "title": "Factorial",
        "description": "Write a function that calculates the factorial of a given number.",
        "difficulty": "Medium",
        "category": "LinkedList",
        "examples": [
            {
                "parameters": ["5"],
                "expected_output": "120"
            }
        ],
        "test_cases": [
            {
                "parameters": ["0"],
                "expected_output": "1"
            },
            {
                "parameters": ["3"],
                "expected_output": "6"
            },
            {
                "parameters": ["7"],
                "expected_output": "5040"
            }
        ],
        "function_config": {
            "name": "factorial",
            "parameters": [
                {
                    "name": "n",
                    "param_type": {"type": "Integer"}
                }
            ],
            "return_type": {
                "type": "Integer"
            }
        }
    },
    {
        "title": "Is Palindrome",
        "description": "Write a function that checks if a given string is a palindrome.",
        "difficulty": "Hard",
        "category": "String",
        "examples": [
            {
                "parameters": ["'racecar'"],
                "expected_output": "true"
            }
        ],
        "test_cases": [
            {
                "parameters": ["'madam'"],
                "expected_output": "true"
            },
            {
                "parameters": ["'hello'"],
                "expected_output": "false"
            },
            {
                "parameters": ["'A man a plan a canal Panama'"],
                "expected_output": "true"
            }
        ],
        "function_config": {
            "name": "isPalindrome",
            "parameters": [
                {
                    "name": "s",
                    "param_type": {"type": "String"}
                }
            ],
            "return_type": {
                "type": "Boolean"
            }
        }
    }
]

# Function to post questions
def post_questions():
    for i, question_data in enumerate(questions, 1):
        question = base_question_template.copy()
        question.update(question_data)
        response = requests.post(backend_url, json=question)
        if response.status_code == 201:
            print(f"Question {i} posted successfully: {question['title']}")
        else:
            print(f"Failed to post question {i}: {question['title']}. Response: {response.text}")

if __name__ == "__main__":
    post_questions()
