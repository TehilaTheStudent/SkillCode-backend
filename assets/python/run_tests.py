from evaluator import evaluate_user_code
import ds_utils as utils

user_code = """def return_one() -> int:"""
test_cases = [{"parameters":[],"expected_output":"1"},{"parameters":[],"expected_output":"1"}]
function_name = "returnOne"

results = evaluate_user_code(user_code, test_cases, function_name)
print(results)