# Implement string matching algorithm based on ngram
# http://www.catalysoft.com/articles/StrikeAMatch.html

# Algorithm overview: 
# 1. split strings into pairs of two consecutive characters (case insentive)
# 2. calculate similarity score  
# 3. s(x,y) = 2 * (pairs(x) intersect pairs(y)) / (ngram(x) + ngram(y))
# 4. value of s will be between 0 and 1

from stops import Busstops # List of bus stops
from itertools import zip_longest
pairs = []
for stop in Busstops:
    s = stop.lower()
    pairs.append(([s[i:i+2] for i in range(len(s)-1)], stop))


def similarity(input):
    result = []
    input = input.lower()
    pair_input = set([input[i:i+2] for i in range(len(input)-1)])
    for stop, name in pairs:
        n_pairs_input = len(pair_input)
        n_pairs_stop = len(stop)
        n_pairs_intersect = len(pair_input & set(stop))
        score = 2 * n_pairs_intersect / (n_pairs_input + n_pairs_stop)
        result.append((score, name))
    result.sort(key=lambda x: x[0], reverse=True)
    return result

print(similarity("elis")[:5])
