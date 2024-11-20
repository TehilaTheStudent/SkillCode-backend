import ast
import json
import os  # Add this for path manipulation
from jsonschema import validate, ValidationError


def parse_input(input_string):
    try:
        return ast.literal_eval(input_string)  # Safely parse input
    except (ValueError, SyntaxError) as e:
        return f"Error parsing input: {str(e)}"


def run_test_cases(compiled_code, test_cases, function_name):
    results = []
    namespace = {}
    all_passed = True  # Track if all test cases pass
    try:
        exec(compiled_code, namespace)  # Execute user code in a separate namespace
    except Exception as e:
        return {"status": "fail", "error": "Execution failed", "details": str(e), "results": []}

    user_function = namespace.get(function_name)
    if not callable(user_function):
        return {
            "status": "fail",
            "error": f"{function_name} is not defined or callable",
            "details": None,
            "results": [],
        }

    for case in test_cases:
        try:
            # Parse each parameter string individually
            inputs = [parse_input(param) for param in case["parameters"]]
            expected_output = parse_input(case["expected_output"])

            # Invoke the user's function
            actual_output = user_function(*inputs)

            # Compare outputs
            if actual_output == expected_output:
                results.append(
                    {
                        "status": "pass",
                        "parameters": case["parameters"],  # Include parameters in the result
                        "expected_output": str(expected_output),
                        "actual_output": str(actual_output),
                    }
                )
            else:
                all_passed = False  # Mark as failed if any test case fails
                results.append(
                    {
                        "status": "fail",
                        "parameters": case["parameters"],  # Include parameters in the result
                        "expected_output": str(expected_output),
                        "actual_output": str(actual_output),
                    }
                )
        except Exception as e:
            all_passed = False  # Mark as failed if any test case raises an exception
            results.append(
                {
                    "status": "fail",
                    "parameters": case["parameters"],  # Include parameters in the result
                    "expected_output": str(case["expected_output"]),
                    "actual_output": f"Error: {str(e)}",
                }
            )

    overall_status = "success" if all_passed else "fail"
    return {"status": overall_status, "results": results, "error": None, "details": None}


def evaluate_user_code(user_code, test_cases, function_name, schema_path="./feedback_schema.json"):
    # Resolve the schema path relative to this script
    script_dir = os.path.dirname(os.path.abspath(__file__))  # Get the script's directory
    resolved_schema_path = os.path.join(script_dir, schema_path)  # Resolve schema path

    # Step 1: Compile the user's code
    user_code = "import ds_utils as utils\n" + user_code
    try:
        compiled_code = compile(user_code, filename="<user_code>", mode="exec")
    except SyntaxError as e:
        return {"status": "fail", "error": "Compilation failed", "details": str(e), "results": []}

    # Step 2: Run test cases
    results = run_test_cases(compiled_code, test_cases, function_name)

    # Step 3: Validate against schema
    try:
        with open(resolved_schema_path, "r") as schema_file:
            schema = json.load(schema_file)
        validate(instance=results, schema=schema)
    except FileNotFoundError:
        return {"status": "fail", "error": "FileNotFoundError", "details": f"Schema file not found at {resolved_schema_path}", "results": []}
    except ValidationError as e:
        return {"status": "fail", "error": "Validation failed", "details": e.message, "results": []}

    return results


if __name__ == "__main__":
    # Example Usage:
    user_code = """
def binarySearch(arr, target):
    left, right = 0, len(arr) - 1
    while left <= right:
        mid = (left + right) // 2
        if arr[mid] == target:
            return mid
        elif arr[mid] < target:
            left = mid + 1
        else:
            right = mid - 1
    return -1
    """

    test_cases = [
        {"parameters": ["[1, 2, 3, 4, 5]", "3"], "expected_output": "2"},
        {"parameters": ["[1, 2, 3, 4, 5]", "6"], "expected_output": "-1"},
    ]

    function_name = "binarySearch"

    results = evaluate_user_code(user_code, test_cases, function_name)
    print(json.dumps(results, indent=2))  # Pretty-print JSON results
