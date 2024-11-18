import ast

def parse_input(input_string):
    try:
        return ast.literal_eval(input_string)  # Safer than eval()
    except (ValueError, SyntaxError) as e:
        return f"Error parsing input: {str(e)}"


def run_test_cases(compiled_code, test_cases, function_name):
    results = []
    namespace = {}
    try:
        exec(compiled_code, namespace)  # Execute user code in a separate namespace
    except Exception as e:
        return {"error": "Execution failed", "details": str(e)}

    user_function = namespace.get(function_name)
    if not callable(user_function):
        return {"error": f"{function_name} is not defined or callable"}

    for case in test_cases:
        try:
            # Parse each parameter string individually
            inputs = [parse_input(param) for param in case["parameters"]]
            expected_output = parse_input(case["expected_output"])

            # Invoke the user's function
            actual_output = user_function(*inputs)

            # Compare outputs
            if actual_output == expected_output:
                results.append({"status": "pass"})
            else:
                results.append({
                    "status": "fail",
                    "expected_output": expected_output,
                    "actual_output": actual_output,
                })
        except Exception as e:
            # Handle errors in execution
            results.append({
                "status": "fail",
                "expected_output": case["expected_output"],
                "actual_output": f"Error: {str(e)}",
            })

    return results


def evaluate_user_code(user_code, test_cases, function_name):
    # Step 1: Compile the user's code
    try:
        compiled_code = compile(user_code, filename="<user_code>", mode="exec")
    except SyntaxError as e:
        return {"error": "Compilation failed", "details": str(e)}

    # Step 2: Run test cases
    results = run_test_cases(compiled_code, test_cases, function_name)
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
    print(results)


    # end of example usage
