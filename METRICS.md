|Name | Mem | Allocs | Total alloc | Duration | XS skipped | Transpose |
|---------------| ------------- |-------------| -----| --------|----------|----------|
| Before#1 |103.45MB|31242256|1.15GB|9.887840037s|21451523||
| Before#2 |103.68MB|31242246|1.15GB|8.282383601s|21451523||
| Before#3 |106.82MB|31242301|1.15GB|7.839726515s|21451523||
| Cache transpose#1 |94.99MB|28704976|1.06GB|7.411978187s|20461309||
| Cache transpose#2 |120.60MB|28704996|1.06GB|7.753794886s|20461309||
| Cache transpose#3 |91.61MB|28705033|1.06GB|7.688985346s|20461309||
| Cache transpose no groups | 125.86MB|31242234|1.15GB|8.057420715s|21451523||
| Cache transpose |131.47MB|31242272|1.15GB|8.167282805s|21451523||
| Cache transpose |106.63MB|31242255|1.15GB|7.53470478s|21451523||
| Break shadow test |129.51MB|21857829|888.41MB|6.782154879s|0|136|

|Name | Mem | Allocs | Total alloc | Duration | XS skipped | Transpose |
|---------------| ------------- |-------------| -----| --------|----------|----------|
| Smooth#1 |251.94MB|16444120|1.14GB|20.683699268s|0|36|
| Smooth#2 |203.22MB|13043039|1023.14MB|22.384036467s|0|36|
| Smooth#3 |177.33MB|13042215|1023.11MB|23.531206572s|0|36|
| Smooth#4 new | 218.96MB|13043028|1023.14MB|20.752756273s|0|36|
| Smooth#5 |179.02MB|13043275|1014.01MB|20.343711954s|0|36|
| Smooth#6 |190.25MB|13050454|1017.83MB|20.488117903s|0|36|
| Cross2 |217.85MB|13048887|1.00GB|21.237598691s|0|36|
| Inline sub |227.09MB|13050751|1.00GB|19.265662085s|0|36|
| Inline sub without reuse |271.56MB|13042883|1013.99MB|20.19283564s|0|36|