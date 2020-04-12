#include <immintrin.h>

#define ALIGN(x) x __attribute__((aligned(32)))

// summing taken from https://stackoverflow.com/questions/49941645/get-sum-of-values-stored-in-m256d-with-sse-avx
// I honestly don't know how these intrisics work, but using this sum
// meant approx 6.2ns/op -> 5.3ns/op
inline
double hsum_double_avx(__m256d v) {
    __m128d vlow  = _mm256_castpd256_pd128(v);
    __m128d vhigh = _mm256_extractf128_pd(v, 1); // high 128
            vlow  = _mm_add_pd(vlow, vhigh);     // reduce down to 128

    __m128d high64 = _mm_unpackhi_pd(vlow, vlow);
    return  _mm_cvtsd_f64(_mm_add_sd(vlow, high64));  // reduce to scalar
}

// performs matrix/vector multiplication, storing the result in result.
void MultiplyMatrixByVec64(double *m, double *vec4, double *result) {
    __m256d vec = _mm256_load_pd(vec4); // load vector into register
    __m256d m1 = _mm256_load_pd(&m[0]); // load each row into a register
    __m256d m2 = _mm256_load_pd(&m[4]);
    __m256d m3 = _mm256_load_pd(&m[8]);
    __m256d m4 = _mm256_load_pd(&m[12]);

    __m256d p1 = _mm256_mul_pd(vec, m1); // multiply each row by vector using AVX2
    __m256d p2 = _mm256_mul_pd(vec, m2);
    __m256d p3 = _mm256_mul_pd(vec, m3);
    __m256d p4 = _mm256_mul_pd(vec, m4);

    double d1 = hsum_double_avx(p1); // sum each vector using AVX2
    double d2 = hsum_double_avx(p2);
    double d3 = hsum_double_avx(p3);
    double d4 = hsum_double_avx(p4);

    _mm256_storeu_pd(result, _mm256_set_pd( d4,d3,d2,d1));
}