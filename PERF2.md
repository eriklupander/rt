## Perf again

| Threading | Optimization | Time | Final mem | Mallocs | Total alloc. | 
| Single | None | 3m14.846395308s | 38.56MB | 3 713 165 195 | 154.14GB |
| Single | Inverse cached | 10.998148908s | 29.74MB | 180 409 064 | 5.96GB |
| Single | Inverse + intersection alloc | 10.567991447s | 19.57MB | 176 292 913 | 5.70GB |                       
| Single | Inverse + intersection alloc | 10.479060737s | 27.82MB | 176 292 923 | 5.70GB |
| Single | All above + worldToPoint ptr | 10.365876513 | 23.79MB | 174227818 | 5.64GB |
| Single | All above + worldToNormal | 10.482145506s | 29.86MB | 174999304 | 5.67GB |
| Single | All above + NormalizePtr | 10.184640323s | 29.84MB | 169309341 | 5.49GB |
| Single | All above + ray transform | 8.657896898s | 30.94MB 125618139 | 4.19GB |
| Single | All above + group ray transform | 8.372439644s | 36.57MB | 118674025 | 3.99GB |
| Single | All above + inner group rays | 7.446529206s | 30.81MB | 90897129 | 3.16GB |
| Single | All above + Sub in sphere xs | 6.956806346s | 18.62MB | 70064378 | 2.54GB |
| Single | All + alloc in plane xs | 6.578297633s | 24.97MB | 63126542 | 2.23GB |
| Single | All + alloc in cube | 6.516764324s | 32.72MB | 63076077 | 2.22GB |
| Single | All + color allocs in light | 6.500401359s | 32.90MB | 61273780 | 2.17GB |
| Single | All + lightVec in light | 6.481354709s | 26.38MB | 60372628 | 2.14GB |
| Single | All + adds in light | 6.453733252s | 27.17MB | 58133700 | 2.08GB |
| Single | All + more in light | 6.426145991s | 30.21MB | 56343100 | 2.02GB |
                                                                                                                                                                                                                      
                                     
                                                                                                                                                                                

| Multi 4/8 | Inverse cached | 4.211724075s | 42.77MB | 180 712 963 | 5.97GB |


