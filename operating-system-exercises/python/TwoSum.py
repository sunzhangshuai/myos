from typing import List


def twoSum(nums: List[int], target: int) -> List[int]:
    for num1 in nums:
        for num2 in nums:
            if num1 != num2 and num1 + num2 == target:
                print(num1, num2)


if __name__ == "__main__":
    a = [1, 2, 3, 4, 5, 7]
    twoSum(a, 10)