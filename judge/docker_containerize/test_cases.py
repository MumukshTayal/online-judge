# import time
# import importlib
# import sys
# import os

# test_cases = [
#     {
#         "input": [range(1,100)],
#         "expected_output": 5050
#     },
#     {
#         "input": [[4, 5, 6]],
#         "expected_output": 15
#     },
#     # Add more test cases as needed
# ]


# def run_test_cases():
#     # Get the function name from the environment variable
#     function_name = os.environ.get('FUNCTION_NAME') or 'add'
#     # function_name = "add"
#     if not function_name:
#         print("FUNCTION_NAME environment variable not set.")
#         return False
    
#     # Import the user's code module
#     user_code = importlib.import_module('code')

#     # Get the function from the user's code
#     user_function = getattr(user_code, function_name)

#     # Call the user's function with the test case inputs
#     start_time = time.time()
#     for i, test_case in enumerate(test_cases):
#         user_output = user_function(*test_case["input"])
#         if user_output != test_case["expected_output"]:
#             print(f"{i}/{len(test_cases)} test case(s) passed! Test case failed: Input={test_case['input']} Expected output={test_case['expected_output']} Actual output={user_output}, Time Elapsed: {time.time() - start_time} seconds")
#             return False
        
#     print(f"{len(test_cases)}/{len(test_cases)} test cases passed! Time Elapsed: {time.time() - start_time} seconds")
#     return True

# if __name__ == "__main__":
#     if len(sys.argv) < 1:
#         print("Usage: python test_cases.py <user_code_file> <function_name>")
#         sys.exit(1)

#     # function_name = sys.argv[1]
#     run_test_cases()

import subprocess


def main():
    with open("input.txt", "r") as file:
        input_data = file.read()

    try:
        output = subprocess.run(
            ["python", "code.py"],
            input=input_data,
            text=True,
            capture_output=True,
            check=True,
        )
        
    except subprocess.CalledProcessError as e:
        print("An error occurred while executing code.py:")
        print(e.stderr)
    
    with open("output.txt", "r") as file:
        output_data = file.read()

    list_output_data = output_data.split("\n")
    # print(list_output_data)
    # print(output.stdout)
    actual_output = output.stdout.split("\n")
    for i in range(len(list_output_data)):
        if int(list_output_data[i]) != int(actual_output[i]):
            print(f"{i}/{len(list_output_data)} test cases Passed") 
            print(f"Expected: {list_output_data[i]}")
            print(f"Actual: {actual_output[i]}")
            break

    print(f"{len(list_output_data)}/{len(list_output_data)} test cases Passed :)")

if __name__ == "__main__":
    main()