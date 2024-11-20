from evaluator import evaluate_user_code
import json
user_code = """def is_palindrome(s: str) -> bool:
    pass














"""
test_cases = [{"parameters":["'madam'"],"expected_output":"true"},{"parameters":["'hello'"],"expected_output":"false"},{"parameters":["'A man a plan a canal Panama'"],"expected_output":"true"}]
function_name = "isPalindrome"

results = evaluate_user_code(user_code, test_cases, function_name)
print(json.dumps(results, indent=2))