# test_cases.py
import importlib
import sys
import os

test_cases = [
    {
        "input": [[1,2,3]],
        "expected_output": 6
    },
    {
        "input": [[4, 5, 6]],
        "expected_output": 15
    },
    # Add more test cases as needed
]


def run_test_cases():
    # Get the function name from the environment variable
    function_name = os.environ.get('FUNCTION_NAME') or 'add'
    # function_name = "add"
    if not function_name:
        print("FUNCTION_NAME environment variable not set.")
        return False
    
    # Import the user's code module
    user_code = importlib.import_module('code')

    # Get the function from the user's code
    user_function = getattr(user_code, function_name)

    # Call the user's function with the test case inputs
    for i, test_case in enumerate(test_cases):
        user_output = user_function(*test_case["input"])
        if user_output != test_case["expected_output"]:
            print(f"{i} test case(s) passed! Test case failed: Input={test_case['input']}, Expected output={test_case['expected_output']}, Actual output={user_output}")
            return False

    print("All test cases passed!")
    return True

if __name__ == "__main__":
    if len(sys.argv) < 1:
        print("Usage: python test_cases.py <user_code_file> <function_name>")
        sys.exit(1)

    # function_name = sys.argv[1]
    run_test_cases()