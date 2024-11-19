from evaluator import evaluate_user_code
import ds_utils as utils

user_code = """def factorial(n: int) -> int:"""
test_cases = [{"parameters":["0"],"expected_output":"1"},{"parameters":["3"],"expected_output":"6"},{"parameters":["7"],"expected_output":"5040"}]
function_name = "factorial"

results = evaluate_user_code(user_code, test_cases, function_name)
print(results)