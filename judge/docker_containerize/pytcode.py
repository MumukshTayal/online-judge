# t = int(input())
# for _ in range(t):
#     l = list(map(int, input().split()))
#     n, target = l[0], l[1]
#     arr = list(map(int, input().split()))
#     hist = {}
#     for i in range(len(arr)):
#         if (target - arr[i]) in hist:
#             print(hist[target - arr[i]], i)
#             break
#         else:
#             hist[arr[i]] = i

# print("Hello World!")

t = int(input())
for _ in range(t):
    l = list(map(int, input().split()))
    n, target = l[0], l[1]
    arr = list(map(int, input().split()))
    hist = {}
    for i in range(len(arr)):
        if (target - arr[i]) in hist:
            print(hist[target - arr[i]], i)
            break
        else:
            hist[arr[i]] = i