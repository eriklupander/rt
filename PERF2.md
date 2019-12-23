## Perf again

| Threading | Optimization | Time | Final mem | Mallocs | Total alloc. | 
| Single | None | 3m14.846395308s | 38.56MB | 3 713 165 195 | 154.14GB |
| Multi 4/8 | None | 1m35.793805161s | 54.14MB | 3 889 554 260 | 161.53GB |
                     
| Single | Inverse cached | 10.998148908s | 29.74MB | 180 409 064 | 5.96GB |
| Single | Inverse + intersection alloc | 10.567991447s | 19.57MB | 176 292 913 | 5.70GB |                       
| Single | Inverse + intersection alloc | 10.479060737s | 27.82MB | 176 292 923 | 5.70GB |
| Multi 4/8 | Inverse cached | 4.211724075s | 42.77MB | 180 712 963 | 5.97GB |


