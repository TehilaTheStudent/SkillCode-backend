from evaluator import evaluate_user_code
import ds_utils as utils

user_code = """def binarySearch(arr, target):
    left, right = 0, len(arr) - 1
    while left <= right:
        mid = (left + right) // 2
        if arr[mid] == target:
            return mid
        elif arr[mid] < target:
            left = mid + 1
        else:
            right = mid - 1
    return -1"""
test_cases = [{"parameters":["[1, 2, 3, 4, 5]","3"],"expected_output":"2"},{"parameters":["[1, 2, 3, 4, 5]","3"],"expected_output":"2"}]
function_name = "binarySearch"

results = evaluate_user_code(user_code, test_cases, function_name)
print(results)