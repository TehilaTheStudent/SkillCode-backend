import ast
import json
import os
from jsonschema import validate, ValidationError
import converter







def run_test_cases(compiled_code, test_cases, function_name,function_config):
    """
    Run the provided test cases against the compiled user code.

    Args:
        compiled_code (code): The compiled user code.
        test_cases (list): A list of test cases, each containing 'parameters' and 'expected_output'.
        function_name (str): The name of the function to test.

    Returns:
        dict: A dictionary containing the overall status, results of each test case, and error details if any.
    """
    results = []
    namespace = {}
    all_passed = True  # Track if all test cases pass

    try:
        exec(compiled_code, namespace)  # Execute user code in a separate namespace
    except Exception as e:
        return {
            "status": "fail",
            "error": "compilation",
            "details": str(e),
            "results": [],
        }

    user_function = namespace.get(function_name)
    if not callable(user_function):
        return {
            "status": "fail",
            "error": "compilation",
            "details": f"{function_name} is not defined or callable",
            "results": [],
        }

    for case in test_cases:
        try:
            # Parse each parameter string individually
            inputs = [converter.listy_to_type(case['parameters'][i],function_config['parameters'][i]['param_type']) for i in range(len(case["parameters"]))]
            expected_output = converter.listy_to_type(case["expected_output"],function_config['return_type'])

            # Invoke the user's function
            actual_output = user_function(*inputs)

            # Compare outputs
            if actual_output == expected_output:
                results.append(
                    {
                        "status": "pass",   
                        "parameters": case["parameters"],
                        "expected_output": str(expected_output),
                        "actual_output": str(actual_output),
                    }
                )
            else:
                all_passed = False
                results.append(
                    {
                        "status": "fail",
                        "parameters": case["parameters"],
                        "expected_output": str(expected_output),
                        "actual_output": str(actual_output),
                    }
                )
        except Exception as e:
            all_passed = False
            results.append(
                {
                    "status": "fail",
                    "parameters": case["parameters"],
                    "expected_output": str(case["expected_output"]),
                    "actual_output": f"Error: {str(e)}",
                }
            )

    overall_status = "success" if all_passed else "fail"
    return {
        "status": overall_status,
        "results": results,
        "error": None if all_passed else "fail tests",
        "details": None,
    }


def evaluate_user_code(
    user_code, test_cases, function_name,function_config, schema_path="../feedback_schema.json"
):
    """
    Evaluate the user's code by compiling it, running test cases, and validating the results against a schema.

    Args:
        user_code (str): The user's code as a string.
        test_cases (list): A list of test cases, each containing 'parameters' and 'expected_output'.
        function_name (str): The name of the function to test.
        schema_path (str): The path to the JSON schema file for validation.

    Returns:
        dict: A dictionary containing the overall status, results of each test case, and error details if any.
    """
    try:
        # Resolve the schema path relative to this script
        script_dir = os.path.dirname(
            os.path.abspath(__file__)
        )  # Get the script's directory
        resolved_schema_path = os.path.join(
            script_dir, schema_path
        )  # Resolve schema path

        # Step 1: Compile the user's code
        user_code = "import ds_utils as utils\n" + user_code
        try:
            compiled_code = compile(user_code, filename="<user_code>", mode="exec")
        except SyntaxError as e:
            return {
                "status": "fail",
                "error": "compilation",
                "details": str(e),
                "results": [],
            }

        # Step 2: Run test cases
        results = run_test_cases(compiled_code, test_cases, function_name,function_config)

        # Step 3: Validate against schema
        try:
            with open(resolved_schema_path, "r") as schema_file:
                schema = json.load(schema_file)
            validate(instance=results, schema=schema)
        except FileNotFoundError:
            return {
                "status": "fail",
                "error": "internal server error",
                "details": f"Schema file not found at {resolved_schema_path}",
                "results": [],
            }
        except ValidationError as e:
            return {
                "status": "fail",
                "error": "internal server error",
                "details": f"Schema validation error: {e.message}",
                "results": [],
            }

        return results
    except Exception as e:
        # Catch any unexpected errors and mark as internal server error
        return {
            "status": "fail",
            "error": "internal server error",
            "details": str(e),
            "results": [],
        }
