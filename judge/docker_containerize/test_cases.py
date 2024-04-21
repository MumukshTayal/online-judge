# import subprocess
# import time


# def main():
#     with open("input.txt", "r") as file:
#         input_data = file.read()

#     try:
#         start_time = time.time()
#         output = subprocess.run(
#             ["python", "code.py"],
#             input=input_data,
#             text=True,
#             capture_output=True,
#             check=True,
#         )
#         total_time = time.time() - start_time
#     except subprocess.CalledProcessError as e:
#         print("An error occurred while executing your code:")
#         print(e.stderr)
#         return
    
#     with open("output.txt", "r") as file:
#         output_data = file.read()

#     list_output_data = output_data.split("\n")
#     count = 0
#     actual_output = output.stdout.split("\n")
#     for i in range(len(list_output_data)):
#         if i >= len(actual_output):
#             break
#         count += 1
#         if list_output_data[i] != actual_output[i]:
#             count -= 1

#     print(f"{count}/{len(list_output_data)} test cases Passed :)")
#     print(f"Time Elapsed: {int(total_time*1000)} msec")

# if __name__ == "__main__":
#     main()

import subprocess
import time
import os
import stat

def main():
    print("Inside the main function")
    with open("input.txt", "r") as file:
        input_data = file.read()

    language = os.environ.get("LANG")
    print("language: ", language)
    try:
        print(language)
        if language == "py":
            start_time = time.time()
            output = subprocess.run(
                ["python", "pytcode.py"],
                input=input_data,
                text=True,
                capture_output=True,
                check=True,
            )
            total_time = time.time() - start_time
        elif language == "cpp":
            # return
            with open("ccode.txt", "r") as file:
                code_data = file.read()
            # print("Code Data: ", code_data)
            # return
            # Create code.cpp file and copy code.txt contents to it
            with open("code.cpp", "w") as file:
                if not os.path.exists("code.cpp"):
                    print("code.cpp file does not exist")
                    return

                file.write(code_data)
            # return


            # Compile the C++ code
            # return

            # os.chmod("code.cpp", stat.S_IRUSR | stat.S_IWUSR)
            compile_output = subprocess.run(
                ["g++", "-o", "code", "code.cpp"],
                capture_output=True,
                check=True,
            )
            
            if compile_output.returncode != 0:
                print("An error occurred while compiling your code:")
                print(compile_output.stderr)
                return
            # Run the compiled C++ code
            start_time = time.time()
            output = subprocess.run(
                ["./code"],
                input=input_data,
                text=True,
                capture_output=True,
                check=True,
            )
            total_time = time.time() - start_time
        else:
            print("Unsupported language")
            return
        
    except subprocess.CalledProcessError as e:
        print("An error occurred while executing your code:")
        print(e.stderr)
        return
    
    with open("output.txt", "r") as file:
        output_data = file.read()

    list_output_data = output_data.split("\n")
    count = 0
    actual_output = output.stdout.split("\n")
    for i in range(len(list_output_data)):
        if i >= len(actual_output):
            break
        count += 1
        if list_output_data[i] != actual_output[i]:
            count -= 1
    if count == len(list_output_data):
        print(f"{count}/{len(list_output_data)} test cases Passed :)")
    else:
        print(f"{count}/{len(list_output_data)} test cases Passed :(")
    
    print(f"Time Elapsed: {total_time*1000:.2f} msec")

if __name__ == "__main__":
    main()