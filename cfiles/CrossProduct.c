#include <immintrin.h>

#define ALIGN(x) x __attribute__((aligned(32)))

// adapted from https://gist.github.com/garrettsickles/85a9ab8385172bd0e762f38e4cfb045f
void CrossProduct(double* a,double* b, double* result) {
    __m256d vec1 = _mm256_load_pd(a);
    __m256d vec2 = _mm256_load_pd(b);

    __m256d c =_mm256_permute4x64_pd(_mm256_sub_pd(
    			_mm256_mul_pd(vec1, _mm256_permute4x64_pd(vec2, _MM_SHUFFLE(3, 0, 2, 1))),
    			_mm256_mul_pd(vec2, _mm256_permute4x64_pd(vec1, _MM_SHUFFLE(3, 0, 2, 1)))
    		), _MM_SHUFFLE(3, 0, 2, 1));

    _mm256_storeu_pd(result, c);
}