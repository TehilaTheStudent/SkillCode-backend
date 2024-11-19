from evaluator import evaluate_user_code


user_code = """def binary_search(arr: list[int], target: int) -> int:"""
test_cases = [{"parameters":["[1, 2, 3, 4, 5]","3"],"expected_output":"2"},{"parameters":["[1, 2, 3, 4, 5]","6"],"expected_output":"-1"}]
function_name = "binarySearch"

results = evaluate_user_code(user_code, test_cases, function_name)
print(results)
