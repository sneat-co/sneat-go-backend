import sys

coverage = {}
try:
    with open('coverage.out', 'r') as f:
        lines = f.readlines()
        if not lines:
            sys.exit(0)
        for line in lines[1:]:
            parts = line.split()
            if len(parts) != 3: continue
            block, num_stmt, count = parts
            pkg = '/'.join(block.split('/')[:-1])
            if pkg not in coverage:
                coverage[pkg] = {'total': 0, 'covered': 0}
            num_stmt = int(num_stmt)
            count = int(count)
            coverage[pkg]['total'] += num_stmt
            if count > 0:
                coverage[pkg]['covered'] += num_stmt

    results = []
    for pkg, data in coverage.items():
        uncovered = data['total'] - data['covered']
        results.append((uncovered, data['covered'], data['total'], pkg))
    
    results.sort(key=lambda x: x[0], reverse=True)
    for res in results[:20]:
        print(f"{res[0]}\t{res[1]}\t{res[2]}\t{res[3]}")
except Exception as e:
    print(e)
