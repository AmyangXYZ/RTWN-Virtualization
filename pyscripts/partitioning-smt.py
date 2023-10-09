import z3
import random
import multiprocessing
import re
import math
import json


class Heuristic:
    def __init__(self, option, interfaces, dags, criteria):
        self.verbose = option['verbose']
        self.num_slots = option['num_slots']
        self.num_apps = option['num_apps']
        self.num_nodes = option['num_nodes']
        self.interfaces = interfaces
        self.dags = dags
        self.criteria = criteria

        self.supply_graphs = [{} for _ in range(self.num_apps)]
        for i in range(self.num_apps):
            for n, intf in self.interfaces[i].items():
                graph = {
                    "node": n,
                    "interface": intf,
                    "supply_func": [0 for _ in range(self.num_slots+1)],
                    "insMax": 0,
                    "insMin": 0,
                    "upper_bound": [round(intf[0]*t+intf[1], 5) for t in range(self.num_slots+1)],
                    "upper_bound_soft": [round(intf[0]*t+intf[1], 5) for t in range(self.num_slots+1)],
                    "lower_bound": [round(intf[0]*t-intf[1], 5) for t in range(self.num_slots+1)],
                    "lower_bound_soft": [round(intf[0]*t-intf[1], 5) for t in range(self.num_slots+1)],
                }
                self.supply_graphs[i][n] = graph

    def make_requests(self, slot, app):
        requests = []
        graphs = self.supply_graphs[app]

        for n, g in graphs.items():
            intf = g["interface"]
            g["supply_func"][slot] = g["supply_func"][slot-1]
            if g["supply_func"][slot] < round(intf[0]*self.num_slots, 5) and \
                    g["upper_bound"][slot] > g["supply_func"][slot]+1:
                t_start_hard = t_start_soft = slot
                t_end_hard = t_end_soft = self.num_slots

                for t in range(slot, self.num_slots+1):
                    if g["upper_bound_soft"][t] >= g["supply_func"][slot]+1:
                        t_start_soft = t
                        break

                for t in range(slot, self.num_slots+1):
                    if g["lower_bound"][t] >= g["supply_func"][slot]:
                        t_end_hard = t
                        break

                for t in range(slot, self.num_slots+1):
                    if g["lower_bound_soft"][t] > g["supply_func"][slot]:
                        t_end_soft = t
                        break

                if self.criteria == "edf":
                    weight = 10/(t_end_hard-slot+1)
                elif self.criteria == "regularity":
                    weight = 10/(t_end_hard-slot+1)
                    if slot >= t_start_soft and slot <= t_end_soft:
                        weight += 10
                requests.append({
                    "app": app,
                    "node": n,
                    "weight": weight
                })
        return requests

    def update_supply_graph(self, slot, app):
        graphs = self.supply_graphs[app]
        for _, g in graphs.items():
            intf = g["interface"]
            ins = round(g["supply_func"][slot] - round(intf[0]*slot, 5), 5)
            if g["insMax"] < ins:
                g["insMax"] = ins
                for t in range(slot, self.num_slots+1):
                    g["lower_bound"][t] = round(round(intf[0]*t, 5)+g["insMax"]-intf[1], 5)
                    g["upper_bound_soft"][t] = round(round(intf[0]*t, 5)+g["insMax"], 5)

            if g["insMin"] > ins:
                g["insMin"] = ins
                for t in range(slot, self.num_slots+1):
                    g["upper_bound"][t] = round(round(intf[0]*t, 5)+g["insMin"]+intf[1], 5)
                    g["lower_bound_soft"][t] = round(round(intf[0]*t, 5)+g["insMin"], 5)

            if g["upper_bound_soft"][slot] > g["upper_bound"][slot] or g["insMax"] == 0:
                for t in range(slot, self.num_slots+1):
                    g["upper_bound_soft"][t] = g["upper_bound"][t]
            if g["lower_bound_soft"][slot] < g["lower_bound"][slot] or g["insMin"] == 0:
                for t in range(slot, self.num_slots+1):
                    g["lower_bound_soft"][t] = g["lower_bound"][t]

    def make_decision(self, slot, requests):
        requests.sort(key=lambda x: x.get('weight'), reverse=True)
        selectedApp = requests[0]["app"]
        return selectedApp

    def run(self):
        partitions = []
        for i in range(self.num_apps):
            p = {}
            for n in self.interfaces[i]:
                p[n] = []
            partitions.append(p)

        for slot in range(1, self.num_slots+1):
            requests = [[] for _ in range(self.num_nodes)]
            for i in range(self.num_apps):
                for req in self.make_requests(slot, i):
                    requests[req["node"]].append(req)

            for n in range(self.num_nodes):
                if len(requests[n]) > 0:
                    selectedApp = self.make_decision(slot, requests[n])
                    if selectedApp != -1:
                        self.supply_graphs[selectedApp][n]["supply_func"][slot] += 1
                        partitions[selectedApp][n].append(slot)

            for i in range(self.num_apps):
                self.update_supply_graph(slot, i)

        if self.verbose:
            for i, p in enumerate(partitions):
                print(f'App-{i}: {p}')
        for i, p in enumerate(partitions):
            for n, pp in p.items():
                if len(pp) < self.num_slots*self.interfaces[i][n][0]:
                    if self.verbose:
                        print("unsat")
                    return 0
                if self.supply_graphs[i][n]['insMax']-self.supply_graphs[i][n]['insMin'] >= self.interfaces[i][n][1]:
                    if self.verbose:
                        print("unsat")
                    return 0
        if self.verbose:
            print("sat")
        return 1


class RRPSolver:
    def __init__(self, option, interfaces, flag):  # flag:0-original rrp (2001), 1-magic7 (2012)
        self.verbose = option['verbose']
        self.num_slots = option['num_slots']
        self.num_apps = option['num_apps']
        self.num_nodes = option['num_nodes']
        self.interfaces = interfaces
        self.flag = flag

    def aaf(self, intf):
        intfs = []
        if intf[0] == 0:
            return intfs
        if intf[0] == 1:
            intfs.append((1, 1))
            return intfs
        if int(intf[1]) == 1:
            intfs.append((1/(2**math.floor(math.log(intf[0], 0.5))), 1))
        else:
            x = 1/(2**math.ceil(math.log(intf[0], 0.5)))
            intfs.append((x, 1))
            intfs += self.aaf((round(intf[0]-x, 5), intf[1]-1))
        return intfs

    def aaf_magic7(self, intf):
        intfs = []
        if intf[0] == 0:
            return intfs
        if intf[0] == 1:
            intfs.append((1, 1))
            return intfs
        if int(intf[1]) == 1:
            if intf[0] < 1/7:
                intfs.append((1/(7*2**(math.floor(math.log(7*intf[0], 0.5)))), 1))
            if 1/7 <= intf[0] <= 6/7:
                intfs.append((math.ceil(7*intf[0])/7, 1))
            if 6/7 < intf[0] < 1:
                intfs.append((1-1/(7*2**(math.ceil(math.log(7*(1-intf[0]), 0.5)))), 1))
        else:
            c = 0
            if intf[0] < 1/7:
                c = 1/(7*2**(math.ceil(math.log(7*intf[0], 0.5))))
            if 1/7 <= intf[0] <= 6/7:
                c = math.floor(7*intf[0])/7
            if 6/7 < intf[0] < 1:
                c = 1-1/(7*2**(math.floor(math.log(7*(1-intf[0]), 0.5))))
            intfs.append((c, 1))
            intfs += self.aaf_magic7((round(intf[0]-c, 5), intf[1]-1))
        return intfs

    def aaf_mulz(self, x, intf):
        intfs = []
        if intf[0] == 0:
            return intfs
        if intf[0] == 1:
            intfs.append((1, 1))
            return intfs
        if int(intf[1]) == 1:
            if intf[0] < 1/x:
                intfs.append((1/(x*2**(math.floor(math.log(x*intf[0], 0.5)))), 1))
            if 1/x <= intf[0] <= 1-1/x:
                intfs.append((math.ceil(x*intf[0])/x, 1))
            if 1-1/x < intf[0] < 1:
                intfs.append((1-1/(x*2**(math.ceil(math.log(x*(1-intf[0]), 0.5)))), 1))
        else:
            c = 0
            if intf[0] < 1/x:
                c = 1/(7*2**(math.ceil(math.log(7*intf[0], 0.5))))
            if 1/x <= intf[0] <= 1-1/x:
                c = math.floor(x*intf[0])/x
            if 1-1/x < intf[0] < 1:
                c = 1-1/(x*2**(math.floor(math.log(x*(1-intf[0]), 0.5))))
            intfs.append((c, 1))
            intfs += self.aaf_mulz(x, (round(intf[0]-c, 5), intf[1]-1))
        return intfs

    def run(self):
        partitions = []
        for n in range(self.num_nodes):
            if self.flag == 2:
                for x in [2, 3, 4, 5, 7]:
                    partitions = []
                    for i in range(self.num_apps):
                        if n in self.interfaces[i]:
                            partitions += self.aaf_mulz(x, self.interfaces[i][n])
                    if sum(intf[0] for intf in partitions) <= 1:
                        if self.verbose:
                            print("sat")
                        return 1
                return 0
            else:
                for i in range(self.num_apps):
                    if n in self.interfaces[i]:
                        if self.flag == 0:
                            partitions += self.aaf(self.interfaces[i][n])
                        elif self.flag == 1:
                            partitions += self.aaf_magic7(self.interfaces[i][n])

                if sum(intf[0] for intf in partitions) > 1:
                    if self.verbose:
                        print("unsat")
                    return 0
        if self.verbose:
            print("sat")
        return 1


class SMTSolver:
    def __init__(self, option, interfaces, dags):
        self.verbose = option['verbose']
        self.num_slots = option['num_slots']
        self.num_apps = option['num_apps']
        self.num_nodes = option['num_nodes']
        self.interfaces = interfaces
        self.dags = dags

    def create_vars(self):
        self.X = []
        self.Supply = []
        self.Slots = []
        for i in range(self.num_apps):
            x = {}
            supply = {}
            slots = {}
            for n in self.interfaces[i]:
                x[n] = [z3.Bool(f'x_{i}^{n}({t})') for t in range(self.num_slots+1)]
                supply[n] = [z3.Int(f'supply_{i}^{n}({t})') for t in range(self.num_slots+1)]
                slots[n] = [z3.Int(f'slots_{i}^{n}_{j}') for j in
                            range(0, int(self.interfaces[i][n][0]*self.num_slots)+1)]
            self.X.append(x)
            self.Supply.append(supply)
            self.Slots.append(slots)

    def constraint_var_range(self):
        constraints = []
        for n in range(self.num_nodes):
            # X mutually exclusive
            for t in range(1, self.num_slots+1):
                xx = []
                for i in range(self.num_apps):
                    if n in self.interfaces[i]:
                        xx.append((self.X[i][n][t], 1))
                constraints += [z3.PbLe(xx, 1)]

            # map X to Supply
            for t in range(1, self.num_slots+1):
                for i in range(self.num_apps):
                    if n in self.interfaces[i]:
                        constraints += [z3.If(
                            self.X[i][n][t],
                            self.Supply[i][n][t] == self.Supply[i][n][t-1]+1,
                            self.Supply[i][n][t] == self.Supply[i][n][t-1]
                        )]

        # map Supply to Slots
        for i in range(self.num_apps):
            for n in self.interfaces[i]:
                for j in range(1, int(self.interfaces[i][n][0]*self.num_slots)+1):
                    for t in range(1, self.num_slots+1):
                        constraints += [z3.If(
                            z3.And(self.Supply[i][n][t] == j, self.Supply[i][n][t-1] == j-1),
                            self.Slots[i][n][j] == t,
                            True
                        )]

        return constraints

    def constraint_interface(self):
        constraints = []
        for n in range(self.num_nodes):
            # instant regularity
            for i in range(self.num_apps):
                if n in self.interfaces[i]:
                    constraints += [self.Supply[i][n][0] == 0]
                    for t1 in range(1, self.num_slots+1):
                        ins1 = self.Supply[i][n][t1]-round(self.interfaces[i][n][0]*t1, 5)
                        for t2 in range(1, self.num_slots+1):
                            if t1 == t2:
                                continue
                            ins2 = self.Supply[i][n][t2]-round(self.interfaces[i][n][0]*t2, 5)
                            constraints += [ins1-ins2 < self.interfaces[i][n][1],
                                            ins2-ins1 < self.interfaces[i][n][1]]

                    constraints += [self.Supply[i][n][self.num_slots]
                                    == self.interfaces[i][n][0]*self.num_slots]
        return constraints

    def constraint_precedence(self):
        constraints = []

        for i in range(self.num_apps):
            D = []
            for j, dag in enumerate(self.dags[i]):
                for k in range(1, int(self.num_slots/dag["period"])+1):
                    d = [z3.Int(f"dag_{i}_{j}_{k}_{h}") for h in range(len(dag["v"]))]
                    D.append(d)

                    for h in range(len(d)-1):
                        if h == 0:
                            constraints += [d[h] >= 1+(k-1)*dag["period"]]
                        if h+1 == len(d)-1:
                            constraints += [d[h+1] <= k*dag["period"]]
                        constraints += [d[h] < d[h+1]]
                    for h, tx in enumerate(dag["v"]):
                        for nn in tx:
                            if nn != -1:
                                constraints += [z3.PbEq([(s == d[h], 1)
                                                        for si, s in enumerate(self.Slots[i][nn]) if si >= 1], 1)]
            for d1i in range(len(D)-1):
                d1 = D[d1i]
                for d2i in range(d1i+1, len(D)):
                    d2 = D[d2i]
                    for v1 in d1:
                        _, j1, _, h1 = re.findall(r'\d+', str(v1))
                        tuple1 = self.dags[i][int(j1)]["v"][int(h1)]
                        for v2 in d2:
                            _, j2, _, h2 = re.findall(r'\d+', str(v2))
                            tuple2 = self.dags[i][int(j2)]["v"][int(h2)]
                            if (tuple1[0] != -1 and tuple1[0] in tuple2) or \
                                    (tuple1[0] != -1 and tuple1[0] in tuple2):
                                constraints += [v1 != v2]
        return constraints

    def constraint_transmission(self):
        constraints = []
        return constraints

    def solve(self):
        self.create_vars()
        solver = z3.Solver()
        solver.set("timeout", 10*60*1000)
        solver.add(self.constraint_var_range())
        solver.add(self.constraint_interface())
        # solver.add(self.constraint_precedence())
        solver.add(self.constraint_transmission())
        flag = solver.check()
        # if self.verbose:
        print(flag)
        if flag == z3.sat:
            if self.verbose:
                self.dump(solver.model())
            return 1
        return 0

    def dump(self, result):
        partitions = []
        for i in range(self.num_apps):
            p = {}
            for n in self.interfaces[i]:
                p[n] = []
            partitions.append(p)

        for i in range(self.num_apps):
            for n in range(self.num_nodes):
                if n in self.interfaces[i]:
                    # for t in range(1, self.num_slots+1):
                    for j in range(1, int(self.interfaces[i][n][0]*self.num_slots)+1):
                        for var in result:
                            if str(var) == f'slots_{i}^{n}_{j}':
                                partitions[i][n].append(result[var])
        for i, p in enumerate(partitions):
            print(f'App-{i}: {p}')


class Trial:
    def __init__(self, option):
        self.option = option
        pass

    def run(self, trial_id):
        random.seed(123+trial_id)

        interfaces = [{} for _ in range(self.option["num_apps"])]

        for n in range(self.option["num_nodes"]):
            n_slot_list = generate_random_ints(self.option['num_apps'], int(
                self.option['num_slots']*self.option["utilization"]))
            for i in range(self.option["num_apps"]):
                interfaces[i][n] = (n_slot_list[i]/self.option["num_slots"],
                                    random.randint(101, 150)/100)
                # self.option["regularity"])
        # interfaces = [
        #     {
        #         0: (0.15, 1),
        #         # 1: (0.1, 1),
        #     },
        #     {
        #         0: (0.2, 1),
        #         # 1: (0.1, 1),
        #     },
        #     # {
        #     #     0: (0.15, 1),
        #     #     # 1: (0.2, 2),
        #     # },
        # ]
        dags = [
            [
                # {
                #     "v": [(0, -1), (1, -1)],
                #     "period": 10,
                # },
                # {
                #     "v": [(0, 1)],
                #     "period": 20,
                # }
            ],
            [
                # {
                #     "v": [(0, 1)],
                #     "period": 10,
                # }
            ],
            [],
            [],
            [],
            []
        ]
        if self.option["verbose"]:
            print(interfaces)
            # print(dags)

        smt = SMTSolver(self.option, interfaces, dags)
        return smt.solve()


def generate_random_ints(num_ints, target_sum):
    random_ints = [random.randint(0, target_sum) for _ in range(num_ints - 1)]
    random_ints.append(0)
    random_ints.append(target_sum)

    random_ints.sort()
    result = [random_ints[i + 1] - random_ints[i] for i in range(num_ints)]

    return result


if __name__ == "__main__":
    num_cores = multiprocessing.cpu_count()

    data = {"smt": [], "aaf": [], "magic7": [], "mulz": [], "heu_edf": [], "heu_reg": []}
    xAxis = []
    need_smt = {
        0.80: [109, 113, 141, 184],
        0.85: [6, 9, 10, 14, 99, 109, 113],
        0.90: [7, 13, 33, 38, 47, 49, 50, 67, 96, 103, 104, 109, 128, 196],
        0.95: [13, 26, 38, 44, 95, 97, 113, 128, 178, 194, 195],
        1.00: [4, 13, 24, 26, 28, 32, 33, 38, 42, 44, 54, 72, 76, 79, 80, 86, 96, 97, 99, 100, 112, 113, 123, 130, 131, 134, 135, 141, 142, 146, 148, 162, 166, 168, 169, 180, 189]
    }
    for u in need_smt:
        pool = multiprocessing.Pool(processes=num_cores)
        option = {
            'verbose': False,
            'num_slots': 28,
            'num_apps': 6,
            'num_nodes': 1,
            'regularity': 1,
            'utilization': u
        }
        xAxis.append(u)
        t = Trial(option)
        res = pool.map(t.run, need_smt[u])
        print(100-(len(res)-sum(res)))
