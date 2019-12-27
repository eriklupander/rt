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
                                             
NEW image after fixing refraction!!

| 4/8 | before cache group xs | 2.474886316s | 42.16MB | 68355717 | 2.45GB |
| 4/8 | after cache group xs | 2.651449543s | 42.16MB | 68138172 | 2.43GB |
| 4/8 | after cache sphere xs | 2.545900359s | 64.21MB | 66940149 | 2.34GB |
| 4/8 | after cache inner rays in group | 2.382903549s | 56.14MB | 52657634 | 1.92GB |
| 4/8 | after compsnegate | 2.349631789s | 38.15MB | 51919729 | 1.89GB |
| 4/8 | after cache normalvec | 2.268705698s | 60.10MB | 51553365 | 1.88GB |
| 4/8 | after cache containers | 2.275932712s | 58.70MB | 48790198 | 1.81GB |
| 4/8 | after mul on tuple | 2.246734595s | 43.40MB | 48183746 | 1.79GB |
| 4/8 | after more mul on tuple | 2.237514183s | 53.22MB | 47326488 | 1.76GB |
| 4/8 | after add fix | 2.140603628s | 62.17MB | 46136453 | 1.73GB |
| 4/8 | after over/under | 2.103660441s  | 46.17MB | 43760459 | 1.66GB |
| 4/8 | after remove sort in hit | 2.016085178s  | 34.80MB | 40923569 | 1.54GB |
| 4/8 | fix adds in light | 1.969865723s | 54.66MB | 39243360 | 1.49GB |
| 4/8 | fix lightVec in light | 2.216685719s | 44.74MB | 35903401 | 1.39GB |
| 4/8 | cylinder xs cache | 1.898069389s | 49.29MB | 33376106 | 1.31GB |

                                

                                                                                                                                       