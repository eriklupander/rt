## Perf again
Single: 66s
Threaded 30.5s

With Inverse cached on camera and shapes
Single: 5.1s
Threaded: 2.5s

With intersection list cached

Mallocs:

With [:0]: 480 / 480 Memory: 23.68MB Mallocs: 89249818 Total alloc: 3.21GB
With nil:  480 / 480 Memory: 23.86MB Mallocs: 89959129 Total alloc: 3.28GB
With Make: 480 / 480 Memory: 31.56MB Mallocs: 90868175 Total alloc: 3.31GB


## Final image
| Threading | Optimization | Time | Final mem | Mallocs | Total alloc. | 
| Single | None | 3m14.846395308s | 38.56MB | 3 713 165 195 | 154.14GB |
| Single | Inverse cached | 10.998148908s | 29.74MB | 180 409 064 | 5.96GB |
| Single | Inverse + intersection alloc | 10.567991447s | 19.57MB | 176 292 913 | 5.70GB |                       
| Single | Inverse + intersection alloc | 10.479060737s | 27.82MB | 176 292 923 | 5.70GB |
| Multi 4/8 | Inverse cached | 4.211724075s | 42.77MB | 180 712 963 | 5.97GB |


