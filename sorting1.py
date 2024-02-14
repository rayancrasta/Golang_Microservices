import math

def bubblesort(customlist): # TC(O(N2))
    # compare adjacent elements and swap them if they are in wrong order
    for i in range(len(customlist)-1):
        for j in range(len(customlist)-1-i):
            if customlist[j] > customlist[j+1]:
                customlist[j],customlist[j+1] = customlist[j+1],customlist[j]
    print(customlist)
    # the highest value basically bubbles up to the last position
    # When to use: Input is like almost sorted,space is concern,easy to implement
    # When to avoid: Average time complexity is poor

def selectionsort(customlist): # TC(ON2)
   # We repeatedly find the minimum element and move it to its position from 0... n 
    for i in range(len(customlist)):
        min_index = i
        for j in range(i+1,len(customlist)):
            if customlist[min_index] > customlist[j]: # list[min] is not minimum value in the list
                min_index = j
        #swap the min value of this round
        customlist[i],customlist[min_index] = customlist[min_index],customlist[i]
    print(customlist)
    
    # when to use: We have insufficient memory
    # when not to use: Time is concern
    
def insertionsort(customlist): # TC(ON2)
    # take element from unsorted part, and put in correct position in sorted part
    for i in range(1,len(customlist)):
        key = customlist[i]
        j = i-1
        while j >= 0 and key < customlist[j]:
            customlist[j+1] = customlist[j] # swaping to put in order, within the sorted list
            j-=1
        # if above conditions are not met, list[i] is already > than i-1,then means 
        customlist[j+1]=key # j+1 is i , nothing changes in that case 
    return (customlist)
            
    # when to use: we have insufficient memory
    # we have continous inflow of elements and we want to keep them sorted
    
    # not to use: time is concern     

def bucketsort(customlist): # TC(ON2) , SC o(N)
    numberofbuckets = round(math.sqrt(len(customlist)))
    maxvalue = max(customlist)
    arr = []
    
    #create buckets
    for i in range(numberofbuckets):
        arr.append([])
    
    #put item in approp bucket
    for j in customlist:
        index_b = math.ceil(j*numberofbuckets/maxvalue)
        arr[index_b-1].append(j) # put item in approp bucket
        
    #sort individual buckets
    for i in range(numberofbuckets):
        arr[i] = insertionsort(arr[i])
        
    #merge the buckets
    k = 0
    for i in range(numberofbuckets):
        for j in range(len(arr[i])):
            customlist[k] = arr[i][j]
            k+=1
    
    print(customlist)
    # when to use: when input data is uniformly distributed over a range
    # unformily dist: in a similar ange like 1,2,4,5,3,7,9 not 1,2,4,91,94
    # not: when space is concern

def bucketSortNegative(customList):
    numberofBuckets = round(math.sqrt(len(customList)))
    minValue = min(customList)
    maxValue = max(customList)
    rangeVal = (maxValue - minValue) / numberofBuckets
 
    buckets = [[] for _ in range(numberofBuckets)]
 
    for j in customList:
        if j == maxValue:
            buckets[-1].append(j)
        else:
            index_b = math.floor((j - minValue) / rangeVal)
            buckets[index_b].append(j)
    
    sorted_array = []
    for i in range(numberofBuckets):
        buckets[i] = insertionSort(buckets[i])
        sorted_array.extend(buckets[i])
    
    return sorted_array

def merge(customlist,first,middle,last):
    n1 = middle - first + 1 
    n2 = last - middle
    
    #Create 2 sub arrays
    L = [0] * (n1)
    R = [0] * (n2)
    
    # fill first sub array
    for i in range(0,n1):
        L[i] = customlist[first+i]
        
    # fill second sub array
    for j in range(0,n2):
        R[j] = customlist[middle+j+1]    

    i = 0  # initial index of first sub array
    j = 0  # initial index of second sub array
    k = first  #initial index of merged sub array
    
    while i < n1 and j < n2:
        if L[i] <= R[j]:
            customlist[k] = L[i]
            i +=1
        else:
            customlist[k] = R[j]
            j+=1
        k+=1
        
    # Copy elements of L and R if any are remaining 
    while i < n1 :
        customlist[k] = L[i]
        i+=1
        k+=1
    
    while j < n2:
        customlist[k] = R[j]
        j+=1
        k+=1
        

def mergesort(customlist,leftindex,rightindex): # TC (O(NlogN)) from masters theorem
    if leftindex < rightindex:
        middleindex = (leftindex + (rightindex-1))//2 # why -1 
        mergesort(customlist,leftindex,middleindex) # TC O(n/2)
        mergesort(customlist,middleindex+1,rightindex) # TC O(n/2)
        merge(customlist,leftindex,middleindex,rightindex)  # TC O(n))
    return customlist

    # TC NlogN , SC N
    # when we need stable sort
    # not when space is concern
    
clist = [2,1,7,3,4,9,8]
# bubblesort(clist)
# selectionsort(clist)    
# insertionsort(clist)
# bucketsort(clist)
print(mergesort(clist,0,6))