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

# Questions for each data structure
questions = [
    # Tree Question
    {
        "title": "Sum of Tree Nodes",
        "description": "Write a function to calculate the sum of all values in a binary tree.",
        "difficulty": "Medium",
        "category": "Tree",
        "examples": [
            {
                "parameters": ["[1, 2, 3, null, 5]"],  # Tree generated as: 1 -> 2, 3 -> 5
                "expected_output": "11"
            }
        ],
        "test_cases": [
            {
                "parameters": ["[1, 2, 3, null, 5]"],
                "expected_output": "11"
            },
            {
                "parameters": ["[4, 9, null, 5, null]"],
                "expected_output": "18"
            },
            {
                "parameters": ["[0, 0, 0]"],
                "expected_output": "0"
            }
        ],
        "function_config": {
            "name": "sumTree",
            "parameters": [
                {
                    "name": "root",
                    "param_type": {
                        "type": "TreeNode"
                    }
                }
            ],
            "return_type": {
                "type": "Integer"
            }
        }
    },
    # Linked List Question
    {
        "title": "Reverse Linked List",
        "description": "Write a function to reverse a singly linked list.",
        "difficulty": "Medium",
        "category": "LinkedList",
        "examples": [
            {
                "parameters": ["[1, 2, 3, 4]"],
                "expected_output": "[4, 3, 2, 1]"
            }
        ],
        "test_cases": [
            {
                "parameters": ["[1, 2, 3, 4]"],
                "expected_output": "[4, 3, 2, 1]"
            },
            {
                "parameters": ["[5, 10]"],
                "expected_output": "[10, 5]"
            },
            {
                "parameters": ["[1]"],
                "expected_output": "[1]"
            }
        ],
        "function_config": {
            "name": "reverseLinkedList",
            "parameters": [
                {
                    "name": "head",
                    "param_type": {
                        "type": "ListNode"
                    }
                }
            ],
            "return_type": {
                "type": "ListNode"
            }
        }
    },
    # Graph Question
    {
        "title": "Find All Neighbors",
        "description": "Write a function to return the list of neighbors for a given node in an undirected graph.",
        "difficulty": "Easy",
        "category": "Graph",
        "examples": [
            {
                "parameters": ["[[1, 2], [2, 3], [3, 4]]", "2"],
                "expected_output": "[1, 3]"
            }
        ],
        "test_cases": [
            {
                "parameters": ["[[1, 2], [2, 3], [3, 4]]", "2"],
                "expected_output": "[1, 3]"
            },
            {
                "parameters": ["[[5, 6], [6, 7], [7, 8], [8, 5]]", "5"],
                "expected_output": "[6, 8]"
            },
            {
                "parameters": ["[[1, 2], [3, 4]]", "3"],
                "expected_output": "[4]"
            }
        ],
        "function_config": {
            "name": "findNeighbors",
            "parameters": [
                {
                    "name": "edges",
                    "param_type": {
                        "type": "Array",
                        "type_children": {
                            "type": "Array",
                            "type_children": {
                                "type": "Integer"
                            }
                        }
                    }
                },
                {
                    "name": "node",
                    "param_type": {
                        "type": "Integer"
                    }
                }
            ],
            "return_type": {
                "type": "Array",
                "type_children": {
                    "type": "Integer"
                }
            }
        }
    },
    # Array Question
    {
        "title": "Find Maximum in Array",
        "description": "Write a function to find the maximum value in an array.",
        "difficulty": "Easy",
        "category": "Array",
        "examples": [
            {
                "parameters": ["[3, 1, 4, 1, 5]"],
                "expected_output": "5"
            }
        ],
        "test_cases": [
            {
                "parameters": ["[3, 1, 4, 1, 5]"],
                "expected_output": "5"
            },
            {
                "parameters": ["[-1, -4, -3, -2]"],
                "expected_output": "-1"
            },
            {
                "parameters": ["[7]"],
                "expected_output": "7"
            }
        ],
        "function_config": {
            "name": "findMax",
            "parameters": [
                {
                    "name": "arr",
                    "param_type": {
                        "type": "Array",
                        "type_children": {
                            "type": "Integer"
                        }
                    }
                }
            ],
            "return_type": {
                "type": "Integer"
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
