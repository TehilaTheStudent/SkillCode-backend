class Solution:
    def merge(nums1: list[int], m: int, nums2: list[int], n: int) -> list[int]:        
        j = m + n - 1
        i, k = m - 1, n - 1
        while i >= 0 and k >= 0:
            if nums1[i] > nums2[k]:
                nums1[j] = nums1[i]
                i -= 1
            else:
                nums1[j] = nums2[k]
                k -= 1
            j -= 1
        while k >= 0:
            nums1[j] = nums2[k]
            k -= 1
            j -= 1
        return nums1
test_cases = [
    {"input": {"nums1": [1, 2, 3, 0, 0, 0], "m": 3, "nums2": [2, 5, 6], "n": 3}, "expected_output": [1, 2, 2, 3, 5, 6]},
    {"input": {"nums1": [4, 5, 6, 0, 0, 0], "m": 3, "nums2": [1, 2, 3], "n": 3}, "expected_output": [1, 2, 3, 4, 5, 6]},
]


# Initialize the Solution class
solution = Solution()

def main():
    # Results container
    results = []

    # Execute all test cases
    for i, test_case in enumerate(test_cases):
        try:
            # Extract inputs dynamically
            inputs = test_case['input']
            expected_output = test_case['expected_output']

            # Dynamically call the user's function
            actual_output = solution.merge(**inputs)  # Adjust function name if needed

            # Compare actual and expected outputs
            if actual_output == expected_output:
                results.append({
                    'status': 'passed',
                    'input': inputs,
                    'expected': expected_output,
                    'actual': actual_output
                })
            else:
                results.append({
                    'status': 'failed',
                    'input': inputs,
                    'expected': expected_output,
                    'actual': actual_output
                })
        except Exception as e:
            # Handle runtime errors gracefully
            results.append({
                'status': 'runtime_error',
                'input': inputs,
                'expected': expected_output,
                'error': str(e)
            })

    # Print results in JSON-like format
    import json
    print(json.dumps(results, indent=4))

if __name__ == "__main__":
    main()
