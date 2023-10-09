from itertools import combinations

T = 20
p = 10
c = 2
r = 1.0
alpha = c/p

sc_size = 0
fs_size = 0
fp_size = 0
unsafe = 0

for s in combinations(range(1, T+1), int(c*T/p)):
    sc_size += 1

    minIns = 0
    maxIns = 0
    supplied = 0
    supply_func = [0]
    for t in range(1, T+1):
        if t > s[-1]:
            break
        if s[supplied] > 0 and t == s[supplied]:
            supplied += 1
        supply_func.append(supplied)
        ins = round(supplied-alpha*t, 5)
        if minIns > ins:
            minIns = ins
        if maxIns < ins:
            maxIns = ins
    l = len(supply_func)
    if l < T+1:
        for i in range(T+1-l):
            supply_func.append(supplied)

    # check feasible schedule
    schedulable_instance = 0
    schedulable = False
    for k in range(0, int(T/p)):
        if supply_func[(k+1)*p] - supply_func[k*p] >= c:
            schedulable_instance += 1
    if schedulable_instance == T/p:
        fs_size += 1
        schedulable = True
        # print(s)

    # check feasible partition
    if abs(round(maxIns-minIns, 5)) < r:
        if not schedulable:
            unsafe += 1
        fp_size += 1
        # print(s)

print(sc_size, fs_size, fp_size, unsafe)
