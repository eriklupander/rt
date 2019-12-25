## Perf again

| Threading | Optimization | Time | Final mem | Mallocs | Total alloc. | 
| Single | None | 3m14.846395308s | 38.56MB | 3 713 165 195 | 154.14GB |
| Multi 4/8 | None | 1m35.793805161s | 54.14MB | 3 889 554 260 | 161.53GB |
                     
| Single | Inverse cached | 10.998148908s | 29.74MB | 180 409 064 | 5.96GB |
| Single | Inverse + intersection alloc | 10.567991447s | 19.57MB | 176 292 913 | 5.70GB |                       
| Single | Inverse + intersection alloc | 10.479060737s | 27.82MB | 176 292 923 | 5.70GB |
| Multi 4/8 | Inverse cached | 4.211724075s | 42.77MB | 180 712 963 | 5.97GB |




| Multi 4/8 | Inverse cached | 4.365706658s | 29.57MB | 180712725 | 5.97GB |
| Multi 4/8 | Pixel cached | 4.201082037s | 35.28MB | 179 791 053 | 5.94GB |
| Multi 4/8 | Pixel cached #2 | 4.55495077s | 32.06MB | 179 791 136 | 5.94GB |
| Multi 4/8 | Pixel cached #3 | 4.574869591s | 31.76MB | 179 483 862 | 5.93GB |                           
| Multi 4/8 | Shade intersections | 5.329562438s | 32.83MB | 174 872 261 | 5.76GB |
| Multi 4/8 | All intersections | 5.270697573s | 35.92MB | 174 611 766 | 5.76GB |  
| Multi 4/8 | Above + caching Sphere/Plane | 3.32846761s | 42.65MB | 143672478 | 4.62GB |
| Multi 4/8 | Above + caching group ray | 3.342463845s | 41.12MB | 136926646 | 4.41GB |
| Multi 4/8 | Above + ray transform | 2.724105283s | 36.71MB | 89707668 | 3.01GB |
                                             
