import z3
import re
import json
import sys


class SMTSolver:
    def __init__(self, option, interfaces, dags, cg):
        self.verbose = option['verbose']
        self.num_slots = option['num_slots']
        self.num_channels = option['num_channels']
        self.num_apps = option['num_apps']
        self.resources = option['resources']
        self.interfaces = interfaces
        self.dags = dags
        self.cg = cg
        self.result = {}

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
        for n in self.resources:
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
        if self.verbose:
            print("added range constraint")
        return constraints

    def constraint_interface(self):
        constraints = []

        for i in range(self.num_apps):
            for n in self.interfaces[i]:
                constraints += [self.Supply[i][n][0] == 0]
                for t1 in range(1, self.num_slots):
                    ins1 = self.Supply[i][n][t1]-round(self.interfaces[i][n][0]*t1, 5)
                    for t2 in range(t1+1, self.num_slots+1):
                        ins2 = self.Supply[i][n][t2]-round(self.interfaces[i][n][0]*t2, 5)
                        constraints += [ins1-ins2 < self.interfaces[i][n][1],
                                        ins2-ins1 < self.interfaces[i][n][1]]

                constraints += [self.Supply[i][n][self.num_slots]
                                == int(self.interfaces[i][n][0]*self.num_slots)]
        if self.verbose:
            print("added interface constraint")
        return constraints

    def constraint_precedence(self):
        constraints = []

        for i in range(self.num_apps):
            D = []
            for j, dag in enumerate(self.dags[i]):
                for k in range(1, int(self.num_slots/dag["p"])+1):
                    d = [z3.Int(f"dag_{i}_{j}_{k}_{h}") for h in range(len(dag["v"]))]
                    D.append(d)

                    for h in range(len(d)):
                        constraints += [k*dag["p"] >= d[h], d[h] >= 1+(k-1)*dag["p"]]
                    for h in range(len(d)-1):
                        constraints += [d[h] < d[h+1]]
                    for h, tx in enumerate(dag["v"]):
                        for nn in tx:
                            if nn in self.interfaces[i]:
                                constraints += [z3.PbEq([(s == d[h], 1)
                                                        for si, s in enumerate(self.Slots[i][nn]) if si > 0], 1)]
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
                            if (tuple1[0] in self.interfaces[i] and tuple1[0] in tuple2) or \
                                    (tuple1[0] in self.interfaces[i] and tuple1[0] in tuple2):
                                constraints += [v1 != v2]

        if self.verbose:
            print("added precedence constraint")
        return constraints

    def constraint_transmission(self):
        constraints = []
        for i in range(self.num_apps):
            for n in self.interfaces[i]:
                if n in self.cg[i]:
                    for n2 in self.cg[i][n]:
                        for t in range(1, self.num_slots+1):
                            constraints += [z3.Not(z3.And(self.X[i][n][t], self.X[i][n2][t]))]
        if self.verbose:
            print("added conflict constraint")
        return constraints

    def constraint_channel(self):
        constraints = []
        for t in range(1, self.num_slots+1):
            constraints += [z3.Sum([z3.If(self.X[i][-1*(i+1)][t], 1, 0) for i in range(self.num_apps)]) <= self.num_channels/4]

        return constraints

    def constraint_alignment(self):
        constraints = []
        for i in range(self.num_apps):
            bw = -1*(i+1)
            for n in self.interfaces[i]:
                if n == bw:
                    continue
                # nodes
                if int(n) >= 0:
                    for si, s in enumerate(self.Slots[i][n]):
                        if si > 0:
                            constraints += [z3.Sum([z3.If(s == self.Slots[i][bw][si2], 1, 0)
                                                    for si2 in range(1, len(self.Slots[i][bw]))]) == 1]
        if self.verbose:
            print("added alignment constraint")
        return constraints

    def solve(self):
        self.create_vars()
        solver = z3.Solver()
        # solver.set("timeout", 5*60*1000)
        solver.add(self.constraint_var_range())
        solver.add(self.constraint_interface())
        solver.add(self.constraint_precedence())
        solver.add(self.constraint_transmission())
        solver.add(self.constraint_channel())
        solver.add(self.constraint_alignment())
        flag = solver.check()
        if self.verbose:
            print(flag)
        if flag == z3.sat:
            self.result = solver.model()
            if self.verbose:
                self.dump()
            return 1
        return 0

    def dump(self):
        partitions = []
        for i in range(self.num_apps):
            p = {}
            for n in self.interfaces[i]:
                p[n] = []
            partitions.append(p)

        for i in range(self.num_apps):
            for n in self.resources:
                if n in self.interfaces[i]:
                    for j in range(1, int(self.interfaces[i][n][0]*self.num_slots)+1):
                        for var in self.result:
                            if str(var) == f'slots_{i}^{n}_{j}':
                                partitions[i][n].append(self.result[var].as_long())
        print(json.dumps({"partitions": partitions}))
        # for i, p in enumerate(partitions):
        #     print(f'App-{i}: {p}')


def jsonKeys2int(x):
    if isinstance(x, dict):
        return {int(k): v for k, v in x.items()}
    return x


if __name__ == "__main__":
    option = {
        'verbose': True,
        'num_slots': 10,
        'num_channels': 4,
        'num_apps': 1,
        'resources': [-1, 1, 2, 3],
    }
    interfaces = [
        {
            -1: (0.2, 1),
            1: (0.1, 1),
            2: (0.1, 1),
            3: (0.1, 1)
        },
        # {
        #     -2: (0.2, 1),
        # },
        # {
        #     -1: (0.12, 1)
        # }
    ]
    dags = [[{"v": [(1, -99), (-99, 2)], "p":10}], [], []]
    cg = [{1: [3], 3:[1]}]

    text = '''
        {"option":{"verbose":false,"num_slots":40, "num_channels":8, "num_apps":2,"resources":[128,-1,-2]},"interfaces":[{"-1":[0.3,1],"128":[0.1,2.0834532861699016]},{"-2":[0.5,1],"128":[0.025,1]}],"dags":[[{"v":[[209,109],[109,11]],"p":40},{"v":[[128,218]],"p":10}],[{"v":[[85,128]],"p":40},{"v":[[132,36],[36,203]],"p":8}]],"cg":[{"128":[]},{"128":[]}]}
    '''
    data = json.loads(sys.argv[1])
    # data = json.loads(text)

    option = data["option"]
    interfaces = data["interfaces"]
    for i, intf in enumerate(interfaces):
        interfaces[i] = jsonKeys2int(intf)
    dags = data["dags"]
    cg = data["cg"]

    smt = SMTSolver(option, interfaces, dags, cg)
    smt.solve()
    smt.dump()
